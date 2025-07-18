package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all game configuration
type Config struct {
	Window   WindowConfig   `json:"window"`
	Graphics GraphicsConfig `json:"graphics"`
	Audio    AudioConfig    `json:"audio"`
	Gameplay GameplayConfig `json:"gameplay"`
	Debug    DebugConfig    `json:"debug"`
}

// WindowConfig contains window-related settings
type WindowConfig struct {
	Width      int32   `json:"width"`
	Height     int32   `json:"height"`
	Title      string  `json:"title"`
	Fullscreen bool    `json:"fullscreen"`
	VSync      bool    `json:"vsync"`
	TargetFPS  int32   `json:"target_fps"`
}

// GraphicsConfig contains graphics-related settings
type GraphicsConfig struct {
	FOV         float32 `json:"fov"`
	DrawWires   bool    `json:"draw_wires"`
	DrawGrid    bool    `json:"draw_grid"`
	AntiAlias   bool    `json:"anti_alias"`
}

// AudioConfig contains audio-related settings
type AudioConfig struct {
	MasterVolume float32 `json:"master_volume"`
	SFXVolume    float32 `json:"sfx_volume"`
	MusicVolume  float32 `json:"music_volume"`
}

// GameplayConfig contains gameplay-related settings
type GameplayConfig struct {
	PlayerSpeed   float32 `json:"player_speed"`
	BulletSpeed   float32 `json:"bullet_speed"`
	BulletLife    float32 `json:"bullet_lifetime"`
	EnemySpeed    float32 `json:"enemy_speed"`
	EnemyHealth   float32 `json:"enemy_health"`
	BulletDamage  float32 `json:"bullet_damage"`
}

// DebugConfig contains debug-related settings
type DebugConfig struct {
	ShowFPS        bool `json:"show_fps"`
	ShowCollision  bool `json:"show_collision"`
	ShowHealthBars bool `json:"show_health_bars"`
}

// Default returns a configuration with sensible defaults
func Default() *Config {
	return &Config{
		Window: WindowConfig{
			Width:      1024,
			Height:     768,
			Title:      "ARPG - 3D Scene",
			Fullscreen: false,
			VSync:      true,
			TargetFPS:  60,
		},
		Graphics: GraphicsConfig{
			FOV:         45.0,
			DrawWires:   false,
			DrawGrid:    true,
			AntiAlias:   true,
		},
		Audio: AudioConfig{
			MasterVolume: 1.0,
			SFXVolume:    1.0,
			MusicVolume:  0.7,
		},
		Gameplay: GameplayConfig{
			PlayerSpeed:  5.0,
			BulletSpeed:  15.0,
			BulletLife:   10.0,
			EnemySpeed:   2.0,
			EnemyHealth:  100.0,
			BulletDamage: 25.0,
		},
		Debug: DebugConfig{
			ShowFPS:        true,
			ShowCollision:  false,
			ShowHealthBars: true,
		},
	}
}

// Load loads configuration from file, or returns default if file doesn't exist
func Load() (*Config, error) {
	configPath := "config.json"
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		cfg := Default()
		if err := cfg.Save(configPath); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the configuration to a file
func (c *Config) Save(path string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}