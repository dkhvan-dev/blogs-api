package api

import (
	"blogs-api/internal/users/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
	}

	input := model.UserCreate{}

	if err := json.Unmarshal(data, &input); err != nil {
		handleError(w, err)
	}

	err = a.users.Service.Create(input)
	if err != nil {
		handleError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	response := []model.UserView{}
	w.Header().Add("Content-Type", "application/json")

	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	size, _ := strconv.Atoi(mux.Vars(r)["size"])

	response = a.users.Service.FindAll(page, size, "")

	data, err := json.Marshal(response)
	if err != nil {
		handleError(w, err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
