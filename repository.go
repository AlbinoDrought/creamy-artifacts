package main

import (
	"io"
	"os"
	"path"
)

type ArtifactRepository interface {
	List() ([]string, error)
	Store(key string, reader io.Reader) error
	Length(key string) (int64, error)
	Pull(key string, writer io.Writer) error
	Remove(key string) error
}

type localArtifactRepository struct {
	basePath string
}

func (repository *localArtifactRepository) localPath(key string) string {
	return path.Join(repository.basePath, path.Clean(key))
}

func (repository *localArtifactRepository) List() ([]string, error) {
	file, err := os.Open(repository.basePath)
	if err != nil {
		return nil, err
	}
	names, err := file.Readdirnames(0)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	return names, err
}

func (repository *localArtifactRepository) Store(key string, reader io.Reader) error {
	file, err := os.OpenFile(repository.localPath(key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	return err
}

func (repository *localArtifactRepository) Length(key string) (int64, error) {
	file, err := os.OpenFile(repository.localPath(key), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return 0, err
	}

	size := int64(0)
	info, err := file.Stat()
	if err1 := file.Close(); err == nil {
		size = info.Size()
		err = err1
	}
	return size, err
}

func (repository *localArtifactRepository) Pull(key string, writer io.Writer) error {
	file, err := os.OpenFile(repository.localPath(key), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	return err
}

func (repository *localArtifactRepository) Remove(key string) error {
	return os.Remove(repository.localPath(key))
}

type Project struct {
	artifactRepository ArtifactRepository
}

func (project *Project) ListArtifacts() ([]string, error) {
	return project.artifactRepository.List()
}

func (project *Project) StoreArtifact(key string, reader io.Reader) error {
	return project.artifactRepository.Store(key, reader)
}

func (project *Project) MeasureArtifacts(keys []string) (int64, error) {
	size := int64(0)
	for _, key := range keys {
		artifactSize, err := project.artifactRepository.Length(key)
		if err != nil {
			return 0, err
		}
		size += artifactSize
	}
	return size, nil
}

func (project *Project) CollateArtifacts(keys []string, writer io.Writer) error {
	for _, key := range keys {
		if err := project.artifactRepository.Pull(key, writer); err != nil {
			return err
		}
	}
	return nil
}

func (project *Project) RemoveArtifact(key string) error {
	return project.artifactRepository.Remove(key)
}
