package repository

import (
	"os"

	"gopkg.in/yaml.v2"

	"app/internal/domain"
)

// Repository defines the main storage interface
type Repository interface {
	// Create operations for entities
	Create(entity interface{}) error
	// Update operations for entities
	Update(entity interface{}) error
	// Delete operations for entities
	Delete(entity interface{}) error

	// LoadWorkspaces loads all workspaces
	LoadWorkspaces() ([]*domain.Workspace, error)
	// GetWorkspace gets a workspace by id
	GetWorkspace(id string) (*domain.Workspace, error)
	// SetActiveWorkspace sets the active workspace
	SetActiveWorkspace(workspace *domain.Workspace) error

	// GetConfig gets the config
	GetConfig() (*domain.Config, error)
	// UpdateConfig updates the config
	UpdateConfig(config *domain.Config) error

	// ReadPreferences reads the preferences
	ReadPreferences() (*domain.Preferences, error)
	// UpdatePreferences updates the preferences
	UpdatePreferences(pref *domain.Preferences) error
}

type FilePath struct {
	Path    string
	NewName string
}

func LoadFromYaml[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	env := new(T)
	if err := yaml.Unmarshal(data, env); err != nil {
		return nil, err
	}
	return env, nil
}

func SaveToYaml[T any](filename string, data *T) error {
	out, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, out, 0644); err != nil {
		return err
	}
	return nil
}
