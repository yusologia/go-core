package handler

import (
	"github.com/gorilla/mux"
	"github.com/yusologia/go-core/filesystem"
	"net/http"
)

type StorageHandler struct{}

func (ctr StorageHandler) ShowFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storage := filesystem.Storage{IsPublic: true}
	storage.ShowFile(w, r, vars["path"])
}
