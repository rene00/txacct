package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"transactionsearch/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Postcode struct {
	Number   string
	Locality string
	State    string
}

func createPostcodes(data [][]string) []Postcode {
	var postcodes []Postcode
	for _, line := range data {
		var postcode Postcode
		for j, field := range line {
			if j == 1 {
				postcode.Number = field
			} else if j == 2 {
				postcode.Locality = field
			} else if j == 4 {
				postcode.State = field
			}
		}
		postcodes = append(postcodes, postcode)
	}
	return postcodes
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	postcodeFile, ok := os.LookupEnv("TS_POSTCODE_FILE")
	if !ok {
		log.Fatal("TS_POSTCODE_FILE not set")
	}

	postgresURI, ok := os.LookupEnv("TS_POSTGRES_URI")
	if !ok {
		log.Fatal(fmt.Errorf("TS_POSTGRES_URI not set"))
	}

	db, err := sql.Open("pgx", postgresURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	boil.SetDB(db)

	if _, err := os.Stat(postcodeFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	f, err := os.Open(postcodeFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range createPostcodes(data) {
		var q []qm.QueryMod

		q = []qm.QueryMod{
			qm.Where("locality=? AND postcode=? AND s.name=?", strings.ToUpper(p.Locality), p.Number, strings.ToUpper(p.State)),
			qm.InnerJoin("state s on postcode.state_id = s.id"),
		}
		exists, err := models.Postcodes(q...).Exists(ctx, db)
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			q = []qm.QueryMod{
				qm.Where("name=?", p.State),
			}
			state, err := models.States(q...).One(ctx, db)
			if err != nil {
				log.Fatal(err)
			}

			postcode := models.Postcode{
				StateID:  state.ID,
				Postcode: p.Number,
				Locality: strings.ToUpper(p.Locality),
			}

			if err := postcode.Insert(ctx, db, boil.Infer()); err != nil {
				switch postcode.Locality {
				// These localities states conflict across data sources. Skip for now.
				case "URIARRA", "ALPURRURULAM", "COTTONVALE", "MINGOOLA":
					continue
				default:
				}
				log.Fatal(fmt.Errorf("failed inserting (state: %#+v) (postcode: %#+v): %w", state, postcode, err))
			}
			fmt.Printf("created %#v\n", postcode)
		}
	}
}
