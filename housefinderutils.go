package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/housefinder/proto"
	redfinlib "github.com/brotherlogic/redfinlib"
	pbrf "github.com/brotherlogic/redfinlib/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	//Backlog - the print queue
	price = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "housefinder_price",
		Help: "The size of the print queue",
	}, []string{"house"})
)

func (s *Server) processHouse(ctx context.Context, number int32) error {
	// Don't process a sold house.
	if val, ok := s.config.FullHistory[number]; ok {
		for _, hist := range val.History {
			if hist.Sold {
				return nil
			}
		}
	}

	body, err := s.getter.get(ctx, fmt.Sprintf("https://www.redfin.com/CA/Albany/1141-Talbot-Ave-94706/home/%v", number))
	if err != nil {
		return err
	}

	stats, _ := redfinlib.Extract(body)
	s.Log(fmt.Sprintf("Got %+v", stats))

	price.With(prometheus.Labels{"house": fmt.Sprintf("%v", number)}).Set(float64(stats.CurrentEstimate))
	housePrice := &pb.HousePrice{
		Id:       number,
		Listed:   stats.CurrentPrice,
		Estimate: stats.CurrentEstimate,
		Date:     time.Now().Unix(),
		Sold:     stats.State == pbrf.Stats_SOLD}

	s.config.LastRun = time.Now().Unix()

	if _, ok := s.config.FullHistory[number]; ok {
		s.config.FullHistory[number].History = append(s.config.FullHistory[number].History, housePrice)
	} else {
		s.config.FullHistory[number] = &pb.HouseHistory{Id: number, History: []*pb.HousePrice{housePrice}}
	}

	return s.save(ctx)
}
