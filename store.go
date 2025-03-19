package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) string {
	hash:=sha1.Sum([]byte(key))
	hashStr:= hex.EncodeToString(hash[:]) //convert byte to slice
	blocksize:=5
	sliceLen:=len(hashStr)/blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from,to:=i*blocksize,(i*blocksize)+blocksize 	
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths,"/")
	

}

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
