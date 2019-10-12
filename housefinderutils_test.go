package main

import (
	"fmt"
	"testing"

	pb "github.com/brotherlogic/housefinder/proto"
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
	s.config.FullHistory = make(map[int32]*pb.HouseHistory)
	return s
}

func TestProcess(t *testing.T) {
	s := InitTest()

	err := s.processHouse(context.Background(), int32(123))

	if err != nil {
		t.Errorf("Error processing house")
	}

	err = s.processHouse(context.Background(), int32(123))

	if err != nil {
		t.Errorf("Error on double process")
	}
}

func TestProcessDone(t *testing.T) {
	s := InitTest()
	s.config.FullHistory[int32(123)] = &pb.HouseHistory{History: []*pb.HousePrice{&pb.HousePrice{Sold: true}}}

	err := s.processHouse(context.Background(), int32(123))

	if err != nil {
		t.Errorf("Error processing house")
	}

	if len(s.config.FullHistory[int32(123)].History) != 1 {
		t.Errorf("History added")
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
