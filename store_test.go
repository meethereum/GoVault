package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestPathTransformFunc(t *testing.T){
	key:="mybestpicture"
	pathname:=CASPathTransformFunc(key)
	fmt.Println(pathname)
	expectedPathName := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"
	if pathname!=expectedPathName{
		t.Errorf("have %s and want %s",pathname,expectedPathName)
	}
}

func TestStore(t *testing.T){
	opts:=StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
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