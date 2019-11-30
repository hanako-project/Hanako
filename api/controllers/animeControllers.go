package controllers

import (
	"github.com/gorilla/mux"
	"github.com/luqmansen/hanako/models/models-postegres"
	u "github.com/luqmansen/hanako/utils"
	"net/http"
	"strconv"
)

var GetAll = func(w http.ResponseWriter, r *http.Request) {

	show := r.URL.Query().Get("show")

	data := models_postegres.GetAll(show)
	if data == nil {
		u.Respond(w, u.Message(http.StatusNoContent, "Not Found"))
	} else {
		resp := u.Message(http.StatusOK, "Success")
		resp["data"] = data
		u.Respond(w, resp)
	}
}

var GetByTitle = func(w http.ResponseWriter, r *http.Request) {

	keyword := r.URL.Query().Get("title")
	data := models_postegres.GetByTitle(keyword)
	if data == nil {
		u.Respond(w, u.Message(http.StatusNoContent, "Not Found"))
	} else {
		resp := u.Message(http.StatusOK, "Success")
		resp["data"] = data
		u.Respond(w, resp)
	}
}

var GetById = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		u.Respond(w, u.Message(http.StatusBadRequest, "There was an error in your request"))
	} else if data := models_postegres.GetByID(uint(id)); data == nil {
		u.Respond(w, u.Message(http.StatusNoContent, "Not Found"))
	} else {

		resp := u.Message(http.StatusOK, "Success")
		resp["data"] = data
		u.Respond(w, resp)
	}
}