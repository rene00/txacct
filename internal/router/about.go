package router

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func aboutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tmpl *template.Template
		var err error

		tmpl, err = template.ParseFiles("./templates/main.tmpl", "./templates/about.tmpl")
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse files"))
			return
		}

		err = tmpl.ExecuteTemplate(ctx.Writer, "main", NewViewData("about"))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to execute template"))
			return
		}
	}
}
