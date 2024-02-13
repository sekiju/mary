package speed_binb

import (
	"559/internal/static"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
)

func handleV016201(uri url.URL, imageChan chan<- static.Image) error {
	log.Trace().Msg("bimb version 016201")
	return fmt.Errorf("v016201 unimplimented")
}
