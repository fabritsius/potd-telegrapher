package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fabritsius/potd-telegrapher/src/wikipedia"
)

func main() {
	today := time.Now().Format("2006-01-02")

	potd, err := wikipedia.ParsePOTD(today)
	if err != nil {
		log.Fatal(err)
	}

	potdStr, err := json.MarshalIndent(potd, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(potdStr))
}
