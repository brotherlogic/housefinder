package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/housefinder/proto"
	redfinlib "github.com/brotherlogic/redfinlib"
)

func (s *Server) processHouse(ctx context.Context, number int32) error {
	body, err := s.getter.get(ctx, fmt.Sprintf("https://www.redfin.com/CA/Albany/1141-Talbot-Ave-94706/home/%v", number))
	if err != nil {
		return err
	}

	stats, _ := redfinlib.Extract(body)
	s.Log(fmt.Sprintf("Got %+v", stats))

	housePrice := &pb.HousePrice{
		Id:       number,
		Listed:   stats.CurrentPrice,
		Estimate: stats.CurrentEstimate,
		Date:     time.Now().Unix()}

	s.config.LastRun = time.Now().Unix()

	if _, ok := s.config.FullHistory[number]; ok {
		s.config.FullHistory[number].History = append(s.config.FullHistory[number].History, housePrice)
	} else {
		s.config.FullHistory[number] = &pb.HouseHistory{Id: number, History: []*pb.HousePrice{housePrice}}
	}

	return s.save(ctx)
}
