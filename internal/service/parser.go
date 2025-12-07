package service

import (
	"fmt"
	"regexp"
)

var trackRegex = regexp.MustCompile(`^(?P<hour>\d{2}:\d{2}) ?: ?(?P<artist>.+?) - (?P<name>.+)$`)

type ParsedTrack struct {
	Hour   string
	Artist *string
	Name   string
}

// ParseLine traite une ligne de l'historique VirtualDL
// pour le transformer en une struct ParsedTrack.
func ParseLine(line string) (*ParsedTrack, error) {
	matches := trackRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("line does not match pattern: %s", line)
	}

	track := &ParsedTrack{}
	for i, name := range trackRegex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		switch name {
		case "hour":
			track.Hour = matches[i]
		case "artist":
			track.Artist = &matches[i]
		case "name":
			track.Name = matches[i]
		}
	}

	return track, nil
}
