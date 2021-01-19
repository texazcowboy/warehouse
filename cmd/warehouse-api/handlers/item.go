package handlers

import (
	"database/sql"
	"net/http"

	"github.com/texazcowboy/warehouse/internal/item"

	"github.com/go-playground/validator/v10"
	"github.com/texazcowboy/warehouse/cmd/warehouse-api/common"

	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type ItemHandler struct {
	*common.BaseHandler
	service item.ServiceInterface
}

func NewItemHandler(base *common.BaseHandler, service item.ServiceInterface) *ItemHandler {
	return &ItemHandler{base, service}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var i item.Item
	err := web.Decode(r, &i)
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

	if err = h.Validate.Struct(&i); err != nil {
		h.LogEntry.WithError(err).Error("Invalid request body")
		if _, ok := err.(*validator.InvalidValidationError); ok {
			h.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			h.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	err = h.service.CreateItem(&i)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to create item")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RenderSuccess(w, http.StatusCreated, &i)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to extract id from request")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetItem(int64(id))
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to get item by id")
		switch err {
		case sql.ErrNoRows:
			h.RenderError(w, http.StatusNotFound, err.Error())
		default:
			h.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	h.RenderSuccess(w, http.StatusOK, result)
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, _ *http.Request) {
	result, err := h.service.GetAllItems()
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to get items")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RenderSuccess(w, http.StatusOK, result)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to extract id from request")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	var i item.Item
	err = web.Decode(r, &i)
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

	if err = h.Validate.Struct(&i); err != nil {
		h.LogEntry.WithError(err).Error("Invalid request body ")
		if _, ok := err.(*validator.InvalidValidationError); ok {
			h.RenderError(w, http.StatusInternalServerError, err.Error())
		}
		if err, ok := err.(validator.ValidationErrors); ok {
			h.RenderError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	i.ID = int64(id)

	res, err := h.service.UpdateItem(&i)
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to update item")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if res == 0 {
		h.RenderError(w, http.StatusNotFound, "item not found")
		return
	}
	h.RenderSuccess(w, http.StatusOK, &i)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := web.ExtractIntFromRequest(r, "id")
	if err != nil {
		h.LogEntry.WithError(err).Error("Unable to extract id from request")
		h.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.DeleteItem(int64(id)); err != nil {
		h.LogEntry.WithError(err).Error("Unable to delete item by id")
		h.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.RenderSuccess(w, http.StatusNoContent, nil)
}
