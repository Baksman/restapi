package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/baksman/rest-api/api/auth"
	"github.com/baksman/rest-api/api/formaterror"
	"github.com/baksman/rest-api/api/models"
	"github.com/baksman/rest-api/api/responses"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = validate.Struct(user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

func (server *Server) Updateuser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	validate := validator.New()
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user.Prepare()
	err = validate.Struct(user)
	//  err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusBadRequest, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

func (s *Server) uploadProfilePic(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myfile")

	if err != nil || file == nil {
		responses.ERROR(w, http.StatusConflict, err)
	}

	defer file.Close()

	dst, err := os.Create(handler.Filename)

	if err != nil {
		responses.ERROR(w, http.StatusConflict, err)
	}

	if _, err = io.Copy(dst, file); err != nil {
		responses.ERROR(w, http.StatusConflict, err)
	}

	responses.JSON(w, http.StatusOK,"image uploaded successfully")

}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["uid"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	tokenID, err := auth.ExtractTokenID(r)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorised"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")

}
