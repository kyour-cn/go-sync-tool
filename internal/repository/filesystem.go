package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"app/internal/domain"
)

const (
	DefaultConfigDir = "chapar"

	environmentsDir = "envs"
	protoFilesDir   = "protofiles"
	collectionsDir  = "collections"
	requestsDir     = "requests"
	preferencesDir  = "preferences"
)

var _ Repository = &Filesystem{}

type Filesystem struct {
	configDir        string
	baseDir          string
	ActiveWorkspace  *domain.Workspace
	requestPaths     map[string]string
	collectionPaths  map[string]string
	environmentPaths map[string]string
	protoFilePaths   map[string]string
	workspacePaths   map[string]string
}

func NewFilesystem(configDir string, baseDir string) (*Filesystem, error) {
	fs := &Filesystem{
		configDir:        configDir,
		baseDir:          baseDir,
		requestPaths:     make(map[string]string),
		collectionPaths:  make(map[string]string),
		environmentPaths: make(map[string]string),
		protoFilePaths:   make(map[string]string),
		workspacePaths:   make(map[string]string),
	}

	config, err := fs.GetConfig()
	if err != nil {
		return nil, err
	}

	cDir, err := fs.getConfigDir()
	if err != nil {
		return nil, err
	}

	if config.Spec.ActiveWorkspace != nil {
		ws, err := fs.GetWorkspace(filepath.Join(cDir, config.Spec.ActiveWorkspace.Name))
		if err != nil {
			return nil, err
		}
		fs.ActiveWorkspace = ws
	}

	// if there is no active workspace, create default workspace
	if fs.ActiveWorkspace == nil {
		ws := domain.NewDefaultWorkspace()
		defaultPath := filepath.Join(cDir, "default")
		fs.workspacePaths[ws.MetaData.ID] = defaultPath
		if err := fs.updateWorkspace(ws); err != nil {
			return nil, err
		}

		fs.ActiveWorkspace = ws
	}

	return fs, nil
}

func (f *Filesystem) getEntityDirectoryInWorkspace(entityType string) (string, error) {
	dir, err := f.CreateConfigDir()
	if err != nil {
		return "", err
	}

	p := filepath.Join(dir, f.ActiveWorkspace.MetaData.Name, entityType)
	if err := makeDir(p); err != nil {
		return "", err
	}

	return p, nil
}

func (f *Filesystem) getNewProtoFilePath(name string) (*FilePath, error) {
	dir, err := f.getEntityDirectoryInWorkspace(protoFilesDir)
	if err != nil {
		return nil, err
	}

	return getNewFilePath(dir, name), nil
}

func (f *Filesystem) SetActiveWorkspace(workspace *domain.Workspace) error {
	config, err := f.GetConfig()
	if err != nil {
		return err
	}

	f.ActiveWorkspace = workspace
	config.Spec.ActiveWorkspace = &domain.ActiveWorkspace{
		ID:   workspace.MetaData.ID,
		Name: workspace.MetaData.Name,
	}
	return f.UpdateConfig(config)
}

func (f *Filesystem) GetConfig() (*domain.Config, error) {
	dir, err := f.getConfigDir()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(dir, "config.yaml")

	// if config file does not exist, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		config := domain.NewConfig()
		if err := SaveToYaml(filePath, config); err != nil {
			return nil, err
		}

		return config, nil
	}

	return LoadFromYaml[domain.Config](filePath)
}

func (f *Filesystem) UpdateConfig(config *domain.Config) error {
	dir, err := f.getConfigDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(dir, "config.yaml")
	return SaveToYaml(filePath, config)
}

func (f *Filesystem) LoadWorkspaces() ([]*domain.Workspace, error) {
	wdir, err := f.getWorkspacesDir()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(wdir)
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Workspace, 0)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirPath := filepath.Join(wdir, dir.Name())
		if ws, err := f.GetWorkspace(dirPath); err != nil {
			return nil, err
		} else {
			out = append(out, ws)
		}
	}

	return out, nil
}

