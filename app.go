package main

import (
	"context"
)

// App is the Wails application struct. UI is fully client-side (IndexedDB via Dexie).
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet is a minimal binding kept for Wails health checks.
func (a *App) Greet(name string) string {
	return "Hello " + name + " from mdmark"
}
