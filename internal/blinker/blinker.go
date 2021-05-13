package blinker

import (
	"context"
	"log"
	"time"
)

type Service struct {
	userLastBlinked time.Time
	lastNotified    time.Time
	blinkedTime     int
	timeout         int
}

func NewService(blinkedTime int) *Service {
	return &Service{
		blinkedTime: blinkedTime,
		timeout:     5,
	}
}

func (s *Service) StartNotifyingToBlink(ctx context.Context) <-chan struct{} {
	const ticketDur = time.Millisecond * 50
	ticker := time.NewTicker(ticketDur)

	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Since(s.userLastBlinked).Seconds() > float64(s.blinkedTime) && time.Since(s.lastNotified).Seconds() > float64(s.timeout) {
					log.Print("Sending notification to blink")
					ch <- struct{}{}
					s.lastNotified = time.Now()
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
