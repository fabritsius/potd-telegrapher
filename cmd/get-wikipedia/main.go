package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fabritsius/potd-telegrapher/src/wikipedia"
)

func main() {
	potd, err := wikipedia.ParsePOTD("2022-10-29")
	if err != nil {
		log.Fatal(err)
	}

	potdStr, err := json.MarshalIndent(potd, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(potdStr))
}
