package user

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

type Handler struct {
	*common.BaseHandler
}

func NewUserHandler(base *common.BaseHandler) *Handler {
	return &Handler{base}
}

func (e *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := web.Decode(r, &u)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
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
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	err = user.Create(&u, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.Password = "***"
	e.RenderSuccess(w, http.StatusCreated, &u)
}

func (e *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var u user.User
	err := web.Decode(r, &u)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
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
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	savedUser, err := user.GetByUsername(u.Username, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if !crypto.Equals(u.Password, savedUser.Password) {
		e.RenderError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	generatedToken, err := security.GenerateToken(map[string]interface{}{
		"user_id":    savedUser.ID,
		"username":   savedUser.Username,
		"authorized": true,
		"exp":        time.Now().Add(time.Minute * 15).Unix(), // nolint
	})
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	e.RenderSuccess(w, http.StatusOK, generatedToken)
}

func (e *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = user.Delete(int64(id), e.DB); err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	e.RenderSuccess(w, http.StatusNoContent, nil)
}
