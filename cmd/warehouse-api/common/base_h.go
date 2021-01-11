package common

import (
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/texazcowboy/warehouse/internal/foundation/logger"
	"github.com/texazcowboy/warehouse/internal/foundation/web"
)

type BaseHandler struct {
	*sql.DB
	*logger.Logger
	*validator.Validate
}

func (e *BaseHandler) RenderError(w http.ResponseWriter, status int, message string) {
	if err := web.RespondError(w, status, message); err != nil {
		e.LogEntry.WithError(err).Error("Unable to write response")
	}
}

func (e *BaseHandler) RenderSuccess(w http.ResponseWriter, status int, payload interface{}) {
	if err := web.Respond(w, status, payload); err != nil {
		e.LogEntry.WithError(err).Error("Unable to write response")
	}
}
