package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type errorHandle func(http.ResponseWriter, *http.Request, httprouter.Params) error

type httpKernel struct {
	project *Project
}

func (kernel *httpKernel) strap(handle errorHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := handle(w, r, ps); err != nil {
			kernel.renderError(w, r, err)
		}
	}
}

func (kernel *httpKernel) renderError(w http.ResponseWriter, r *http.Request, err error) {
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("artifact not found"))
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("%v %v: %v\n", r.Method, r.URL.String(), err)
}

func (kernel *httpKernel) renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func (kernel *httpKernel) renderMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("method not allowed"))
}

func (kernel *httpKernel) collateArtifacts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	rawArtifactsValue := r.URL.Query().Get("artifacts")
	if rawArtifactsValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	return kernel.project.CollateArtifacts(strings.Split(rawArtifactsValue, ","), w)
}

func (kernel *httpKernel) listArtifacts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	artifactKeys, err := kernel.project.ListArtifacts()
	if err != nil {
		return err
	}
	sort.Strings(artifactKeys)
	return json.NewEncoder(w).Encode(artifactKeys)
}

func (kernel *httpKernel) showArtifact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	key := ps.ByName("artifact")
	if err := kernel.project.CollateArtifacts([]string{key}, w); err != nil {
		return err
	}
	return nil
}

func (kernel *httpKernel) storeArtifact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	key := ps.ByName("artifact")
	if err := kernel.project.StoreArtifact(key, r.Body); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (kernel *httpKernel) removeArtifact(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	key := ps.ByName("artifact")
	if err := kernel.project.RemoveArtifact(key); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (kernel *httpKernel) router() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(kernel.renderNotFound)
	router.MethodNotAllowed = http.HandlerFunc(kernel.renderMethodNotAllowed)

	router.GET("/collation", kernel.strap(kernel.collateArtifacts))
	router.GET("/artifacts", kernel.strap(kernel.listArtifacts))
	router.GET("/artifacts/:artifact", kernel.strap(kernel.showArtifact))
	router.PUT("/artifacts/:artifact", kernel.strap(kernel.storeArtifact))
	router.DELETE("/artifacts/:artifact", kernel.strap(kernel.removeArtifact))

	return router
}
