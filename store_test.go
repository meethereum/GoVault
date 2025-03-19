package main

import (
	"bytes"
	"os"
	"testing"
)

func TestStore(t *testing.T){
	opts:=StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s:=NewStore(opts)

	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpictures1",data);err!=nil{
		t.Error(err)
	}
	if err:=os.RemoveAll("myspecialpictures1");err!=nil{
		t.Error(err)
	}
}