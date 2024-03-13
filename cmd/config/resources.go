package config

import (
	"encoding/json"
	"io"
	"os"

	"autenticacion-ms/cmd/config/model"
)

func GetArtifactResources(pathArtifactResources string) model.ArtifactResources {
	resourceFile, err := os.Open(pathArtifactResources)
	if err != nil {
		panic("Error opening resources microservice file")
	}
	defer resourceFile.Close()

	var resources model.ArtifactResources
	byteResourceFile, _ := io.ReadAll(resourceFile)
	err = json.Unmarshal(byteResourceFile, &resources)
	if err != nil {
		panic("Error mapping resources microservice data")
	}
	return resources
}
