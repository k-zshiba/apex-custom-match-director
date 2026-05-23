package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/koshi/apex-custom-match-director/internal/discord"
	"github.com/koshi/apex-custom-match-director/internal/domain"
	"github.com/koshi/apex-custom-match-director/internal/liveapi"
	"github.com/koshi/apex-custom-match-director/internal/storage"
)

type Service struct {
	mu        sync.Mutex
	settings  domain.Settings
	store     *storage.SettingsStore
	receiver  *liveapi.Receiver
	aggregate *domain.Aggregator
	discord   *discord.Client
}

func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			base = "."
		}
		dataDir = filepath.Join(base, "apex-custom-match-director")
	}
	store, err := storage.NewSettingsStore(filepath.Join(dataDir, "settings.sqlite"))
	if err != nil {
		return nil, err
	}
	settings, err := store.Load(context.Background())
	if err != nil {
		return nil, err
	}
	return &Service{
		settings:  settings,
		store:     store,
		aggregate: domain.NewAggregator(),
		discord:   discord.NewClient(),
	}, nil
}

func (s *Service) Shutdown(ctx context.Context) {
	_ = s.StopReceiver()
	_ = s.store.Close()
}

func (s *Service) GetPublicSettings() domain.PublicSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.settings.Public()
}

func (s *Service) SavePublicSettings(public domain.PublicSettings) error {
	settings := domain.NormalizeSettings(domain.SettingsFromPublic(public))

	s.mu.Lock()
	wasRunning := s.receiver != nil
	oldReceiver := s.receiver
	s.receiver = nil
	s.settings = settings
	s.mu.Unlock()

	if oldReceiver != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_ = oldReceiver.Shutdown(ctx)
		cancel()
	}
	if err := s.store.Save(context.Background(), settings); err != nil {
		return err
	}
	if wasRunning {
		return s.StartReceiver()
	}
	return nil
}

func (s *Service) GetMatchState() domain.MatchState {
	return s.aggregate.State()
}

func (s *Service) GenerateTeams() ([]domain.Team, error) {
	settings := s.currentSettings()
	state := s.aggregate.State()
	teams, err := domain.BalanceTeams(state.Rankings, settings.TeamCount, settings.TeamSizeCap)
	if err != nil {
		return nil, err
	}
	s.aggregate.SetTeams(teams)
	return teams, nil
}

func (s *Service) StartReceiver() error {
	s.mu.Lock()
	if s.receiver != nil {
		s.mu.Unlock()
		return nil
	}
	settings := s.settings
	receiver := liveapi.NewReceiver(settings.LiveAPIPort, s.handleEvent)
	s.receiver = receiver
	s.mu.Unlock()

	if err := receiver.Start(); err != nil {
		s.mu.Lock()
		if s.receiver == receiver {
			s.receiver = nil
		}
		s.mu.Unlock()
		return err
	}
	return nil
}

func (s *Service) StopReceiver() error {
	s.mu.Lock()
	receiver := s.receiver
	s.receiver = nil
	s.mu.Unlock()
	if receiver == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return receiver.Shutdown(ctx)
}

func (s *Service) ReceiverStatus() domain.ReceiverStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.receiver == nil {
		return domain.ReceiverStatus{}
	}
	return domain.ReceiverStatus{Running: true, Address: s.receiver.Address()}
}

func (s *Service) handleEvent(event domain.LiveEvent) {
	state := s.aggregate.Apply(event)
	settings := s.currentSettings()
	if !settings.AutoPostDiscord || !state.Completed || state.Posted {
		return
	}
	teams, err := domain.BalanceTeams(state.Rankings, settings.TeamCount, settings.TeamSizeCap)
	if err != nil {
		s.aggregate.SetLastError(err.Error())
		return
	}
	state = s.aggregate.SetTeams(teams)
	if err := s.discord.PostResult(settings.DiscordWebhookURL, state, teams); err != nil {
		s.aggregate.SetLastError(fmt.Sprintf("%v", discord.RedactSecrets(err)))
		return
	}
	s.aggregate.MarkPosted()
}

func (s *Service) currentSettings() domain.Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.settings
}
