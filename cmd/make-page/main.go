package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fabritsius/potd-telegrapher/src/telegraph"
)

func main() {
	today := time.Now().Format("2006-01-02")

	body, err := telegraph.MakeArticle(today)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(body)
}
