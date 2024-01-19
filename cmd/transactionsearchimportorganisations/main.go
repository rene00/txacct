package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"transactionsearch/internal/dataimporter"
	"transactionsearch/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/sync/errgroup"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	postgresURI, ok := os.LookupEnv("TS_POSTGRES_URI")
	if !ok {
		log.Fatal(fmt.Errorf("TS_POSTGRES_URI not set"))
	}

	organisationsFile, ok := os.LookupEnv("TS_ORGANISATIONS_FILE")
	if !ok {
		log.Fatal("TS_ORGANISATIONS_FILE not set")
	}

	worksheet, ok := os.LookupEnv("TS_ORGANISATIONS_WORKSHEET")
	if !ok {
		log.Fatal("TS_ORGANISATIONS_WORKSHEET not set")
	}

	state, ok := os.LookupEnv("TS_ORGANISATIONS_STATE")
	if !ok {
		log.Fatal("TS_ORGANISATIONS_STATE not set")
	}

	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	boil.SetDB(db)

	tsDebug := false
	if d, ok := os.LookupEnv("TS_DEBUG"); ok && d == "1" {
		tsDebug = true
	}
	boil.DebugMode = tsDebug

	const num = 8
	c := make(chan dataimporter.Row, num)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(c)
		return dataimporter.ProcessExcelData(ctx, c, organisationsFile, worksheet)
	})

	for i := 0; i < num; i++ {
		g.Go(func() error {
			for r := range c {
				switch {
				case strings.ToLower(state) == "nsw":
					d := dataimporter.NewDataImporter(models.OrganisationStateNSW{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "vic":
					d := dataimporter.NewDataImporter(models.OrganisationStateVic{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "qld":
					d := dataimporter.NewDataImporter(models.OrganisationStateQLD{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "tasmania":
					d := dataimporter.NewDataImporter(models.OrganisationStateTasmanium{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "act":
					d := dataimporter.NewDataImporter(models.OrganisationStateAct{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "sa":
					d := dataimporter.NewDataImporter(models.OrganisationStateSa{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "nt":
					d := dataimporter.NewDataImporter(models.OrganisationStateNT{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				case strings.ToLower(state) == "wa":
					d := dataimporter.NewDataImporter(models.OrganisationStateWa{})
					if err := d.Do(ctx, db, r); err != nil {
						return err
					}
				default:
					return fmt.Errorf("state not supported")
				}

			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
