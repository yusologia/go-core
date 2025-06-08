package logiacore

import (
	"github.com/gorilla/mux"
	"github.com/yusologia/go-core/v2/handler"
	logiamdw "github.com/yusologia/go-core/v2/middleware"
)

type CallbackRouter func(*mux.Router)

func RegisterRouter(router *mux.Router, callbacks ...CallbackRouter) {
	router.Use(logiamdw.PanicHandler)
	router.Use(logiamdw.PrepareRequestHandler)

	h := handler.Handler{}
	router.HandleFunc("/health-check", h.HealthCheck).Methods("GET")
	router.HandleFunc("/storages/{path:.*}", h.StorageShowFile).Methods("GET")
	router.HandleFunc("/{path:.*}/log-active", h.LogActivate).Methods("POST")

	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			callback(router)
		}
	}
}
