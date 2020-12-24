package handlers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
	"github.com/texazcowboy/warehouse/internal/item"
	"log"
	"net/http"
	"strconv"
)

func (e *Env) CreateItem(w http.ResponseWriter, r *http.Request) {
	var i item.Item
	err := web.Decode(r, &i)
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid payload"); err != nil {
			log.Println(err)
			return
		}
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	err = item.Create(&i, e.DB)
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if err = web.Respond(w, http.StatusCreated, &i); err != nil {
		log.Println(err)
		return
	}
}

func (e *Env) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			log.Println(err)
			return
		}
		return
	}

	result, err := item.Get(int64(id), e.DB)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			if err := web.RespondError(w, http.StatusNotFound, "Item not found"); err != nil {
				log.Println(err)
				return
			}
			return
		default:
			if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
				log.Println(err)
				return
			}
			return
		}
	}

	if err = web.Respond(w, http.StatusOK, result); err != nil {
		log.Println(err)
		return
	}
}

func (e *Env) GetItems(w http.ResponseWriter, r *http.Request) {
	result, err := item.GetAll(e.DB)
	if err != nil {
		log.Println(err)
		if err = web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}
	if err = web.Respond(w, http.StatusOK, result); err != nil {
		log.Println(err)
		return
	}
}

func (e *Env) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			log.Println(err)
			return
		}
		return
	}

	var i item.Item
	err = web.Decode(r, &i)
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid payload"); err != nil {
			log.Println(err)
			return
		}
		return
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	i.Id = int64(id)

	err = item.Update(&i, e.DB)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			if err := web.RespondError(w, http.StatusNotFound, "Item not found"); err != nil {
				log.Println(err)
				return
			}
			return
		default:
			if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
				log.Println(err)
				return
			}
			return
		}
	}

	if err = web.Respond(w, http.StatusOK, &i); err != nil {
		log.Println(err)
		return
	}
}

func (e *Env) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusBadRequest, "Invalid item id"); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if err = item.Delete(int64(id), e.DB); err != nil {
		log.Println(err)
		if err := web.RespondError(w, http.StatusInternalServerError, err.Error()); err != nil {
			log.Println(err)
			return
		}
		return
	}

	if err = web.Respond(w, http.StatusNoContent, nil); err != nil {
		log.Println(err)
		return
	}
}
