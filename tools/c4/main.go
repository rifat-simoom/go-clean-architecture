package main

import (
	"context"
	"fmt"
	trainerService "github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/infrastructure/configs"
	trainingsService "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/infrastructure/configs"
	"os"

	"github.com/krzysztofreczek/go-structurizr/pkg/scraper"
	"github.com/krzysztofreczek/go-structurizr/pkg/view"
)

const (
	scraperConfig = "scraper.yml"
	viewConfig    = "view.yml"
	outputFile    = "out/view-%s.plantuml"
)

func main() {
	ctx := context.Background()

	trainerApp := trainerService.NewApplication(ctx)
	scrape(trainerApp, "trainer")

	trainingsApp, _ := trainingsService.NewApplication(ctx)
	scrape(trainingsApp, "trainings")
}

func scrape(app interface{}, name string) {
	s, err := scraper.NewScraperFromConfigFile(scraperConfig)
	if err != nil {
		panic(err)
	}

	structure := s.Scrape(app)

	outFileName := fmt.Sprintf(outputFile, name)
	outFile, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = outFile.Close()
	}()

	v, err := view.NewViewFromConfigFile(viewConfig)
	if err != nil {
		panic(err)
	}

	err = v.RenderStructureTo(structure, outFile)
	if err != nil {
		panic(err)
	}
}
