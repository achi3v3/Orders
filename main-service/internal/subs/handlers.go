package subs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orders/internal/models"
	"strings"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	service *Service
	logger  *logrus.Logger
}

func NewHandler(service *Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}
func (h *Handler) Create(ctx context.Context, jsonOrder *models.OrderJson) error {
	return h.service.Create(ctx, jsonOrder)

}

func (h *Handler) GetOrderFromHttp(w http.ResponseWriter, r *http.Request) {

	if strings.ToLower(r.Method) != "get" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderUID := getParamFromPath(r.URL.Path)
	if orderUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(context.Background(), orderUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Handler.GetOrderFromHttp: %v", err)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Handler.GetOrderFromHttp: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	fmt.Println(order)
}

func (h *Handler) GetOrder(orderUID string) {
	order, err := h.service.GetOrder(context.Background(), orderUID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(order)
}

func getParamFromPath(path string) string {
	param := path[strings.LastIndex(path, "/")+1:]
	return param
}
