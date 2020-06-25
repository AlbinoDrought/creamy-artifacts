package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

type sendFunc func(req *http.Request) *http.Response
type testFunc func(send sendFunc)

func setUp(callback testFunc) error {
	basePath := path.Join(os.TempDir(), "AlbinoDrought/creamy-artifacts/http_kernel_test.go")
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return err
	}
	defer os.RemoveAll(basePath)

	artifactRepository := &localArtifactRepository{basePath}
	project := &Project{artifactRepository}
	kernel := &httpKernel{project}
	router := kernel.router()

	callback(func(req *http.Request) *http.Response {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		return resp.Result()
	})

	return nil
}

func assertStatus(t *testing.T, resp *http.Response, expectedStatus int, action string) {
	if resp.StatusCode != expectedStatus {
		t.Errorf("expected %v after %v but received %v", expectedStatus, action, resp.StatusCode)
	}
}

func assertBody(t *testing.T, resp *http.Response, expectedBody []byte, action string) {
	actualBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("unexpected error when reading body after %v", action)
	}
	if bytes.Compare(expectedBody, actualBody) != 0 {
		t.Errorf("expected body of %v after %v but received %v (\"%v\" vs \"%v\")", expectedBody, action, actualBody, string(expectedBody), string(actualBody))
	}
}

func TestWriteReadCollateDelete(t *testing.T) {
	err := setUp(func(send sendFunc) {
		var resp *http.Response

		firstArtifact := "we should stick to the plan"
		secondArtifact := "we've come to stay"
		thirdArtifact := "we want a slice of the cake"

		// write #1
		resp = send(httptest.NewRequest("PUT", "http://artifacts.localhost/artifacts/v1.2.0", strings.NewReader(firstArtifact)))
		assertStatus(t, resp, http.StatusNoContent, "successful write #1")

		// list #1
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts", nil))
		assertStatus(t, resp, http.StatusOK, "successful list #1")
		assertBody(t, resp, []byte("[\"v1.2.0\"]\n"), "successful list #1")

		// write #2
		resp = send(httptest.NewRequest("PUT", "http://artifacts.localhost/artifacts/v1.3.0", strings.NewReader(secondArtifact)))
		assertStatus(t, resp, http.StatusNoContent, "successful write #2")

		// list #2
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts", nil))
		assertStatus(t, resp, http.StatusOK, "successful list #2")
		assertBody(t, resp, []byte("[\"v1.2.0\",\"v1.3.0\"]\n"), "successful list #2")

		// write #3
		resp = send(httptest.NewRequest("PUT", "http://artifacts.localhost/artifacts/v1.3.1", strings.NewReader(thirdArtifact)))
		assertStatus(t, resp, http.StatusNoContent, "successful write #3")

		// list #3
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts", nil))
		assertStatus(t, resp, http.StatusOK, "successful list #3")
		assertBody(t, resp, []byte("[\"v1.2.0\",\"v1.3.0\",\"v1.3.1\"]\n"), "successful list #2")

		// read #1
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.2.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful read #1")
		assertBody(t, resp, []byte(firstArtifact), "successful read #1")

		// read #2
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.3.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful read #2")
		assertBody(t, resp, []byte(secondArtifact), "successful read #2")

		// read #3
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.3.1", nil))
		assertStatus(t, resp, http.StatusOK, "successful read #3")
		assertBody(t, resp, []byte(thirdArtifact), "successful read #3")

		// collation #1: 1, single things are ok
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/collation?artifacts=v1.2.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful collation #1")
		assertBody(t, resp, []byte(firstArtifact), "successful collation #1")

		// collation #2: 1+2+3, multiple things are ok
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/collation?artifacts=v1.2.0,v1.3.0,v1.3.1", nil))
		assertStatus(t, resp, http.StatusOK, "successful collation #2")
		assertBody(t, resp, append([]byte(firstArtifact), append([]byte(secondArtifact), []byte(thirdArtifact)...)...), "successful collation #2")

		// collation #3: 3+1+2, order is important
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/collation?artifacts=v1.3.1,v1.2.0,v1.3.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful collation #3")
		assertBody(t, resp, append([]byte(thirdArtifact), append([]byte(firstArtifact), []byte(secondArtifact)...)...), "successful collation #3")

		// delete #1
		resp = send(httptest.NewRequest("DELETE", "http://artifacts.localhost/artifacts/v1.2.0", nil))
		assertStatus(t, resp, http.StatusNoContent, "successful delete #1")

		// read-after-delete #1 should fail
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.2.0", nil))
		assertStatus(t, resp, http.StatusNotFound, "read-after-delete #1")

		// read undeleted #2
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.3.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful read of undeleted #2")
		assertBody(t, resp, []byte(secondArtifact), "successful read of undeleted #2")

		// read undeleted #3
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.3.1", nil))
		assertStatus(t, resp, http.StatusOK, "successful read of undeleted #3")
		assertBody(t, resp, []byte(thirdArtifact), "successful read of undeleted #3")

		// overwrite #2 with #1
		resp = send(httptest.NewRequest("PUT", "http://artifacts.localhost/artifacts/v1.3.0", strings.NewReader(firstArtifact)))
		assertStatus(t, resp, http.StatusNoContent, "successful overwrite #2 with #1")
		resp = send(httptest.NewRequest("GET", "http://artifacts.localhost/artifacts/v1.3.0", nil))
		assertStatus(t, resp, http.StatusOK, "successful read of overwritten #2")
		assertBody(t, resp, []byte(firstArtifact), "successful read of overwritten #2")

		// delete #2
		resp = send(httptest.NewRequest("DELETE", "http://artifacts.localhost/artifacts/v1.3.0", nil))
		assertStatus(t, resp, http.StatusNoContent, "successful delete #2")

		// delete #3
		resp = send(httptest.NewRequest("DELETE", "http://artifacts.localhost/artifacts/v1.3.1", nil))
		assertStatus(t, resp, http.StatusNoContent, "successful delete #3")

		// re-delete #3, should 404
		resp = send(httptest.NewRequest("DELETE", "http://artifacts.localhost/artifacts/v1.3.1", nil))
		assertStatus(t, resp, http.StatusNotFound, "re-delete #3")
	})

	if err != nil {
		t.Error(err)
	}
}
