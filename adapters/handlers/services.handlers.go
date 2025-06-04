package handlers

import "github.com/medicue/core/services"

type HTTPHandler struct {
	service services.ServicesHandler
}

func HttpAdapter(serviceHandler *services.ServicesHandler) *HTTPHandler {
	return &HTTPHandler{
		service: *serviceHandler,
	}
}
