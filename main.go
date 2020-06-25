package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	artifactRepository := &localArtifactRepository{
		basePath: "data/default",
	}
	if err := os.MkdirAll(artifactRepository.basePath, os.ModePerm); err != nil {
		panic(err)
	}

	project := &Project{
		artifactRepository: artifactRepository,
	}

	kernel := &httpKernel{project}

	log.Fatal(http.ListenAndServe(":8080", kernel.router()))
}
