package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/illfate2/health-image-processor/internal/blinker"
	"github.com/illfate2/health-image-processor/internal/image"
)

type Server struct {
	http.Handler
	service   *blinker.Service
	processor *image.Processor
	upgrader  *websocket.Upgrader
}

type notifyMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func newNotifyMessage(t string, message string) notifyMessage {
	return notifyMessage{Type: t, Message: message}
}

func NewServer(service *blinker.Service, processor *image.Processor) *Server {
	s := &Server{
		service:   service,
		upgrader:  &websocket.Upgrader{},
		processor: processor,
	}
	engine := gin.New()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	engine.GET("/events", s.handleEvents)
	s.Handler = engine
	return s
}

func (s *Server) handleEvents(c *gin.Context) {
	conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	ctx, cancelF := context.WithCancel(c.Request.Context())
	defer cancelF()
	notifyCh := s.service.StartNotifyingToBlink(ctx)
	imageCh := s.processor.GetNotificationCh()
	for {
		select {
		case <-notifyCh:
			err = conn.WriteJSON(newNotifyMessage("blink", "Pls, blink"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-imageCh:
			err = conn.WriteJSON(newNotifyMessage("crooked", "Pls, sit up straight"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
