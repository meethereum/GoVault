package main

import (
	"bytes"
	//"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) Pathkey {
	hash:=sha1.Sum([]byte(key))
	hashStr:= hex.EncodeToString(hash[:]) //convert byte to slice
	blocksize:=5
	sliceLen:=len(hashStr)/blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from,to:=i*blocksize,(i*blocksize)+blocksize 	
		paths[i] = hashStr[from:to]
	}

	return Pathkey{
		Pathname: strings.Join(paths,"/"),
		Original:hashStr,
	}


}

type Pathkey struct{
	Pathname string
	Original string
}

type PathTransformFunc func(string) Pathkey  

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}


func (p Pathkey) Filename() string{
	return fmt.Sprintf("%s/%s",p.Pathname,p.Original); 
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{opts}
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser,error){
	pathkKey:=s.PathTransformFunc(key)

	f,err:=os.Open(pathkKey.Filename())
	if err!=nil {
		return nil,err
	}

	return f,err

}


func (s *Store) WriteStream(key string, r io.Reader) error {
	
	pathKey := s.PathTransformFunc(key)
	
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	 pathAndFilename := pathKey.Filename() 

	f, err := os.Create(pathAndFilename)
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
