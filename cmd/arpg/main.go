package main

import (
	"log"

	"arpg/internal/game"
	"arpg/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize and run the game with scene system
	g := game.New(cfg)
	if err := g.Run(); err != nil {
		log.Fatal("Game error:", err)
	}
}

