package body

import (
	"context"
	"log"
	"time"
)

type Processor struct {
	userLastBlinked   time.Time
	lastBlinkNotified time.Time
	blinkedTime       int
	blinkTimeout      int

	noseCrookedCh           chan struct{}
	lastCrookedNoseNotified time.Time
	noseCrookedTime         int
	noseCrookedTimeout      int

	backCrookedCh           chan struct{}
	lastCrookedBackNotified time.Time
	backCrookedTime         int
	backCrookedTimeout      int
}

func NewService(blinkedTime int) *Processor {
	return &Processor{
		blinkedTime:   blinkedTime,
		blinkTimeout:  5,
		noseCrookedCh: make(chan struct{}, 100000),
		backCrookedCh: make(chan struct{}, 100000),
	}
}

type NotifyMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func newNotifyMessage(t string, message string) NotifyMessage {
	return NotifyMessage{Type: t, Message: message}
}

func (s *Processor) StartNotifying(ctx context.Context) <-chan NotifyMessage {
	const ticketDur = time.Millisecond * 50
	blinkTicker := time.NewTicker(ticketDur)

	ch := make(chan NotifyMessage)
	go func() {
		for {
			select {
			case <-blinkTicker.C:
				if time.Since(s.userLastBlinked).Seconds() > float64(s.blinkedTime) && time.Since(s.lastBlinkNotified).Seconds() > float64(s.blinkTimeout) {
					log.Print("Sending notification to blink")
					ch <- NotifyMessage{
						Type:    "blink",
						Message: "Pls, blink",
					}
					s.lastBlinkNotified = time.Now()
				}

			case <-s.backCrookedCh:
				ch <- NotifyMessage{
					Type:    "crookedBack",
					Message: "Pls, sit straight",
				}
			case <-s.noseCrookedCh:
				ch <- NotifyMessage{
					Type:    "crookedHead",
					Message: "Pls, set your head straight",
				}
			case <-ctx.Done():
				blinkTicker.Stop()
				close(ch)
				return
			}
		}
	}()
	return ch
}

func (s *Processor) Blinked(isFaceRecognized bool) {
	if !isFaceRecognized {
		log.Print("face is not recognized for blinking")
		return
	}
	log.Print("Blinking...")
	s.userLastBlinked = time.Now()
}

func (s *Processor) NoseCrooked(isCrooked, isFaceRecognized bool) {
	if !isFaceRecognized {
		log.Print("face is not recognized for nose")
		return
	}
	log.Print("Nose crooked...", isCrooked)
	if isCrooked {
		s.noseCrookedCh <- struct{}{}
	}
}

func (s *Processor) BackCrooked(isCrooked, isFaceRecognized bool) {
	if !isFaceRecognized {
		log.Print("face is not recognized for back")
		return
	}
	log.Print("Back crooked...", isCrooked)
	if isCrooked {
		s.backCrookedCh <- struct{}{}
	}
}
