package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/metastore"
	"github.com/denbeigh2000/jfsi/utils"

	"github.com/gorilla/mux"
)

type CreateRequest struct {
	NChunks int `json:"n"`
}

type HTTP struct {
	Store metastore.MetaStore

	mux *mux.Router
}

func NewHTTP(store metastore.MetaStore) http.Handler {
	router := mux.NewRouter()

	handler := &HTTP{
		Store: store,

		mux: router,
	}

	router.HandleFunc("/{id}", handler.HandleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", handler.HandleRetrieve).Methods(http.MethodGet)
	router.HandleFunc("/{id}", handler.HandleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{id}", handler.HandleDelete).Methods(http.MethodDelete)

	return handler
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *HTTP) HandleCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.RespondError(w, "No ID given", 400)
	}

	ID := jfsi.ID(id)

	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		txt := fmt.Sprintf("Couldn't parse JSON body (%v)", err.Error())
		utils.RespondError(w, txt, 400)
	}

	record, err := h.Store.Create(ID, req.NChunks)
	switch err.(type) {
	case nil:
		utils.RespondDataSuccess(w, record)
	case metastore.KeyAlreadyExistsErr, metastore.ZeroLengthCapacityRecordErr:
		utils.RespondError(w, err.Error(), 400)
	default:
		utils.RespondError(w, err.Error(), 500)
	}
}

func (h *HTTP) HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.RespondError(w, "No ID given", 400)
	}

	ID := jfsi.ID(id)

	record, err := h.Store.Retrieve(ID)
	switch err.(type) {
	case nil:
		utils.RespondDataSuccess(w, record)
	case metastore.KeyNotFoundErr:
		utils.RespondError(w, err.Error(), 404)
	default:
		utils.RespondError(w, err.Error(), 500)
	}
}

func (h *HTTP) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.RespondError(w, "No ID given", 400)
	}

	ID := jfsi.ID(id)

	var record metastore.Record
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		txt := fmt.Sprintf("Couldn't parse JSON body (%v)", err.Error())
		utils.RespondError(w, txt, 400)
	}

	err = h.Store.Update(ID, record)
	switch err.(type) {
	case nil:
		utils.RespondDataSuccess(w, record)
	case metastore.KeyNotFoundErr:
		utils.RespondError(w, err.Error(), 404)
	default:
		utils.RespondError(w, err.Error(), 500)
	}
}

func (h *HTTP) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.RespondError(w, "No ID given", 400)
	}

	ID := jfsi.ID(id)

	err := h.Store.Delete(ID)
	switch err.(type) {
	case nil:
		utils.RespondSuccess(w, "OK")
	case metastore.KeyNotFoundErr:
		utils.RespondError(w, err.Error(), 404)
	default:
		utils.RespondError(w, err.Error(), 500)
	}
}
