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

	"golang.org/x/vuln/client"
	"gopkg.in/yaml.v3"
)

type flags struct {
	configFile *string
	url        *string
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
}

func main() {
	cmd := flag.NewFlagSet("transactionsearchtests", flag.ExitOnError)
	flags := flags{}
	flags.url = cmd.String("url", "", "url")
	flags.configFile = cmd.String("config-file", "", "config file path")
	cmd.Parse(os.Args[1:])

	if _, err := os.Stat(*flags.configFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	f, err := os.ReadFile(*flags.configFile)
	if err != nil {
		log.Fatal(err)
	}

	var c config.Config
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(*flags.url)
	if err != nil {
		log.Fatal(err)
	}

	tsClient := client.NewClient(u)
	for _, v := range c.Tests {

		requestBody := RequestBody{
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

		var responseBody ResponseBody

		if err := json.Unmarshal(body, &responseBody); err != nil {
			log.Fatal(err)
		}

		result := Compare(v, responseBody)

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
	test         config.Test
	responseBody ResponseBody
	mismatches   []string
}

func Compare(t config.Test, r ResponseBody) Result {
	var result Result
	result.test = t
	result.responseBody = r

	if t.Address != r.Address {
		result.mismatches = append(result.mismatches, fmt.Sprintf("address: want=%s got=%s", t.Address, r.Address))
	}

	if t.Description != r.Description {
		result.mismatches = append(result.mismatches, fmt.Sprintf("description: want=%s got=%s", t.Description, r.Description))
	}

	if t.Organisation != r.Organisation {
		result.mismatches = append(result.mismatches, fmt.Sprintf("organisation: want=%s got=%s", t.Organisation, r.Organisation))
	}

	return result
}
