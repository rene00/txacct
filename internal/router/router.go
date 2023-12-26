package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/about", aboutHandler())
	r.GET("/", homeHandler())
	r.POST("/transactions", transactionHandler(db))
	return r
}

// ViewData is passed into template.ExecuteTemplate.
type ViewData struct {
	CurrentPage string
}

func NewViewData(currentPage string) ViewData {
	return ViewData{CurrentPage: currentPage}
}
