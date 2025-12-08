package service

import (
	"bufio"
	"djtracker/internal/config"
	"djtracker/internal/model"
	"djtracker/internal/repository"
	"djtracker/internal/utils"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Service struct {
	log           *slog.Logger
	config        *config.Config
	repo          *repository.Repository
	reader        *bufio.Reader
	liveTracklist chan string

	trackBroadcaster *Broadcaster[*model.TrackDTO]
}

func New(log *slog.Logger, config *config.Config, repo *repository.Repository) *Service {
	return &Service{
		log:           log,
		config:        config,
		repo:          repo,
		liveTracklist: make(chan string, 1),

		trackBroadcaster: NewBroadcaster[*model.TrackDTO](log),
	}
}

// SubscribeForTracks Créer un nouveau channel abonné à la réception des tracks
func (s *Service) SubscribeForTracks() (chan *model.TrackDTO, func()) {
	return s.trackBroadcaster.Subscribe(1)
}

func (s *Service) StartTracking() error {
	tracklistFile, err := os.Open(s.config.Tracker.History.Path)
	if err != nil {
		return err
	}

	stat, err := tracklistFile.Stat()
	if err != nil {
		return err
	}
	_, err = tracklistFile.Seek(stat.Size(), 0)
	if err != nil {
		return err
	}

	go s.readTracks(tracklistFile)
	go s.analyseTracks()
	return nil
}

// analyseTracks Lit les tracks brutes reçues de liveTracklist
// Transfer les informations de la TrackDTO vers le channel Tracks
func (s *Service) analyseTracks() {
	for trackText := range s.liveTracklist {
		parsedTrack, err := ParseLine(trackText)
		if err != nil {
			s.log.Error("Unable to parse track line", "track_line", trackText)
			continue
		}

		track := &model.Track{
			Artist: utils.SafeTrim(parsedTrack.Artist),
			Name:   strings.TrimSpace(parsedTrack.Name),
			PlayAt: time.Now(),
		}

		fileTrack, err := s.findTrackFile(track.Name)
		if err != nil {
			s.log.Error("Track file not found", "track", track.Name)
		} else {
			s.log.Debug("Track file found", "track", fileTrack)

			trackDuration, err := s.findTrackDuration(fileTrack)
			if err != nil {
				s.log.Error("Failed to retrieve time data", "path", fileTrack.Path)
			} else {
				track.Duration = &trackDuration
			}
		}

		s.repo.AddTrackToHistory(track)
		s.trackBroadcaster.Broadcast(track.ToDTO())
	}
}

func (s *Service) findTrackFile(track string) (*FileTrack, error) {
	for _, sourceFolder := range s.config.Tracker.Source.Paths {
		if file, err := FindFile(sourceFolder, track); err == nil {
			return file, nil
		}
	}
	return nil, fmt.Errorf("track file not found for %s", track)
}

// readTracks Lit le fichier tracklist de VirtualDJ
// Transfer les informations brutes vers le channel liveTracklist
func (s *Service) readTracks(file *os.File) {
	reader := bufio.NewReader(file)
	defer utils.SafeDeferClose(file, s.log)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			s.log.Error("Error while reading file", err)
		}
		data = strings.TrimRight(data, "\r\n")

		s.liveTracklist <- data
	}
}
