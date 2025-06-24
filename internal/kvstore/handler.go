package kvstore

import (
	"errors"
	"net/http"

	"github.com/sollunar/kvstore-api/pkg/req"
	"github.com/sollunar/kvstore-api/pkg/res"
)

type KVStoreHandler struct {
	kvservice *KVService
}

func NewKVStoreHandler(router *http.ServeMux, kvservice *KVService) {
	h := &KVStoreHandler{
		kvservice: kvservice,
	}

	router.HandleFunc("/get", h.Get)
	router.HandleFunc("/set", h.Set)
	router.HandleFunc("/delete", h.Delete)
}

// Get godoc
// @Summary Get value by key
// @Description Retrieve the value for a given key
// @Tags kvstore
// @Accept  json
// @Produce  json
// @Param   key query string true "Key"
// @Success 200 {object} kvstore.GetResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Key not found"
// @Failure 500 {string} string "Internal error"
// @Router /get [get]
func (h *KVStoreHandler) Get(w http.ResponseWriter, r *http.Request) {
	key, err := req.GetQueryParam(r, "key")
	if err != nil {
		res.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	val, err := h.kvservice.Get(key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
		res.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res.Json(w, GetResponse{key, val}, http.StatusOK)
}

// Set godoc
// @Summary Set a key-value pair
// @Description Store a key-value entry in the store
// @Tags kvstore
// @Accept  json
// @Produce  json
// @Param   data body kvstore.SetRequest true "Key-Value Pair"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal error"
// @Router /set [post]
func (h *KVStoreHandler) Set(w http.ResponseWriter, r *http.Request) {
	body, err := req.HandleBody[SetRequest](&w, r)
	if err != nil {
		return
	}
	defer r.Body.Close()

	if err := h.kvservice.Set(body.Key, body.Value); err != nil {
		res.Error(w, "failed to set", http.StatusInternalServerError)
		return
	}

	res.Json(w, nil, http.StatusCreated)
}

// Delete godoc
// @Summary Delete key
// @Description Delete a key-value entry by key
// @Tags kvstore
// @Accept  json
// @Produce  json
// @Param   key query string true "Key"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Key not found"
// @Failure 500 {string} string "Internal error"
// @Router /delete [delete]
func (h *KVStoreHandler) Delete(w http.ResponseWriter, r *http.Request) {
	key, err := req.GetQueryParam(r, "key")
	if err != nil {
		res.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.kvservice.Delete(key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			res.Error(w, "key not found", http.StatusNotFound)
			return
		}
		res.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res.Json(w, nil, http.StatusNoContent)
}
