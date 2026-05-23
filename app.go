package main

import (
	"context"
	"fmt"

	appsvc "github.com/koshi/apex-custom-match-director/internal/app"
	"github.com/koshi/apex-custom-match-director/internal/domain"
)

type App struct {
	ctx context.Context
	svc *appsvc.Service
}

func NewApp() (*App, error) {
	svc, err := appsvc.NewService("")
	if err != nil {
		return nil, err
	}
	return &App{svc: svc}, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("[apex-custom-match-director] startup")
}

func (a *App) domReady(ctx context.Context) {
	fmt.Println("[apex-custom-match-director] dom ready")
}

func (a *App) shutdown(ctx context.Context) {
	fmt.Println("[apex-custom-match-director] shutdown")
	a.svc.Shutdown(ctx)
}

func (a *App) GetSettings() domain.PublicSettings {
	return a.svc.GetPublicSettings()
}

func (a *App) SaveSettings(settings domain.PublicSettings) error {
	return a.svc.SavePublicSettings(settings)
}

func (a *App) GetMatchState() domain.MatchState {
	return a.svc.GetMatchState()
}

func (a *App) GenerateTeams() ([]domain.Team, error) {
	return a.svc.GenerateTeams()
}

func (a *App) StartReceiver() error {
	return a.svc.StartReceiver()
}

func (a *App) StopReceiver() error {
	return a.svc.StopReceiver()
}

func (a *App) ReceiverStatus() domain.ReceiverStatus {
	return a.svc.ReceiverStatus()
}
