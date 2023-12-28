package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"transactionsearch/db/migrations"
	"transactionsearch/internal/router"
	"transactionsearch/internal/transaction"
	"transactionsearch/models"

	"github.com/datasapiens/cachier"
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

	_cacheLRUSize, ok := os.LookupEnv("TS_CACHE_LRU_SIZE")
	if !ok {
		_cacheLRUSize = "300"
	}

	cacheLRUSize, err := strconv.Atoi(_cacheLRUSize)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	boil.SetDB(db)

	if doMigrate {
		err = migrations.DoMigrateDb(postgresURI)
		if err != nil {
			log.Fatal(err)
		}
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
			log.Printf("Server listening on http://%s", listenAddrPort)
			if err := s.ListenAndServe(); err != http.ErrServerClosed {
				panic(err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				log.Print("router shutting down")
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

		lc, err := cachier.NewLRUCache(cacheLRUSize,
			func(value interface{}) ([]byte, error) {
				return json.Marshal(value)
			},
			func(b []byte, value *interface{}) error {
				var res []models.OrganisationSlice
				return json.Unmarshal(b, &res)
				*value = res
				return nil
			},
			nil)
		if err != nil {
			return err
		}

		cache := cachier.MakeCache[[]models.OrganisationSlice](lc)
		fmt.Println(cache)

		store := transaction.Store{db, cache}

		for req := range workerCh {
			resp := router.WorkerResponse{Error: nil}

			transactionHandlers := []transaction.TransactionHandler{
				transaction.NewTransactionState(),
				transaction.NewTransactionPostcode(),
				transaction.NewTransactionOrganisation(),
			}

			for _, handler := range transactionHandlers {
				if err := handler.Handle(ctx, store, req.Transaction); err != nil {
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
		log.Printf("errgroup received error: %v", err)
		cancel()
		os.Exit(1)
	}

}
