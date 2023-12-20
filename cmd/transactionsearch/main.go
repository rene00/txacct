package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"transactionsearch/db/migrations"
	"transactionsearch/internal/router"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	rtr := router.NewRouter(db)
	listenAddrPort := fmt.Sprintf("%s:%s", httpHost, httpPort)
	log.Printf("Server listening on http://%s", listenAddrPort)
	if err := http.ListenAndServe(listenAddrPort, rtr); err != nil {
		log.Fatalf("%v", err)
	}
}
