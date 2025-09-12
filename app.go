package main

import (
	"context"
	"log"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("App is starting up...")
}

func (a *App) shutdown(ctx context.Context) {
	log.Println("App is shutting down...")
}
