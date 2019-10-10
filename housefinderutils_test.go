package main

import (
	"testing"

	"golang.org/x/net/context"
)

type testGetter struct {
	httpGet func(ctx context.Context, url string) (string, error)
}

func (t *testGetter) get(ctx context.Context, url string) (string, error) {
	return "got", nil
}

func InitTest() *Server {
	s := Init()
	s.getter = &testGetter{}
	return s
}

func TestProcess(t *testing.T) {
	s := InitTest()

	err := s.processHouse(context.Background(), int32(123))

	if err != nil {
		t.Errorf("Error processing house")
	}
}
