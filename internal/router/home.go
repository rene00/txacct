package router

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func homeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tmpl *template.Template
		var err error

		tmpl, err = template.ParseFiles("./templates/main.tmpl", "./templates/home.tmpl")
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to parse files"))
			return
		}

		err = tmpl.ExecuteTemplate(ctx.Writer, "main", NewViewData("home"))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to execute template"))
			return
		}
	}
}
