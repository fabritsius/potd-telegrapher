package main

import (
	"os"
	"time"

	"github.com/fabritsius/potd-telegrapher/src/telegram"
)

func main() {
	date, found := os.LookupEnv("ARTICLE_DATE")

	if !found || date == "" {
		date = time.Now().Format("2006-01-02")
	}

	telegram.PostArticle(date)
}
