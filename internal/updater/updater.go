package updater

import (
	"strings"

	"github.com/rs/zerolog/log"

	"559/pkg/request"
)

const Version = "development"

func Check() error {
	log.Trace().Msgf("current version: %s", Version)

	if Version == "development" {
		return nil
	}

	res, err := request.Get[[]GithubRelease]("https://api.github.com/repos/sekiju/mary/releases")
	if err != nil {
		return err
	}

	if isOutdated(Version, res.Body[0].TagName) {
		log.Warn().Msg("current version outdated, update to new version: https://github.com/sekiju/mary/releases")
	}

	return nil
}

func isOutdated(current, external string) bool {
	currentSlice := strings.Split(current, ".")
	externalSlice := strings.Split(external, ".")

	if len(currentSlice) != len(externalSlice) {
		return false
	}

	for i := range len(currentSlice) {
		if currentSlice[i] < externalSlice[i] {
			return true
		}
	}

	return false
}
