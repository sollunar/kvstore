package kvstore

import (
	"net/http"

	"github.com/sollunar/kvstore-api/pkg/req"
	"github.com/sollunar/kvstore-api/pkg/res"
)

type Handler struct {
	kvservice *KVService
}

func NewHandler(router *http.ServeMux, kvservice *KVService) {
	h := &Handler{
		kvservice: kvservice,
	}

	router.HandleFunc("/get", h.Get)
	router.HandleFunc("/set", h.Set)
	router.HandleFunc("/delete", h.Delete)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		res.Error(w, "key not found", http.StatusBadRequest)
		return
	}

	val, err := h.kvservice.Get(key)
	if err != nil {
		if err == ErrKeyNotFound {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
		res.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res.Json(w, GetResponse{key, val}, http.StatusOK)
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[SetRequest](&w, r)
	if err != nil {
		return
	}

	if body.Key == "" || body.Value == "" {
		res.Error(w, "missing key or value", http.StatusBadRequest)
		return
	}

	if err := h.kvservice.Set(body.Key, body.Value); err != nil {
		res.Error(w, "failed to set", http.StatusInternalServerError)
		return
	}
	res.Json(w, nil, http.StatusCreated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		res.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	err := h.kvservice.Delete(key)
	if err != nil {
		if err == ErrKeyNotFound {
			res.Error(w, "key not found", http.StatusNotFound)
			return
		}
		res.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res.Json(w, nil, http.StatusNoContent)
}
