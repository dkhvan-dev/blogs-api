package api

import (
	"blogs-api/internal/users/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

const (
	APPLICATION_JSON = "application/json"
	CONTENT_TYPE     = "Content-type"
)

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(CONTENT_TYPE, APPLICATION_JSON)

	data, err := io.ReadAll(r.Body)

	if err != nil {
		handleError(w, err)
	}

	input := model.UserLoginRequest{}

	err = json.Unmarshal(data, &input)
	if err != nil {
		handleError(w, err)
	}

	tokensMap, err := a.users.Service.Login(input)
	if err != nil {
		handleError(w, err)
	}

	response, err := json.Marshal(tokensMap)
	if err != nil {
		handleError(w, err)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(CONTENT_TYPE, APPLICATION_JSON)

	data, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
	}

	input := model.UserCreate{}

	if err = json.Unmarshal(data, &input); err != nil {
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
	w.Header().Add(CONTENT_TYPE, APPLICATION_JSON)

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

func (a *API) FindUserByLogin(w http.ResponseWriter, r *http.Request) {
	response := &model.UserView{}
	w.Header().Add(CONTENT_TYPE, APPLICATION_JSON)

	login := mux.Vars(r)["login"]
	response, err := a.users.Service.FindByLogin(login)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(response)
	if err != nil {
		handleError(w, err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	input := model.UserUpdate{}
	login := mux.Vars(r)["login"]
	data, err := io.ReadAll(r.Body)

	if err != nil {
		handleError(w, err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		handleError(w, err)
	}

	if err = a.users.Service.Update(login, input); err != nil {
		handleError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) ChangePassword(w http.ResponseWriter, r *http.Request) {
	input := model.UserChangePassword{}
	login := mux.Vars(r)["login"]

	data, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		handleError(w, err)
	}

	if err = a.users.Service.ChangePassword(login, input); err != nil {
		handleError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	login := mux.Vars(r)["login"]

	if err := a.users.Service.Delete(login); err != nil {
		handleError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
