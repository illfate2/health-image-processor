package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/illfate2/health-image-processor/internal/body"
)

type Server struct {
	http.Handler
	service  *body.Processor
	upgrader *websocket.Upgrader
}

func NewServer(service *body.Processor) *Server {
	s := &Server{
		service:  service,
		upgrader: &websocket.Upgrader{},
	}
	s.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
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
	notifyCh := s.service.StartNotifying(ctx)
	for {
		select {
		case msg := <-notifyCh:
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
