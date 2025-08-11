package handlers

import "github.com/diagnoxix/core/services"

type HTTPHandler struct {
	service services.ServicesHandler
}

func HttpAdapter(serviceHandler *services.ServicesHandler) *HTTPHandler {
	return &HTTPHandler{
		service: *serviceHandler,
	}
}
