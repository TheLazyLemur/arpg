package scenes

import (
	"encoding/json"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Vector3Data represents a 3D vector in JSON
type Vector3Data struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

// ToVector3 converts Vector3Data to rl.Vector3
func (v Vector3Data) ToVector3() rl.Vector3 {
	return rl.Vector3{X: v.X, Y: v.Y, Z: v.Z}
}

// FromVector3 creates Vector3Data from rl.Vector3
func FromVector3(v rl.Vector3) Vector3Data {
	return Vector3Data{X: v.X, Y: v.Y, Z: v.Z}
}

// PlayerData represents player configuration in JSON
type PlayerData struct {
	SpawnPoint Vector3Data `json:"spawn_point"`
	Speed      float32     `json:"speed"`
	Health     float32     `json:"health"`
	MaxHealth  float32     `json:"max_health"`
}

// EnemyData represents enemy configuration in JSON
type EnemyData struct {
	ID       string      `json:"id"`
	Position Vector3Data `json:"position"`
	Health   float32     `json:"health"`
	Speed    float32     `json:"speed"`
}

// ObstacleData represents obstacle configuration in JSON
type ObstacleData struct {
	ID       string      `json:"id"`
	Type     string      `json:"type"` // "box" or "cylinder"
	Position Vector3Data `json:"position"`
	Size     Vector3Data `json:"size,omitempty"`     // For box obstacles
	Radius   float32     `json:"radius,omitempty"`   // For cylinder obstacles
	Height   float32     `json:"height,omitempty"`   // For cylinder obstacles
	Color    string      `json:"color"`              // "brown", "green", etc.
}

// HealthPickupData represents health pickup configuration in JSON
type HealthPickupData struct {
	ID         string      `json:"id"`
	Position   Vector3Data `json:"position"`
	HealAmount float32     `json:"heal_amount"`
	Radius     float32     `json:"radius"`
}

// SceneData represents the complete scene configuration
type SceneData struct {
	Metadata struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
	} `json:"metadata"`

	Player PlayerData `json:"player"`

	Entities struct {
		Enemies       []EnemyData        `json:"enemies"`
		Obstacles     []ObstacleData     `json:"obstacles"`
		HealthPickups []HealthPickupData `json:"health_pickups"`
	} `json:"entities"`
}

// LoadSceneFromJSON loads scene data from a JSON file
func LoadSceneFromJSON(filename string) (*SceneData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var sceneData SceneData
	err = json.Unmarshal(data, &sceneData)
	if err != nil {
		return nil, err
	}

	return &sceneData, nil
}

// SaveSceneToJSON saves scene data to a JSON file
func SaveSceneToJSON(sceneData *SceneData, filename string) error {
	data, err := json.MarshalIndent(sceneData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ParseColor converts color string to rl.Color
func ParseColor(colorStr string) rl.Color {
	switch colorStr {
	case "brown":
		return rl.Brown
	case "green":
		return rl.Green
	case "darkgreen":
		return rl.DarkGreen
	case "red":
		return rl.Red
	case "blue":
		return rl.Blue
	case "yellow":
		return rl.Yellow
	case "white":
		return rl.White
	case "gray":
		return rl.Gray
	case "black":
		return rl.Black
	default:
		return rl.Gray // Default color
	}
}