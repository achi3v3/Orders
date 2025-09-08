package router

import (
	"fmt"
	"net/http"
	"orders/internal/config"
	"orders/internal/subs"

	"github.com/sirupsen/logrus"
)

type Server struct {
	handler *subs.Handler
	logger  *logrus.Logger
}

func NewServer(handler *subs.Handler, logger *logrus.Logger) *Server {
	return &Server{
		handler: handler,
		logger:  logger,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/order/{order_uid}", s.handler.GetOrderFromHttp)
	port := config.GetEnv("PORT", "8081")
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	if err := server.ListenAndServe(); err != nil {
		s.logger.Errorf("Server.Run: error with listen server %v", err)
	}

	s.logger.Infof("Server.Run: Server UP: http://localhost:8081/order")
}
