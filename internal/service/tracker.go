package service

import (
	"bufio"
	"djtracker/internal/config"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Service struct {
	log    *slog.Logger
	config *config.Config
	reader *bufio.Reader
	Tracks chan string
}

func New(log *slog.Logger, config *config.Config) *Service {
	return &Service{
		log:    log,
		config: config,
		Tracks: make(chan string, 1),
	}
}

func (s *Service) StartTracking() error {
	file, err := os.Open(s.config.Tracker.Path)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	_, err = file.Seek(stat.Size(), 0)
	if err != nil {
		return err
	}

	go s.readTracks(file, s.Tracks)
	return nil
}

func (s *Service) readTracks(file *os.File, channel chan string) {
	reader := bufio.NewReader(file)
	defer s.handleClose(file)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			s.log.Error("Error while reading file", err)
		}
		data = strings.TrimRight(data, "\r\n")

		channel <- data
	}
}

func (s *Service) handleClose(file *os.File) {
	err := file.Close()
	if err != nil {
		s.log.Error("Failed to close file", err)
	}
}
