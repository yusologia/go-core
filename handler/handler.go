package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	logiafs "github.com/yusologia/go-core/v2/filesystem"
	logiapkg "github.com/yusologia/go-core/v2/pkg"
	logiares "github.com/yusologia/go-core/v2/response"
	"net/http"
)

type Handler struct{}

func (Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func (Handler) StorageShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storage := logiafs.Storage{IsPublic: true}
	storage.ShowFile(w, r, vars["path"])
}

func (Handler) LogActivate(w http.ResponseWriter, r *http.Request) {
	logiapkg.LOG_ACTIVE = !logiapkg.LOG_ACTIVE
	status := "inactive"
	if logiapkg.LOG_ACTIVE {
		status = "active"
	}

	res := logiares.Response{Object: map[string]interface{}{"log": status}}

	res.Success(w)
}
