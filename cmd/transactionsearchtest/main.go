package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"transactionsearch/internal/testclient"
	"transactionsearch/internal/testconfig"
	"transactionsearch/internal/transaction"

	"gopkg.in/yaml.v3"
)

type flags struct {
	configFile         *string
	url                *string
	insecureSkipVerify *bool
}

type RequestBody struct {
	Memo string `json:"memo"`
}

type ResponseBody struct {
	Address      string `json:"address"`
	Description  string `json:"description"`
	Memo         string `json:"memo"`
	Organisation string `json:"organisation"`
	Postcode     string `json:"postcode"`
	State        string `json:"state"`
}

func main() {
	cmd := flag.NewFlagSet("transactionsearchtests", flag.ExitOnError)
	flags := flags{}
	flags.url = cmd.String("url", "", "url")
	flags.configFile = cmd.String("config-file", "", "config file path")
	flags.insecureSkipVerify = cmd.Bool("insecure-skip-verify", false, "skip verifying tls")
	cmd.Parse(os.Args[1:])

	if _, err := os.Stat(*flags.configFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	f, err := os.ReadFile(*flags.configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c testconfig.Config
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(*flags.url)
	if err != nil {
		log.Fatal(err)
	}

	clientOpts := []testclient.OptFunc{
		testclient.WithInsecureSkipVerify(*flags.insecureSkipVerify),
	}

	tsClient := testclient.NewClient(u, clientOpts...)
	for _, v := range c.Tests {

		requestBody := transaction.TransactionJSONRequest{
			Memo: v.Memo,
		}

		marshalled, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatal(err)
		}

		req, err := tsClient.NewRequest("POST", *u, bytes.NewReader(marshalled))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := tsClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var transactionJSONResponse transaction.TransactionJSONResponse

		if err := json.Unmarshal(body, &transactionJSONResponse); err != nil {
			log.Fatal(err)
		}

		result := Compare(v, transactionJSONResponse)

		if len(result.mismatches) == 0 {
			log.Printf("%s:\t[PASS]\n", v.Memo)
			continue
		}

		log.Printf("%s:\t[FAIL]\n", v.Memo)
		for _, mismatch := range result.mismatches {
			log.Printf("\t%s\n", mismatch)
		}
	}
}

type Result struct {
	test         testconfig.Test
	responseBody transaction.TransactionJSONResponse
	mismatches   []string
}

func Compare(t testconfig.Test, r transaction.TransactionJSONResponse) Result {
	var result Result
	result.test = t
	result.responseBody = r

	if t.State != r.GetState() {
		result.mismatches = append(result.mismatches, fmt.Sprintf("state: want=[%s] got=[%s]", t.State, r.GetState()))
	}

	if t.Postcode != r.GetPostcode() {
		result.mismatches = append(result.mismatches, fmt.Sprintf("postcode: want=[%d] got=[%d]", t.Postcode, r.GetPostcode()))
	}

	if t.Address != r.GetAddress() {
		result.mismatches = append(result.mismatches, fmt.Sprintf("address: want=[%s] got=[%s]", t.Address, r.GetAddress()))
	}

	if t.Description != r.GetDescription() {
		result.mismatches = append(result.mismatches, fmt.Sprintf("description: want=%s got=%s", t.Description, r.GetDescription()))
	}

	if t.Organisation != r.GetOrganisation() {
		result.mismatches = append(result.mismatches, fmt.Sprintf("organisation: want=%s got=%s", t.Organisation, r.GetOrganisation()))
	}

	return result
}
