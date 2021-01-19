package handlers

import (
	"net/http"
	"time"

	"github.com/texazcowboy/warehouse/cmd/warehouse-api/common"

	"github.com/texazcowboy/warehouse/internal/foundation/security"

	"github.com/texazcowboy/warehouse/internal/foundation/crypto"

	"github.com/texazcowboy/warehouse/internal/user"

	"github.com/go-playground/validator/v10"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type UserHandler struct {
	*common.BaseHandler
	service user.ServiceInterface
}

func NewUserHandler(base *common.BaseHandler, service user.ServiceInterface) *UserHandler {
	return &UserHandler{base, service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := web.Decode(r, &u)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to decode request body")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			h.LogEntry.WithError(err).Error("Unable to close request body")
		}
	}()

	if err = h.Validate.Struct(&u); err != nil {
		h.LogEntry.WithError(err).Error("Invalid request body ")
		if _, ok := err.(*validator.InvalidValidationError); ok {
			h.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			h.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	err = h.service.CreateUser(&u)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to create user")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.Password = "***"
	h.RenderSuccess(w, http.StatusCreated, &u)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := web.Decode(r, &u)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to decode request body")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			h.LogEntry.WithError(err).Error("Unable to close request body")
		}
	}()

	if err = h.Validate.Struct(&u); err != nil {
		h.LogEntry.WithError(err).Error("Invalid request body ")
		if _, ok := err.(*validator.InvalidValidationError); ok {
			h.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			h.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	savedUser, err := h.service.GetUserByUsername(u.Username)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to get user by username")
		h.RenderError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if !crypto.Equals(u.Password, savedUser.Password) {
		h.RenderError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	generatedToken, err := security.GenerateToken(map[string]interface{}{
		"user_id":    savedUser.ID,
		"username":   savedUser.Username,
		"authorized": true,
		"exp":        time.Now().Add(time.Minute * 15).Unix(), // nolint
	})
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to generate token")
		h.RenderError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	h.RenderSuccess(w, http.StatusOK, generatedToken)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to extract id from request")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.DeleteUserByID(int64(id)); err != nil {
		h.LogEntry.WithError(err).Error("Unable to delete user")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.RenderSuccess(w, http.StatusNoContent, nil)
}
