package item

import (
	"database/sql"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/common"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/texazcowboy/warehouse/internal/foundation/web"
	"github.com/texazcowboy/warehouse/internal/item"
)

type Handler struct {
	*common.BaseHandler
}

func NewItemHandler(base *common.BaseHandler) *Handler {
	return &Handler{base}
}

func (e *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var i item.Item
	err := web.Decode(r, &i)
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

	if err = e.Validate.Struct(&i); err != nil {
		e.LogEntry.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	err = item.Create(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.RenderSuccess(w, http.StatusCreated, &i)
}

func (e *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := item.Get(int64(id), e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			e.RenderError(w, http.StatusNotFound, err.Error())
		default:
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	e.RenderSuccess(w, http.StatusOK, result)
}

func (e *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	result, err := item.GetAll(e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.RenderSuccess(w, http.StatusOK, result)
}

func (e *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	var i item.Item
	err = web.Decode(r, &i)
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

	if err = e.Validate.Struct(&i); err != nil {
		e.LogEntry.Error(err)
		if _, ok := err.(*validator.InvalidValidationError); ok {
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			e.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	i.ID = int64(id)

	err = item.Update(&i, e.DB)
	if err != nil {
		e.LogEntry.Error(err)
		switch err {
		case sql.ErrNoRows:
			e.RenderError(w, http.StatusNotFound, err.Error())
		default:
			e.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	e.RenderSuccess(w, http.StatusOK, &i)
}

func (e *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = item.Delete(int64(id), e.DB); err != nil {
		e.LogEntry.Error(err)
		e.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	e.RenderSuccess(w, http.StatusNoContent, nil)
}
