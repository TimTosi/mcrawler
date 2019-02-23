package main

import (
	"log"

	crawler "github.com/timtosi/mcrawler/internal"
	"github.com/timtosi/mcrawler/internal/domain"
	"github.com/timtosi/mcrawler/internal/extractor"
)

func main() {
	t := domain.NewTarget("localhost:8080") // TO UPDATE

	c := crawler.NewCrawler()
	w := crawler.NewWorker()
	e := extractor.NewExtractor(extractor.GetImg, extractor.GetLinkBasic)
	a := crawler.NewArchiver()
	f, err := crawler.NewFollower(t.BaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Run(t, w, e, f, a); err != nil {
		log.Fatal(err)
	}
	log.Printf("shutdown")
}
