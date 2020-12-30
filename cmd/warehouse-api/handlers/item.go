package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"

	"github.com/texazcowboy/warehouse/internal/foundation/web"
	"github.com/texazcowboy/warehouse/internal/item"
)

type ItemHandler struct {
	*sql.DB
	*logger.Logger
	*validator.Validate
}

func NewItemHandler(db *sql.DB, logger *logger.Logger, validator *validator.Validate) *ItemHandler {
	return &ItemHandler{DB: db, Logger: logger, Validate: validator}
}

func (e *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var i item.Item
	err := web.Decode(r, &i)
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

	if err = e.Validate.Struct(&i); err != nil {
		e.LogEntry.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			e.renderError(w, http.StatusInternalServerError, err)
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.renderError(w, http.StatusBadRequest, err)
		}
		return
	}

	err = item.Create(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusInternalServerError, err)
		return
	}
	e.renderSuccess(w, http.StatusCreated, &i)
}

func (e *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusBadRequest, err)
		return
	}

	result, err := item.Get(int64(id), e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			e.renderError(w, http.StatusNotFound, err)
		default:
			e.renderError(w, http.StatusInternalServerError, err)
		}
		return
	}
	e.renderSuccess(w, http.StatusOK, result)
}

func (e *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	result, err := item.GetAll(e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusInternalServerError, err)
		return
	}
	e.renderSuccess(w, http.StatusOK, result)
}

func (e *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusBadRequest, err)
		return
	}

	var i item.Item
	err = web.Decode(r, &i)
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

	if err = e.Validate.Struct(&i); err != nil {
		e.LogEntry.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			e.renderError(w, http.StatusInternalServerError, err)
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.renderError(w, http.StatusBadRequest, err)
		}
		return
	}

	i.ID = int64(id)

	err = item.Update(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			e.renderError(w, http.StatusNotFound, err)
		default:
			e.renderError(w, http.StatusInternalServerError, err)
		}
		return
	}
	e.renderSuccess(w, http.StatusOK, &i)
}

func (e *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusBadRequest, err)
		return
	}

	if err = item.Delete(int64(id), e.DB); err != nil {
		e.LogEntry.Error(err)
		e.renderError(w, http.StatusInternalServerError, err)
		return
	}
	e.renderSuccess(w, http.StatusNoContent, nil)
}

func (e *ItemHandler) renderError(w http.ResponseWriter, status int, err error) {
	if err := web.RespondError(w, status, err.Error()); err != nil {
		e.LogEntry.Error(err)
	}
}

func (e *ItemHandler) renderSuccess(w http.ResponseWriter, status int, payload interface{}) {
	if err := web.Respond(w, status, payload); err != nil {
		e.LogEntry.Error(err)
	}
}
