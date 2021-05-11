package blinker

import (
	"context"
	"log"
	"time"
)

type Service struct {
	userLastBlinked time.Time
	blinkedTime     int
}

func NewService(blinkedTime int) *Service {
	return &Service{blinkedTime: blinkedTime}
}

func (s *Service) StartNotifyingToBlink(ctx context.Context) <-chan struct{} {
	const ticketDur = time.Millisecond * 50
	ticker := time.NewTicker(ticketDur)

	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Since(s.userLastBlinked).Seconds() > float64(s.blinkedTime) {
					log.Print("Sending notification to blink")
					ch <- struct{}{}
				}
			case <-ctx.Done():
				ticker.Stop()
				close(ch)
				return
			}
		}
	}()
	return ch
}

func (s *Service) Blinked() {
	log.Print("Blinking...")
	s.userLastBlinked = time.Now()
}
