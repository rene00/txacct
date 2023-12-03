package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"transactionsearch/internal/transaction"

	"github.com/gin-gonic/gin"
)

func transactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transactionRequest transaction.TransactionJSONRequest
		if err := ctx.ShouldBindJSON(&transactionRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("")})
			return
		}

		t := transaction.NewTransaction(transactionRequest.Memo, db)

		transactionHandlers := []transaction.TransactionHandler{
			transaction.NewTransactionState(),
			transaction.NewTransactionPostcode(),
			transaction.NewTransactionOrganisation(),
		}

		for _, handler := range transactionHandlers {
			if err := handler.Handle(ctx, db, t); err != nil {
				ctx.String(http.StatusInternalServerError, "failed handler")
				log.Printf("failed handler: %v", err)
				return
			}
		}

		transactionJSONResponse, err := transaction.NewTransactionJSONResponse(*t)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed transaction response")
			log.Printf("failed transaction response: %v", err)
			return
		}

		ctx.JSON(http.StatusOK, transactionJSONResponse)
	}
}
