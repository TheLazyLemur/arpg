# ARPG - 3D Action RPG Game

A 3D action RPG game built with Go and Raylib, featuring proper collision detection, entity management, and a maintainable architecture.

## Features

- **3D Graphics**: Top-down 3D perspective with player, enemies, and obstacles
- **Collision System**: Proper 3D collision detection using Raylib's built-in functions
- **Entity System**: Modular entity management for players, enemies, bullets, and obstacles
- **Configuration System**: JSON-based configuration with sensible defaults
- **Input System**: Keyboard and mouse input handling
- **Camera System**: Follow camera with smooth tracking
- **Rendering System**: Organized rendering pipeline with debug options

## Game Mechanics

- **Movement**: WASD keys to move the player
- **Aiming**: Mouse to aim
- **Shooting**: Left mouse button to shoot bullets
- **Collision**: Bullets can pass over short obstacles and enemies when fired at the right height
- **Enemy AI**: Enemies chase the player automatically
- **Health System**: Enemies have health bars and take damage from bullets

## Project Structure

```
arpg/
├── cmd/arpg/              # Main application entry point
├── pkg/                   # Public packages
│   ├── entities/          # Game entities (Player, Enemy, Bullet, Obstacle)
│   ├── collision/         # Collision detection system
│   ├── rendering/         # Rendering system
│   ├── input/            # Input handling
│   ├── camera/           # Camera management
│   └── config/           # Configuration system
├── internal/             # Private packages
│   └── game/             # Main game logic
├── assets/               # Game assets
│   ├── models/           # 3D models
│   ├── textures/         # Textures
│   └── sounds/           # Audio files
├── docs/                 # Documentation
└── scripts/              # Build and utility scripts
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Raylib dependencies (handled automatically)

### Building and Running

```bash
# Clone and navigate to the project
cd arpg

# Build the game
make build

# Run the game
make run

# Or build and run directly
go run ./cmd/arpg
```

### Development

```bash
# Set up development environment
make setup

# Run with hot reload (requires air)
make dev

# Format code
make fmt

# Run tests
make test
```

## Configuration

The game creates a `config.json` file on first run with default settings. You can modify this file to customize:

- **Window settings**: Size, fullscreen, FPS
- **Graphics settings**: FOV, wireframes, grid display
- **Gameplay settings**: Movement speed, bullet speed, enemy health
- **Debug settings**: FPS display, collision visualization

Example `config.json`:

```json
{
  "window": {
    "width": 1024,
    "height": 768,
    "title": "ARPG - 3D Scene",
    "fullscreen": false,
    "vsync": true,
    "target_fps": 60
  },
  "graphics": {
    "fov": 45.0,
    "draw_wires": false,
    "draw_grid": true,
    "anti_alias": true
  },
  "gameplay": {
    "player_speed": 5.0,
    "bullet_speed": 15.0,
    "bullet_lifetime": 10.0,
    "enemy_speed": 2.0,
    "enemy_health": 100.0,
    "bullet_damage": 25.0
  },
  "debug": {
    "show_fps": true,
    "show_collision": false,
    "show_health_bars": true
  }
}
```

## Controls

- **WASD** or **Arrow Keys**: Move player
- **Mouse**: Aim direction
- **Left Mouse Button**: Shoot bullets
- **F11**: Toggle fullscreen
- **F3**: Toggle debug info
- **ESC**: Exit game

## Architecture

### Entity System
- **Player**: Main character with position, rotation, health, and collision
- **Enemy**: AI-controlled entities that chase the player
- **Bullet**: Projectiles with physics and collision detection
- **Obstacle**: Static collidable objects (boxes and cylinders)

### Collision System
- Uses Raylib's built-in 3D collision detection
- Proper height-based collision for bullets passing over obstacles
- Efficient bounding box calculations

### Rendering Pipeline
1. Begin frame
2. 3D rendering mode
3. Draw ground, obstacles, entities
4. End 3D mode
5. Draw UI elements
6. End frame

## Build System

The project includes a comprehensive Makefile with the following targets:

```bash
make build       # Build the project
make run         # Build and run the game
make clean       # Clean build artifacts
make fmt         # Format code
make test        # Run tests
make build-all   # Build for all platforms
make release     # Create release build
make dev         # Run with hot reload
```

## Development Guidelines

### Adding New Entities
1. Create entity struct in `pkg/entities/`
2. Implement required methods (Update, GetBoundingBox, etc.)
3. Add rendering logic in `pkg/rendering/`
4. Update collision detection in `pkg/collision/`

### Adding New Systems
1. Create system package in `pkg/`
2. Implement initialization and update methods
3. Integrate with main game loop in `internal/game/`

### Code Style
- Use `gofmt` for formatting
- Follow Go naming conventions
- Add comments for public APIs
- Use meaningful variable names

## Troubleshooting

### Build Issues
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to resolve dependencies
- Check that Raylib builds correctly on your system

### Performance Issues
- Reduce enemy count in `initializeEnemies()`
- Lower FPS target in config
- Disable wireframe rendering

### Graphics Issues
- Update graphics drivers
- Try different window sizes
- Disable fullscreen mode

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make fmt` and `make test`
6. Submit a pull request

## License

This project is available under the MIT License.

## Future Enhancements

- [ ] Audio system integration
- [ ] Particle effects
- [ ] Level system
- [ ] Save/load functionality
- [ ] Multiplayer support
- [ ] Advanced AI behaviors
- [ ] Inventory system
- [ ] Skill trees
- [ ] Better graphics and models
- [ ] Performance optimizations