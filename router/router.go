package router

import (
	"github.com/gorilla/mux"
	"github.com/yusologia/go-core/handler"
	"github.com/yusologia/go-core/middleware"
)

type CallbackRouter func(*mux.Router)

func RegisterRouter(router *mux.Router, callback CallbackRouter) {
	router.Use(middleware.PanicHandler)
	router.Use(middleware.PrepareRequestHandler)

	// Storage route
	storageHnd := handler.StorageHandler{}
	router.HandleFunc("/storages/{path:.*}", storageHnd.ShowFile).Methods("GET")

	callback(router)
}
