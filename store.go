package main

import (
	"io"
	"log"
	"os"
)

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc PathTransformFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{opts}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.Mkdir(pathName, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"
	pathAndFilename := pathName + "/" + filename
	_, err := os.Stat(pathAndFilename)
	if err == nil {
		return os.ErrExist
	}

	

	f, err := os.Create(pathName + "/" + filename)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d bytes) to %s", n, pathAndFilename)
	return nil
}
