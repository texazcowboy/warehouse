package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/texazcowboy/warehouse/internal/foundation/logger"

	"github.com/gorilla/mux"
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
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid payload"); err != nil {
			e.LogEntry.Error(err)
			return
		}
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
			if err := web.RespondError(w, http.StatusInternalServerError, "Invalid validation input"); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		}
		var message string
		for _, err = range err.(validator.ValidationErrors) {
			message += err.Error() + ". "
		}
		if err := web.RespondError(w, http.StatusBadRequest, message); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	err = item.Create(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	if err = web.Respond(w, http.StatusCreated, &i); err != nil {
		e.LogEntry.Error(err)
		return
	}
}

func (e *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	result, err := item.Get(int64(id), e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			if err := web.RespondError(w, http.StatusNotFound, "Item not found"); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		default:
			if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		}
	}

	if err = web.Respond(w, http.StatusOK, result); err != nil {
		e.LogEntry.Error(err)
		return
	}
}

func (e *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	result, err := item.GetAll(e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		if err = web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}
	if err = web.Respond(w, http.StatusOK, result); err != nil {
		e.LogEntry.Error(err)
		return
	}
}

func (e *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	var i item.Item
	err = web.Decode(r, &i)
	if err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid payload"); err != nil {
			e.LogEntry.Error(err)
			return
		}
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
			if err := web.RespondError(w, http.StatusInternalServerError, "Invalid validation input"); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		}
		var message string
		for _, err = range err.(validator.ValidationErrors) {
			message += err.Error() + ". "
		}
		if err := web.RespondError(w, http.StatusBadRequest, message); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	i.ID = int64(id)

	err = item.Update(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			if err := web.RespondError(w, http.StatusNotFound, "Item not found"); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		default:
			if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
				e.LogEntry.Error(err)
				return
			}
			return
		}
	}

	if err = web.Respond(w, http.StatusOK, &i); err != nil {
		e.LogEntry.Error(err)
		return
	}
}

func (e *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	if err = item.Delete(int64(id), e.DB); err != nil {
		e.LogEntry.Error(err)
		if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			e.LogEntry.Error(err)
			return
		}
		return
	}

	if err = web.Respond(w, http.StatusNoContent, nil); err != nil {
		e.LogEntry.Error(err)
		return
	}
}
