package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Masterminds/semver"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/constant"
	"github.com/sekiju/mdl/downloader"
	"github.com/sekiju/mdl/extractor"
	"github.com/sekiju/mdl/internal/util"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"resty.dev/v3"
	"sort"
	"strings"
	"time"
)

var version = "1.0.0"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	isGUI := !util.IsRunningFromCLI()

	if err := run(); err != nil {
		log.Error().Err(err).Send()
		if isGUI {
			waitForInput()
		}
		os.Exit(1)
	}

	if isGUI {
		waitForInput()
	}
}

func waitForInput() {
	fmt.Println("\nPress Enter to exit...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getChapterURLs() []string {
	if len(config.Params.DownloadChapters) > 0 {
		return config.Params.DownloadChapters
	}

	fmt.Print("Please enter the chapter URL: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Split(strings.TrimSpace(input), " ")
}

func parse() string {
	rootFlags := flag.NewFlagSet(constant.MDL, flag.ExitOnError)
	primaryCookie := rootFlags.String("cookie", "", "Cookie string for the current session")
	configPath := rootFlags.String("config", "config.hcl", "Path to the config file")

	if len(os.Args) > 1 && os.Args[1] == "chapters" {
		config.Params.ListChaptersMode = true
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	if err := rootFlags.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Send()
	}

	config.Params.DownloadChapters = rootFlags.Args()
	if primaryCookie != nil && *primaryCookie != "" {
		config.Params.PrimaryCookie = primaryCookie
	}

	return *configPath
}

func run() error {
	config.Load(parse())

	if config.Params.Application.CheckUpdates {
		if err := checkForUpdates(); err != nil {
			log.Warn().Err(err).Msg("Failed to check for updates, continuing")
		}
	}

	chapterURLs := getChapterURLs()

	if config.Params.ListChaptersMode {
		for _, chapterURL := range chapterURLs {
			parsedURL, err := url.Parse(chapterURL)
			if err != nil {
				return err
			}

			ext, err := extractor.NewExtractor(parsedURL.Hostname())
			if err != nil {
				return err
			}

			if listing, ok := ext.(manga.ChapterListingFeature); ok && !listing.SupportsChapterListing() {
				log.Warn().Str("url", chapterURL).Msg("This site does not support chapter listing, skipping")
				continue
			}

			chapters, err := ext.FindChapters(chapterURL)
			if err != nil {
				return err
			}

			for _, chapter := range chapters {
				fmt.Println(chapter.ID, chapter.Title, chapter.URL)
			}
		}
	} else {
		// Default download mode

		loader := downloader.NewDownloader()

		for _, chapterURL := range chapterURLs {
			loader.Queue(chapterURL)
		}

		loader.Stop()
	}

	return nil
}

func checkForUpdates() error {
	log.Trace().Msgf("Current version: %s | Checking for updates...", version)

	res, err := resty.New().R().Get("https://api.github.com/repos/sekiju/mdl/tags")
	if err != nil {
		return err
	}

	var tags []map[string]interface{}
	if err = json.Unmarshal(res.Bytes(), &tags); err != nil {
		return err
	}

	currentVersion := semver.MustParse(version)

	var versions []*semver.Version
	for _, tag := range tags {
		name, ok := tag["name"].(string)
		if !ok {
			continue
		}

		if v, err := semver.NewVersion(name); err == nil && v.GreaterThan(currentVersion) {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return nil
	}

	sort.Sort(semver.Collection(versions))

	log.Info().Msgf("New downloader version available: %s - download release from: https://github.com/sekiju/mdl/releases", versions[len(versions)-1].String())

	return nil
}
