package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) processHouse(ctx context.Context, number int32) error {
	_, err := s.getter.get(ctx, fmt.Sprintf("https://www.redfin.com/CA/Albany/1141-Talbot-Ave-94706/home/%v", number))
	return err
}