func (f *Filesystem) GetWorkspace(dirPath string) (*domain.Workspace, error) {
	// if directory is not exist, create it
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(dirPath, "_workspace.yaml")

	// if workspace file does not exist, create it
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ws := domain.NewWorkspace(filepath.Base(dirPath))
		f.workspacePaths[ws.MetaData.ID] = filePath
		if err := SaveToYaml(filePath, ws); err != nil {
			return nil, err
		}

		return ws, nil
	}

	ws, err := LoadFromYaml[domain.Workspace](filePath)
	if err != nil {
		return nil, err
	}

	f.workspacePaths[ws.MetaData.ID] = filePath
	return ws, nil
}

func (f *Filesystem) getWorkspacesDir() (string, error) {
	dir, err := f.CreateConfigDir()
	if err != nil {
		return "", err
	}

	// all folders in the config directory are workspaces
	return dir, nil
}

func (f *Filesystem) updateWorkspace(workspace *domain.Workspace) error {
	filePath, exists := f.workspacePaths[workspace.MetaData.ID]
	if !exists {
		return fmt.Errorf("workspace path not found for %s", workspace.MetaData.ID)
	}

	if err := SaveToYaml(filePath, workspace); err != nil {
		return err
	}

	// Get the directory name
	dirName := filepath.Dir(filePath)
	// Change the directory name to the workspace name
	if workspace.MetaData.Name != filepath.Base(dirName) {
		// replace last part of the path with the new name
		newDirName := filepath.Join(filepath.Dir(dirName), workspace.MetaData.Name)
		if err := os.Rename(dirName, newDirName); err != nil {
			return err
		}
		f.workspacePaths[workspace.MetaData.ID] = filepath.Join(newDirName, "_workspace.yaml")
	}

	return nil
}

func (f *Filesystem) GetNewWorkspaceDir(name string) (*FilePath, error) {
	wDir, err := f.getWorkspacesDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(wDir, name)
	if !dirExist(dir) {
		return &FilePath{
			Path:    dir,
			NewName: name,
		}, nil
	}

	// If the file exists, append a number to the filename.
	for i := 1; ; i++ {
		newDirName := fmt.Sprintf("%s%d", dir, i)
		if !dirExist(newDirName) {
			return &FilePath{
				Path:    newDirName,
				NewName: fmt.Sprintf("%s%d", name, i),
			}, nil
		}
	}
}

func (f *Filesystem) ReadPreferences() (*domain.Preferences, error) {
	dir, err := f.getConfigDir()
	if err != nil {
		return nil, err
	}
	pdir := filepath.Join(dir, f.ActiveWorkspace.MetaData.Name, preferencesDir)
	filePath := filepath.Join(pdir, "preferences.yaml")

	preferences, err := LoadFromYaml[domain.Preferences](filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default preferences if file doesn't exist
			preferences = domain.NewPreferences()
			if err := f.UpdatePreferences(preferences); err != nil {
				return nil, err
			}
			return preferences, nil
		}
		return nil, err
	}
	return preferences, nil
}

func (f *Filesystem) UpdatePreferences(pref *domain.Preferences) error {
	dir, err := f.getConfigDir()
	if err != nil {
		return err
	}

	pdir := filepath.Join(dir, f.ActiveWorkspace.MetaData.Name, preferencesDir)
	if err := makeDir(pdir); err != nil {
		return err
	}

	filePath := filepath.Join(pdir, "preferences.yaml")
	return SaveToYaml[domain.Preferences](filePath, pref)
}

func (f *Filesystem) getNewRequestFilePath(name string) (*FilePath, error) {
	dir, err := f.getEntityDirectoryInWorkspace(requestsDir)
	if err != nil {
		return nil, err
	}
	return getNewFilePath(dir, name), nil
}

func getNewFilePath(dir, name string) *FilePath {
	fileName := filepath.Join(dir, name)
	fName := generateNewFileName(fileName, "yaml")

	return &FilePath{
		Path:    fName,
		NewName: getFileNameWithoutExt(fName),
	}
}

