package handler

import (
	"io"
	"net/http"

	"github.com/denbeigh2000/jfsi"
	"github.com/denbeigh2000/jfsi/application"
	// TODO: Wrap these errors in application errors
	"github.com/denbeigh2000/jfsi/storage"
	"github.com/gorilla/mux"
)

type HTTP struct {
	Node application.Node

	mux *mux.Router
}

func NewHTTP(node application.Node) http.Handler {
	router := mux.NewRouter()

	handler := &HTTP{
		Node: node,

		mux: router,
	}

	router.HandleFunc("/", handler.HandleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{id}", handler.HandleRetrieve).Methods(http.MethodGet)
	router.HandleFunc("/{id}", handler.HandleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{id}", handler.HandleDelete).Methods(http.MethodDelete)

	return handler
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *HTTP) HandleCreate(w http.ResponseWriter, r *http.Request) {
	id, err := h.Node.Create(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	http.Error(w, id.String(), 200)
}

func (h *HTTP) HandleRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID given", 400)
	}

	ID, err := jfsi.IDFromString(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	data, err := h.Node.Retrieve(ID)
	switch err.(type) {
	case nil:
		io.Copy(w, data)
	case storage.NotFoundErr:
		http.Error(w, err.Error(), 404)
	default:
		http.Error(w, err.Error(), 500)
	}
}

func (h *HTTP) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID given", 400)
	}

	ID, err := jfsi.IDFromString(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.Node.Update(ID, r.Body)
	switch err.(type) {
	case nil:
		http.Error(w, "OK", 200)
	case storage.NotFoundErr:
		http.Error(w, err.Error(), 404)
	default:
		http.Error(w, err.Error(), 500)
	}
}

func (h *HTTP) HandleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No ID given", 400)
	}

	ID, err := jfsi.IDFromString(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.Node.Delete(ID)
	switch err.(type) {
	case nil:
		http.Error(w, "OK", 200)
	case storage.NotFoundErr:
		http.Error(w, err.Error(), 404)
	default:
		http.Error(w, err.Error(), 500)
	}
}
