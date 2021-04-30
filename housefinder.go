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
	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/housefinder/proto"
)

const (
	// KEY - where the config is stored
	KEY = "/github.com/brotherlogic/housefinder/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	houses []int32
	getter
	config *pb.Config
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		houses:   []int32{int32(1054537)},
	}
	s.getter = &prodGetter{s.HTTPGet}
	s.config = &pb.Config{}
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
	if master {
		return s.load(ctx)
	}

	return nil
}

func (s *Server) save(ctx context.Context) error {
	return s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)
	if err != nil {
		return err
	}
	config = data.(*pb.Config)

	if config.FullHistory == nil {
		config.FullHistory = make(map[int32]*pb.HouseHistory)
	}

	s.config = config

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "tracked", Value: int64(len(s.config.FullHistory))},
		&pbg.State{Key: "last_run", Value: s.config.LastRun},
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
	var init = flag.Bool("init", false, "Prep server")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("housefinder", false, true)
	if err != nil {
		return
	}

	if *init {
		ctx, cancel := utils.BuildContext("housefinder", "housefinder")
		defer cancel()

		server.config.LastRun = time.Now().Unix()
		err := server.save(ctx)
		fmt.Printf("%v\n", err)
		return
	}

	go func() {
		for true {
			ctx, cancel := utils.BuildContext("hf", "hf")
			server.processHouses(ctx)
			cancel()

			time.Sleep(time.Minute * 30)
		}
	}()

	fmt.Printf("%v", server.Serve())
}
