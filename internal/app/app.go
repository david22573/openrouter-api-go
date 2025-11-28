package app

import (
	"github.com/david22573/openrouter-api-go/internal/config"
	"github.com/david22573/openrouter-api-go/pkg/openrouter"
)

type App struct {
	Config *config.Config
	Client *openrouter.Client
}

var A App
