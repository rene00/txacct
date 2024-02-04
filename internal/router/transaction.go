package router

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"transactionsearch/internal/transaction"
	"transactionsearch/internal/validation"

	"github.com/gin-gonic/gin"
)

type TransactionCache struct {
	Memo string
}

func transactionHandler(db *sql.DB, ch chan WorkerRequest) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var transactionRequest transaction.TransactionJSONRequest
		if err := ctx.ShouldBindJSON(&transactionRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("")})
			return
		}

		validator := validation.NewValidator()
		for _, rule := range []validation.Rule{
			validation.ValidateMaxLength(64),
			validation.ValidateMinLength(12),
			validation.ValidateRegexp(regexp.MustCompile("^[a-zA-Z0-9*\\-\\s]+$")),
		} {
			validator.Add(rule)
		}
		errors := validator.Validate(transactionRequest.Memo)
		if len(errors) > 0 {
			// pop first error and return
			ctx.JSON(http.StatusBadRequest, gin.H{"error": errors[0].Error()})
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