// generateNewFileName takes the original file name and generates a new file name
// with the first possible numeric postfix if the original file exists.
func generateNewFileName(filename, ext string) string {
	if !fileExists(filename + "." + ext) {
		return filename + "." + ext
	}

	// If the file exists, append a number to the filename.
	for i := 1; ; i++ {
		newFilename := fmt.Sprintf("%s%d.%s", filename, i, ext)
		if !fileExists(newFilename) {
			return newFilename
		}
	}
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExist(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func (f *Filesystem) getConfigDir() (string, error) {
	if f.baseDir != "" {
		path := filepath.Join(f.baseDir, f.configDir)
		return path, makeDir(path)
	}

	dir, err := userConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, f.configDir)
	return path, makeDir(path)
}

func (f *Filesystem) CreateConfigDir() (string, error) {
	dir, err := f.getConfigDir()
	if err != nil {
		return "", err
	}

	if err := makeDir(dir); err != nil {
		return "", err
	}

	return dir, nil
}

func makeDir(dir string) error {
	dir = filepath.FromSlash(dir)
	fnMakeDir := func() error { return os.MkdirAll(dir, os.ModePerm) }
	info, err := os.Stat(dir)
	switch {
	case err == nil:
		if info.IsDir() {
			return nil // The directory exists
		} else {
			return fmt.Errorf("path exists but is not a directory: %s", dir)
		}
	case os.IsNotExist(err):
		return fnMakeDir()
	default:
		return err
	}
}

func userConfigDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir += "/lib"

	default: // Unix
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			}
			dir += "/.config"
		}
	}

	return dir, nil
}

func (f *Filesystem) Create(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.Workspace:
		return f.createWorkspace(e)
	default:
		return fmt.Errorf("unsupported entity type: %T", entity)
	}
}

// Helper function to generate unique names
func (f *Filesystem) generateUniqueName(name string) string {
	// Start with the original name
	newName := name
	counter := 1

	// Keep trying new names until we find one that doesn't exist
	for {
		// Check if this name exists in various locations
		exists, err := f.nameExists(newName)
		if err != nil || !exists {
			break
		}

		// If it exists, try the next number
		newName = fmt.Sprintf("%s%d", name, counter)
		counter++
	}

	return newName
}

// Helper function to check if a name exists across different types
func (f *Filesystem) nameExists(name string) (bool, error) {
	// Get all directories we need to check
	reqDir, err := f.getEntityDirectoryInWorkspace(requestsDir)
	if err != nil {
		return false, err
	}

	cDir, err := f.getEntityDirectoryInWorkspace(collectionsDir)
	if err != nil {
		return false, err
	}

	envDir, err := f.getEntityDirectoryInWorkspace(environmentsDir)
	if err != nil {
		return false, err
	}

	// Check in requests directory
	if fileExists(filepath.Join(reqDir, name+".yaml")) {
		return true, nil
	}

	// Check in collections directory
	if dirExist(filepath.Join(cDir, name)) {
		return true, nil
	}

	// Check in environments directory
	if fileExists(filepath.Join(envDir, name+".yaml")) {
		return true, nil
	}

	return false, nil
}

func (f *Filesystem) createWorkspace(workspace *domain.Workspace) error {
	// Get workspaces directory
	workspaceDir, err := f.getWorkspacesDir()
	if err != nil {
		return err
	}

	// Generate directory path internally
	dirPath := filepath.Join(workspaceDir, workspace.MetaData.Name)
	f.workspacePaths[workspace.MetaData.ID] = filepath.Join(dirPath, "_workspace.yaml")

	// Create the workspace directory
	if err := makeDir(dirPath); err != nil {
		return fmt.Errorf("failed to create collection directory: %w", err)
	}

	return f.updateWorkspace(workspace)
}

func (f *Filesystem) Delete(entity interface{}) error {
	deleteFn := func(mp map[string]string, id string) error {
		filePath, exists := mp[id]
		if !exists {
			return fmt.Errorf("collection path not found")
		}
		err := os.RemoveAll(filepath.Dir(filePath))
		if err == nil {
			delete(mp, id)
		}
		return err
	}

	switch e := entity.(type) {
	case *domain.Workspace:
		return deleteFn(f.workspacePaths, e.MetaData.ID)
	default:
		return fmt.Errorf("unsupported entity type: %T", entity)
	}
}

func getFileNameWithoutExt(filePath string) string {
	_, file := filepath.Split(filePath)
	extension := filepath.Ext(file)
	return file[:len(file)-len(extension)]
}

func (f *Filesystem) Update(entity interface{}) error {
	switch e := entity.(type) {
	case *domain.Workspace:
		return f.updateWorkspace(e)
	default:
		return fmt.Errorf("unsupported entity type: %T", entity)
	}
}
