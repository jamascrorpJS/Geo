package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"jamascrorpJS/gwatch/internal/interactors"
	"jamascrorpJS/gwatch/internal/models"
)

type Router interface {
	Save(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type routes struct {
	mux         *http.ServeMux
	coordinator interactors.Coordinator
}

func NewRoutes(mux *http.ServeMux, coordinator interactors.Coordinator) Router {
	return &routes{
		mux:         mux,
		coordinator: coordinator,
	}
}

func (routes *routes) Save(w http.ResponseWriter, r *http.Request) {
	var coordinate models.Position
	var response models.Response

	if err := json.NewDecoder(r.Body).Decode(&coordinate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err := routes.coordinator.SendCoordinate(r.Context(), coordinate)
	if err != nil {
		if errors.Is(err, models.ErrInternal) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response.Mes = "Created"
	response.Code = http.StatusCreated
	response.Data = coordinate
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func (routes *routes) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("deviceId")
	var response models.Response
	res, err := routes.coordinator.ReceiveCoordinate(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrInternal) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	response.Mes = "Created"
	response.Code = http.StatusOK
	response.Data = res
	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
