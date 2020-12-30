package handlers

import (
	"net/http"

	"github.com/texazcowboy/warehouse/internal/user"

	"github.com/go-playground/validator/v10"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type UserHandler struct {
	*BaseHandler
}

func NewUserHandler(base *BaseHandler) *UserHandler {
	return &UserHandler{base}
}

func (e *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := web.Decode(r, &u)
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusBadRequest, err)
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			e.LogEntry.Error(err)
		}
	}()

	if err = e.Validate.Struct(&u); err != nil {
		e.LogEntry.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			e.renderError(w, http.StatusInternalServerError, err)
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.renderError(w, http.StatusBadRequest, err)
		}
		return
	}

	err = user.Create(&u, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusInternalServerError, err)
		return
	}
	u.Password = "***"
	e.renderSuccess(w, http.StatusCreated, &u)
}

func (e *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusBadRequest, err)
		return
	}

	if err = user.Delete(int64(id), e.DB); err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusInternalServerError, err)
		return
	}

	e.renderSuccess(w, http.StatusNoContent, nil)
}
