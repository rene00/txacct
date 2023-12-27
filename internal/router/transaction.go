package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"transactionsearch/internal/transaction"

	"github.com/gin-gonic/gin"
)

func transactionHandler(db *sql.DB, ch chan WorkerRequest) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transactionRequest transaction.TransactionJSONRequest
		if err := ctx.ShouldBindJSON(&transactionRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("")})
			return
		}

		maxCap := cap(ch)
		if len(ch) == maxCap {
			ctx.String(429, "retry later")
			return
		}

		responseCh := make(chan WorkerResponse, 1)
		defer close(responseCh)

		t := transaction.NewTransaction(transactionRequest.Memo, db)
		req := WorkerRequest{Chan: responseCh, Transaction: t}
		ch <- req

		resp := <-responseCh
		if resp.Error != nil {
			ctx.String(http.StatusInternalServerError, "failed transaction response")
			log.Printf("failed transaction response: %v", resp.Error)
			return
		}

		transactionJSONResponse, err := transaction.NewTransactionJSONResponse(*resp.Transaction)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "failed transaction response")
			log.Printf("failed transaction response: %v", err)
			return
		}

		ctx.JSON(http.StatusOK, transactionJSONResponse)
	}
}
