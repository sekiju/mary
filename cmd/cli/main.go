package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	semver "github.com/Masterminds/semver/v3"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/constant"
	"github.com/sekiju/mdl/downloader"
	"github.com/sekiju/mdl/extractor"
	"github.com/sekiju/mdl/internal/tui"
	"github.com/sekiju/mdl/internal/util"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"os/signal"
	"resty.dev/v3"
	"sort"
	"strings"
	"time"
)

var version = "dev"

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
	if len(config.Params.Runtime.DownloadChapters) > 0 {
		return config.Params.Runtime.DownloadChapters
	}

	fmt.Print("Please enter the chapter URL: ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Split(strings.TrimSpace(input), " ")
}

func newFlagSet(name string) (fs *flag.FlagSet, primaryCookie, configPath *string) {
	fs = flag.NewFlagSet(name, flag.ExitOnError)
	primaryCookie = fs.String("cookie", "", "Cookie string for the current session")
	configPath = fs.String("config", "config.hcl", "Path to the config file")
	return fs, primaryCookie, configPath
}

func parse() string {
	var fs *flag.FlagSet
	var primaryCookie, configPath *string
	var args []string

	switch {
	case len(os.Args) > 1 && os.Args[1] == "chapters":
		config.Params.Runtime.ListChaptersMode = true
		fs, primaryCookie, configPath = newFlagSet("chapters")
		args = os.Args[2:]
	default:
		fs, primaryCookie, configPath = newFlagSet(constant.MDL)
		args = os.Args[1:]
	}

	if err := fs.Parse(args); err != nil {
		log.Fatal().Err(err).Send()
	}

	config.Params.Runtime.DownloadChapters = fs.Args()
	if primaryCookie != nil && *primaryCookie != "" {
		config.Params.Runtime.PrimaryCookie = primaryCookie
	}

	return *configPath
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	configPath := parse()

	isTUI := len(os.Args) == 1 && len(config.Params.Runtime.DownloadChapters) == 0

	config.Load(configPath)

	var statusMessages []string

	if config.Params.File.Application.CheckUpdates {
		if msg, err := checkForUpdates(); err != nil {
			log.Warn().Err(err).Msg("Failed to check for updates, continuing")
		} else if msg != "" {
			statusMessages = append(statusMessages, msg)
		}
	}

	if isTUI {
		return tui.Run(ctx, stop, nil, version, statusMessages)
	}

	chapterURLs := getChapterURLs()

	if config.Params.Runtime.ListChaptersMode {
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

			chapters, err := ext.FindChapters(ctx, chapterURL)
			if err != nil {
				return err
			}

			for _, chapter := range chapters {
				fmt.Println(chapter.ID, chapter.Title, chapter.URL)
			}
		}
	} else {
		loader := downloader.NewDownloader(ctx)

		for _, chapterURL := range chapterURLs {
			loader.Queue(chapterURL)
		}

		loader.Stop()
	}

	return nil
}

func checkForUpdates() (string, error) {
	currentVersion, err := semver.NewVersion(version)
	if err != nil {
		return "", nil
	}

	log.Trace().Msgf("Current version: %s | Checking for updates...", version)

	res, err := resty.New().R().Get("https://api.github.com/repos/sekiju/mdl/tags")
	if err != nil {
		return "", err
	}

	var tags []map[string]interface{}
	if err = json.Unmarshal(res.Bytes(), &tags); err != nil {
		return "", err
	}

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
		return "", nil
	}

	sort.Sort(semver.Collection(versions))

	msg := fmt.Sprintf("New version available: %s — https://github.com/sekiju/mdl/releases", versions[len(versions)-1].String())
	return msg, nil
}
