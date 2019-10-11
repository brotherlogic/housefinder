package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	redfinlib "github.com/brotherlogic/redfinlib"
)

func (s *Server) processHouse(ctx context.Context, number int32) error {
	body, err := s.getter.get(ctx, fmt.Sprintf("https://www.redfin.com/CA/Albany/1141-Talbot-Ave-94706/home/%v", number))
	if err != nil {
		return err
	}

	stats, _ := redfinlib.Extract(body)
	s.Log(fmt.Sprintf("Got %+v", stats))

	s.config.LastRun = time.Now().Unix()
	return s.save(ctx)
}
