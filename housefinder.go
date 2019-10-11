package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	houses []int32
	getter
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		houses:   []int32{int32(17187297)},
	}
	s.getter = &prodGetter{s.HTTPGet}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	//Pass
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "what", Value: int64(123)},
	}
}

type getter interface {
	get(ctx context.Context, url string) (string, error)
}

type prodGetter struct {
	httpGet func(ctx context.Context, url string, header string) (string, error)
}

func (p *prodGetter) get(ctx context.Context, url string) (string, error) {
	return p.httpGet(ctx, url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36")
}

func (s *Server) processHouses(ctx context.Context) error {
	for _, house := range s.houses {
		err := s.processHouse(ctx, house)
		if err != nil {
			return err
		}
		time.Sleep(time.Minute)
	}
	return nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	server.RPCTracing = true

	err := server.RegisterServer("housefinder", false)
	if err != nil {
		log.Fatalf("Registration Error: %v", err)
	}

	server.RegisterRepeatingTask(server.processHouses, "process_houses", time.Minute*10)

	fmt.Printf("%v", server.Serve())
}
