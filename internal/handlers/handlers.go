package handlers

import (
	"database/sql"
	"transactionsearch/internal/logwrap"
)

type Handlers struct {
	DB     *sql.DB
	Logger *logwrap.LogWrap
}
