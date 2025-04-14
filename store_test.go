package main

import (
	"bytes"
	// "os"
	"testing"
	"io"
	"fmt"
)

func TestPathTransformFunc(t *testing.T){
	key:="mybestpicture"
	pathKey:=CASPathTransformFunc(key)
	expectedPathName := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	expectedOriginalKey:= "be17b32c2870b1c0c73b59949db6a3be7814dd23"
	if pathKey.Pathname!=expectedPathName{
		t.Errorf("have %s and want %s",pathKey.Pathname,expectedPathName)
	}
	if pathKey.Original!=expectedOriginalKey{
		t.Errorf("have %s and want %s",pathKey.Original,expectedOriginalKey)
	}
	
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momsspecials"
	data := []byte("some jpg bytes")

	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}
}