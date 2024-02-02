package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"transactionsearch/db/migrations"
	"transactionsearch/internal/handlers"
	"transactionsearch/internal/logwrap"
	"transactionsearch/internal/router"
	"transactionsearch/internal/transaction"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/sync/errgroup"
)

func main() {
	var err error

	postgresURI, ok := os.LookupEnv("TS_POSTGRES_URI")
	if !ok {
		log.Fatal(fmt.Errorf("TS_POSTGRES_URI not set"))
	}

	httpHost, ok := os.LookupEnv("TS_HTTP_HOST")
	if !ok {
		httpHost = "127.0.0.1"
	}

	httpPort, ok := os.LookupEnv("TS_HTTP_PORT")
	if !ok {
		httpPort = "3000"
	}

	doMigrate := false
	if s, ok := os.LookupEnv("TS_DO_MIGRATE"); ok && s == "1" {
		doMigrate = true
	}

	_workerQueueLength, ok := os.LookupEnv("TS_WORKER_QUEUE_LENGTH")
	if !ok {
		_workerQueueLength = "2"
	}

	workerQueueLength, err := strconv.Atoi(_workerQueueLength)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	boil.SetDB(db)

	logger := logwrap.New("logger", os.Stdout, true)
	logger.SetLevel(logwrap.INFO)

	tsDebug := false
	if d, ok := os.LookupEnv("TS_DEBUG"); ok && d == "1" {
		tsDebug = true
		logger.SetLevel(logwrap.DEBUG)
	}
	boil.DebugMode = tsDebug

	if doMigrate {
		err = migrations.DoMigrateDb(postgresURI)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("completed database migrations")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)

	workerCh := make(chan router.WorkerRequest, workerQueueLength)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		rtr := router.NewRouter(db, workerCh)
		listenAddrPort := fmt.Sprintf("%s:%s", httpHost, httpPort)
		s := &http.Server{
			Addr:    listenAddrPort,
			Handler: rtr,
		}

		go func() {
			logger.Info(fmt.Sprintf("Server listening on http://%s", listenAddrPort))
			if err := s.ListenAndServe(); err != http.ErrServerClosed {
				panic(err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				logger.Info("router shutting down")
				return s.Shutdown(ctx)
			}
		}
		return nil
	})

	g.Go(func() error {
		for {
			select {
			case <-sigC:
				close(workerCh)
				return fmt.Errorf("received SIGTERM/SIGINT, exiting")
			case <-ctx.Done():
				close(workerCh)
				return nil
			}
		}
	})

	g.Go(func() error {

		handlers := handlers.Handlers{db, logger}

		for req := range workerCh {
			resp := router.WorkerResponse{Error: nil}

			transactionHandlers := []transaction.TransactionHandler{
				transaction.NewTransactionState(),
				transaction.NewTransactionPostcode(),
				transaction.NewTransactionOrganisation(),
			}

			for _, handler := range transactionHandlers {
				if err := handler.Handle(ctx, handlers, req.Transaction); err != nil {
					resp.Error = err
					break
				}
			}

			resp.Transaction = req.Transaction
			req.Chan <- resp
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Info(fmt.Sprintf("errgroup received error: %v", err))
		cancel()
		os.Exit(1)
	}

}
