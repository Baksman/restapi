package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/baksman/rest-api/api/auth"
	"github.com/baksman/rest-api/api/formaterror"
	"github.com/baksman/rest-api/api/models"
	"github.com/baksman/rest-api/api/responses"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}

	json.Unmarshal(body, &user)

	user.Prepare()
	validate := validator.New()
	err = validate.Struct(user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignUp(&user)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"token":    token,
		"email":    user.Email,
		"password": user.Password,
		"nickname": user.Nickname,
		"created":  user.CreatedAt,
		"id":       user.ID,
	})

}

func (server *Server) SignUp(user *models.User) (string, error) {
	err := server.DB.Debug().Model(&user).Error

	if err != nil {
		return "", err

	} else {
		err = server.DB.Save(user).Error
		if err != nil {
			// responses.ERROR(w, http.StatusInternalServerError, err)
			return "", err
		} else {
			return auth.CreateToken(user.ID)
		}
	}
}
