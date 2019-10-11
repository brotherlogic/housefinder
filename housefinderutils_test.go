package main

import (
	"fmt"
	"testing"

	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
)

type testGetter struct {
	fail    bool
	httpGet func(ctx context.Context, url string) (string, error)
}

func (t *testGetter) get(ctx context.Context, url string) (string, error) {
	if t.fail {
		return "", fmt.Errorf("Built to fail")
	}
	return "got", nil
}

func InitTest() *Server {
	s := Init()
	s.getter = &testGetter{}
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	return s
}

func TestProcess(t *testing.T) {
	s := InitTest()

	err := s.processHouse(context.Background(), int32(123))

	if err != nil {
		t.Errorf("Error processing house")
	}
}

func TestProcessReadFail(t *testing.T) {
	s := InitTest()
	s.getter = &testGetter{fail: true}

	err := s.processHouse(context.Background(), int32(123))

	if err == nil {
		t.Errorf("No error on read fail")
	}
}
