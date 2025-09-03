package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ProjectDatabase manages project data using file-based storage
type ProjectDatabase struct {
	projectDir string
	nextID     map[string]int
}

// ProjectFile represents a file in the project
type ProjectFile struct {
	ID         int       `json:"id"`
	Path       string    `json:"path"`
	Content    string    `json:"content"`
	Hash       string    `json:"hash"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	IsActive   bool      `json:"is_active"`
}

// EditHistory tracks all edits made to files
type EditHistory struct {
	ID        int       `json:"id"`
	FileID    int       `json:"file_id"`
	OldHash   string    `json:"old_hash"`
	NewHash   string    `json:"new_hash"`
	Diff      string    `json:"diff"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
}

// Checkpoint represents a project snapshot
type Checkpoint struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	FileCount   int       `json:"file_count"`
}

// MCPIntegration represents an MCP integration
type MCPIntegration struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	IsActive  bool                   `json:"is_active"`
}

// MultiplayerAction represents a multiplayer action
type MultiplayerAction struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Action    string    `json:"action"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// NewProjectDatabase creates a new database connection (file-based)
func NewProjectDatabase(dbPath string) (*ProjectDatabase, error) {
	projectDir := filepath.Dir(dbPath)
	if projectDir == "." {
		projectDir = "eve_project_data"
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	pdb := &ProjectDatabase{
		projectDir: projectDir,
		nextID:     make(map[string]int),
	}

	// Initialize nextID counters
	pdb.loadNextIDs()

	return pdb, nil
}

// Close closes the database connection (no-op for file-based storage)
func (pdb *ProjectDatabase) Close() error {
	pdb.saveNextIDs()
	return nil
}

// loadNextIDs loads the next ID counters from file
func (pdb *ProjectDatabase) loadNextIDs() {
	idFile := filepath.Join(pdb.projectDir, "next_ids.json")
	if data, err := ioutil.ReadFile(idFile); err == nil {
		json.Unmarshal(data, &pdb.nextID)
	}
}

// saveNextIDs saves the next ID counters to file
func (pdb *ProjectDatabase) saveNextIDs() {
	idFile := filepath.Join(pdb.projectDir, "next_ids.json")
	if data, err := json.MarshalIndent(pdb.nextID, "", "  "); err == nil {
		ioutil.WriteFile(idFile, data, 0644)
	}
}

// getNextID returns the next ID for a table
func (pdb *ProjectDatabase) getNextID(table string) int {
	id := pdb.nextID[table]
	pdb.nextID[table] = id + 1
	return id + 1
}

// SaveFile saves a file to the project
func (pdb *ProjectDatabase) SaveFile(path, content string) error {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(content)))

	// Check if file exists
	existing, err := pdb.GetFile(path)
	if err == nil && existing.IsActive {
		// Update existing file
		existing.Content = content
		existing.Hash = hash
		existing.Version++
		existing.ModifiedAt = time.Now()
		return pdb.saveFileRecord(existing)
	}

	// Create new file
	file := &ProjectFile{
		ID:         pdb.getNextID("files"),
		Path:       path,
		Content:    content,
		Hash:       hash,
		Version:    1,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		IsActive:   true,
	}

	return pdb.saveFileRecord(file)
}

// saveFileRecord saves a file record to disk
func (pdb *ProjectDatabase) saveFileRecord(file *ProjectFile) error {
	dir := filepath.Join(pdb.projectDir, "files")
	os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, fmt.Sprintf("%d.json", file.ID))
	data, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// GetFile retrieves a file from the project
func (pdb *ProjectDatabase) GetFile(path string) (*ProjectFile, error) {
	dir := filepath.Join(pdb.projectDir, "files")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var f ProjectFile
			if err := json.Unmarshal(data, &f); err != nil {
				continue
			}

			if f.Path == path && f.IsActive {
				return &f, nil
			}
		}
	}

	return nil, fmt.Errorf("file not found")
}

// ListFiles returns all active files in the project
func (pdb *ProjectDatabase) ListFiles() ([]*ProjectFile, error) {
	dir := filepath.Join(pdb.projectDir, "files")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []*ProjectFile
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var f ProjectFile
			if err := json.Unmarshal(data, &f); err != nil {
				continue
			}

			if f.IsActive {
				result = append(result, &f)
			}
		}
	}

	return result, nil
}

// CreateCheckpoint creates a new checkpoint
func (pdb *ProjectDatabase) CreateCheckpoint(name, description string) (*Checkpoint, error) {
	files, err := pdb.ListFiles()
	if err != nil {
		return nil, err
	}

	checkpoint := &Checkpoint{
		ID:          pdb.getNextID("checkpoints"),
		Name:        name,
		Description: description,
		Timestamp:   time.Now(),
		FileCount:   len(files),
	}

	return checkpoint, pdb.saveCheckpointRecord(checkpoint)
}

// saveCheckpointRecord saves a checkpoint record
func (pdb *ProjectDatabase) saveCheckpointRecord(checkpoint *Checkpoint) error {
	dir := filepath.Join(pdb.projectDir, "checkpoints")
	os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, fmt.Sprintf("%d.json", checkpoint.ID))
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// ListCheckpoints returns all checkpoints
func (pdb *ProjectDatabase) ListCheckpoints() ([]*Checkpoint, error) {
	dir := filepath.Join(pdb.projectDir, "checkpoints")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []*Checkpoint
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var c Checkpoint
			if err := json.Unmarshal(data, &c); err != nil {
				continue
			}

			result = append(result, &c)
		}
	}

	// Sort by timestamp (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	return result, nil
}

// AddMCPIntegration adds a new MCP integration
func (pdb *ProjectDatabase) AddMCPIntegration(name, mcpType string, config map[string]interface{}) error {
	integration := &MCPIntegration{
		ID:        pdb.getNextID("mcp"),
		Name:      name,
		Type:      mcpType,
		Config:    config,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	return pdb.saveMCPRecord(integration)
}

// saveMCPRecord saves an MCP integration record
func (pdb *ProjectDatabase) saveMCPRecord(integration *MCPIntegration) error {
	dir := filepath.Join(pdb.projectDir, "mcp")
	os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, fmt.Sprintf("%d.json", integration.ID))
	data, err := json.MarshalIndent(integration, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// GetMCPIntegrations returns all MCP integrations
func (pdb *ProjectDatabase) GetMCPIntegrations() ([]*MCPIntegration, error) {
	dir := filepath.Join(pdb.projectDir, "mcp")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []*MCPIntegration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var m MCPIntegration
			if err := json.Unmarshal(data, &m); err != nil {
				continue
			}

			if m.IsActive {
				result = append(result, &m)
			}
		}
	}

	return result, nil
}

// RecordMultiplayerAction records a multiplayer action
func (pdb *ProjectDatabase) RecordMultiplayerAction(user, action, data string) error {
	mpAction := &MultiplayerAction{
		ID:        pdb.getNextID("multiplayer"),
		User:      user,
		Action:    action,
		Data:      data,
		Timestamp: time.Now(),
	}

	return pdb.saveMultiplayerRecord(mpAction)
}

// saveMultiplayerRecord saves a multiplayer action record
func (pdb *ProjectDatabase) saveMultiplayerRecord(action *MultiplayerAction) error {
	dir := filepath.Join(pdb.projectDir, "multiplayer")
	os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, fmt.Sprintf("%d.json", action.ID))
	data, err := json.MarshalIndent(action, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

// GetMultiplayerHistory returns multiplayer action history
func (pdb *ProjectDatabase) GetMultiplayerHistory(limit int) ([]*MultiplayerAction, error) {
	dir := filepath.Join(pdb.projectDir, "multiplayer")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []*MultiplayerAction
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			filePath := filepath.Join(dir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}

			var a MultiplayerAction
			if err := json.Unmarshal(data, &a); err != nil {
				continue
			}

			result = append(result, &a)
		}
	}

	// Sort by timestamp (newest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	// Apply limit
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// BackupProject creates a backup of the project
func (pdb *ProjectDatabase) BackupProject(backupPath string) error {
	// For file-based storage, we can just copy the directory
	return fmt.Errorf("backup not implemented for file-based storage")
}

// RestoreProject restores a project from backup
func (pdb *ProjectDatabase) RestoreProject(backupPath string) error {
	return fmt.Errorf("restore not implemented for file-based storage")
}
