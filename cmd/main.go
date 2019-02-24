package main

import (
	"log"
	"os"

	crawler "github.com/timtosi/mcrawler/internal"
	"github.com/timtosi/mcrawler/internal/domain"
	"github.com/timtosi/mcrawler/internal/extractor"
)

func main() {
	if len(os.Args[1]) == 0 {
		log.Fatal(`usage: ./mcrawler <BASE_URL>`)
	}

	t := domain.NewTarget(os.Args[1])
	m := crawler.NewMapper()
	f, err := crawler.NewFollower(t.BaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := crawler.NewCrawler().Run(
		t,
		crawler.NewArchiver(),
		m,
		f,
		crawler.NewWorker(),
		extractor.NewExtractor(extractor.GetImg, extractor.GetLinkBasic),
	); err != nil {
		log.Fatal(err)
	}

	m.Render()
	log.Printf("shutdown")
}
