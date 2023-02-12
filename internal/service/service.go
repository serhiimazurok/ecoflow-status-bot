package service

import "github.com/serhiimazurok/ecoflow-status-bot/internal/repository"

type Services struct {
}

type Deps struct {
	Repos *repository.Repositories
}

func NewServices(deps Deps) *Services {
	return &Services{}
}
