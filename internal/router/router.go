package router

import (
	"database/sql"
	"transactionsearch/internal/transaction"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB, ch chan WorkerRequest) *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/about", aboutHandler())
	r.GET("/", homeHandler())
	r.POST("/transactions", transactionHandler(db, ch))
	return r
}

// ViewData is passed into template.ExecuteTemplate.
type ViewData struct {
	CurrentPage string
}

func NewViewData(currentPage string) ViewData {
	return ViewData{CurrentPage: currentPage}
}

type WorkerResponse struct {
	Transaction *transaction.Transaction
	Error       error
}

type WorkerRequest struct {
	Chan        chan WorkerResponse
	Transaction *transaction.Transaction
}
