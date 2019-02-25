package main

import (
	"log"
	"os"

	"github.com/timtosi/mcrawler/internal"
	"github.com/timtosi/mcrawler/internal/crawler"
	"github.com/timtosi/mcrawler/internal/domain"
	"github.com/timtosi/mcrawler/internal/extractor"
	"github.com/timtosi/mcrawler/internal/mapper"
)

func main() {
	if len(os.Args[1]) == 0 {
		log.Fatal(`usage: ./mcrawler <BASE_URL>`)
	}

	t := domain.NewTarget(os.Args[1])
	m := mapper.NewMapper()
	f, err := internal.NewFollower(t.BaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := crawler.NewCrawler().Run(
		t,
		internal.NewArchiver(),
		m,
		f,
		internal.NewWorker(),
		extractor.NewExtractor(extractor.GetImg, extractor.GetLinkNoFollow),
	); err != nil {
		log.Fatal(err)
	}

	m.Render()
	log.Printf("shutdown")
}
