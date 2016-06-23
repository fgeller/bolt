package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var dataDir string
var dataFilePath string

const (
	dataFileName = "data.json"
)

func mustReadPersisted() []string {
	_, err := os.Open(dataFilePath)
	switch {
	case os.IsNotExist(err):
		err = os.MkdirAll(dataDir, 0744)
		if err != nil {
			log.Fatalf("Error creating data dir err=%v", err)
		}

		_, err := os.Create(dataFilePath)
		if err != nil {
			log.Fatalf("Error creating data file err=%v", err)
		}

	case err != nil:
		log.Fatalf("Error checking data dir %#v err=%v", dataDir, err)
	}

	byts, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		log.Fatalf("Cannot read data err=%v", err)
	}

	if len(byts) == 0 {
		return []string{}
	}

	var dc []string
	err = json.Unmarshal(byts, &dc)
	if err != nil {
		log.Fatalf("Cannot unmarshal data err=%v", err)
	}

	return dc
}

func main() {
	flag.StringVar(&dataDir, "data", ".bolt", "directory to persist data to")
	flag.Parse()

	dataFilePath = filepath.Join(dataDir, dataFileName)

	serve(&bolt{data: mustReadPersisted()}, ":8077")
}
