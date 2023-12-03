package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/transactions", transactionHandler(db))
	return r
}
