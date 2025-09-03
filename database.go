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
	ID          int       `json:"id"`
	Path        string    `json:"path"`
	Content     string    `json:"content"`
	Hash        string    `json:"hash"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	IsActive    bool      `json:"is_active"`
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
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	IsActive    bool                   `json:"is_active"`
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

// NewProjectDatabase creates a new database connection (file-based fallback)
func NewProjectDatabase(dbPath string) (*ProjectDatabase, error) {
	// Use directory-based storage instead of SQLite
	projectDir := filepath.Dir(dbPath)
	if projectDir == "." {
		projectDir = "eve_project_data"
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	return &ProjectDatabase{projectDir: projectDir}, nil
}

// Close closes the database connection (no-op for file-based storage)
func (pdb *ProjectDatabase) Close() error {
	return nil
}

// Helper methods for file-based storage
func (pdb *ProjectDatabase) getFilePath(table string, id interface{}) string {
	return filepath.Join(pdb.projectDir, fmt.Sprintf("%s_%v.json", table, id))
}

func (pdb *ProjectDatabase) getTableDir(table string) string {
	dir := filepath.Join(pdb.projectDir, table)
	os.MkdirAll(dir, 0755)
	return dir
}

func (pdb *ProjectDatabase) saveRecord(table string, id interface{}, data interface{}) error {
	filePath := pdb.getFilePath(table, id)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}

func (pdb *ProjectDatabase) loadRecord(table string, id interface{}, data interface{}) error {
	filePath := pdb.getFilePath(table, id)
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, data)
}

func (pdb *ProjectDatabase) listRecords(table string, data interface{}) error {
	dir := pdb.getTableDir(table)
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var records []interface{}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(dir, file.Name())
			jsonData, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}
			var record interface{}
			if err := json.Unmarshal(jsonData, &record); err != nil {
				continue
			}
			records = append(records, record)
		}
	}

	// This is a simplified implementation - in practice you'd need type assertions
	return nil
}

// initTables creates all necessary database tables
func (pdb *ProjectDatabase) initTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS project_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL UNIQUE,
			content TEXT,
			hash TEXT,
			version INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_active BOOLEAN DEFAULT 1
		)`,
		`CREATE TABLE IF NOT EXISTS edit_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			file_id INTEGER,
			old_hash TEXT,
			new_hash TEXT,
			diff TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			user TEXT DEFAULT 'eve',
			FOREIGN KEY (file_id) REFERENCES project_files(id)
		)`,
		`CREATE TABLE IF NOT EXISTS checkpoints (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			file_count INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS checkpoint_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			checkpoint_id INTEGER,
			file_id INTEGER,
			content TEXT,
			hash TEXT,
			FOREIGN KEY (checkpoint_id) REFERENCES checkpoints(id),
			FOREIGN KEY (file_id) REFERENCES project_files(id)
		)`,
		`CREATE TABLE IF NOT EXISTS multiplayer_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL UNIQUE,
			user_id TEXT NOT NULL,
			action TEXT NOT NULL,
			data TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS mcp_integrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			endpoint TEXT NOT NULL,
			auth_token TEXT,
			config TEXT,
			is_active BOOLEAN DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range tables {
		if _, err := pdb.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (pdb *ProjectDatabase) Close() error {
	if pdb.db != nil {
		return pdb.db.Close()
	}
	return nil
}

// SaveFile saves or updates a file in the database
func (pdb *ProjectDatabase) SaveFile(path, content, hash string) error {
	now := time.Now()

	// Check if file exists
	var existingID int
	var existingVersion int
	err := pdb.db.QueryRow("SELECT id, version FROM project_files WHERE path = ? AND is_active = 1", path).Scan(&existingID, &existingVersion)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing file: %w", err)
	}

	if err == sql.ErrNoRows {
		// Insert new file
		_, err = pdb.db.Exec(
			"INSERT INTO project_files (path, content, hash, version, created_at, modified_at) VALUES (?, ?, ?, 1, ?, ?)",
			path, content, hash, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert new file: %w", err)
		}
	} else {
		// Update existing file
		newVersion := existingVersion + 1
		_, err = pdb.db.Exec(
			"UPDATE project_files SET content = ?, hash = ?, version = ?, modified_at = ? WHERE id = ?",
			content, hash, newVersion, now, existingID,
		)
		if err != nil {
			return fmt.Errorf("failed to update file: %w", err)
		}
	}

	return nil
}

// GetFile retrieves a file from the database
func (pdb *ProjectDatabase) GetFile(path string) (*ProjectFile, error) {
	var file ProjectFile
	err := pdb.db.QueryRow(
		"SELECT id, path, content, hash, version, created_at, modified_at FROM project_files WHERE path = ? AND is_active = 1",
		path,
	).Scan(&file.ID, &file.Path, &file.Content, &file.Hash, &file.Version, &file.CreatedAt, &file.ModifiedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return &file, nil
}

// ListFiles returns all active files in the project
func (pdb *ProjectDatabase) ListFiles() ([]ProjectFile, error) {
	rows, err := pdb.db.Query(
		"SELECT id, path, content, hash, version, created_at, modified_at FROM project_files WHERE is_active = 1 ORDER BY path",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}
	defer rows.Close()

	var files []ProjectFile
	for rows.Next() {
		var file ProjectFile
		err := rows.Scan(&file.ID, &file.Path, &file.Content, &file.Hash, &file.Version, &file.CreatedAt, &file.ModifiedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	return files, nil
}

// RecordEdit records an edit in the history
func (pdb *ProjectDatabase) RecordEdit(fileID int, oldHash, newHash, diff string) error {
	_, err := pdb.db.Exec(
		"INSERT INTO edit_history (file_id, old_hash, new_hash, diff) VALUES (?, ?, ?, ?)",
		fileID, oldHash, newHash, diff,
	)
	if err != nil {
		return fmt.Errorf("failed to record edit: %w", err)
	}
	return nil
}

// GetEditHistory returns edit history for a file
func (pdb *ProjectDatabase) GetEditHistory(fileID int) ([]EditHistory, error) {
	rows, err := pdb.db.Query(
		"SELECT id, file_id, old_hash, new_hash, diff, timestamp, user FROM edit_history WHERE file_id = ? ORDER BY timestamp DESC",
		fileID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get edit history: %w", err)
	}
	defer rows.Close()

	var history []EditHistory
	for rows.Next() {
		var edit EditHistory
		err := rows.Scan(&edit.ID, &edit.FileID, &edit.OldHash, &edit.NewHash, &edit.Diff, &edit.Timestamp, &edit.User)
		if err != nil {
			return nil, fmt.Errorf("failed to scan edit: %w", err)
		}
		history = append(history, edit)
	}

	return history, nil
}

// CreateCheckpoint creates a snapshot of the current project state
func (pdb *ProjectDatabase) CreateCheckpoint(name, description string) error {
	// Count active files
	var fileCount int
	err := pdb.db.QueryRow("SELECT COUNT(*) FROM project_files WHERE is_active = 1").Scan(&fileCount)
	if err != nil {
		return fmt.Errorf("failed to count files: %w", err)
	}

	// Insert checkpoint
	result, err := pdb.db.Exec(
		"INSERT INTO checkpoints (name, description, file_count) VALUES (?, ?, ?)",
		name, description, fileCount,
	)
	if err != nil {
		return fmt.Errorf("failed to create checkpoint: %w", err)
	}

	checkpointID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get checkpoint ID: %w", err)
	}

	// Save all current files to checkpoint
	files, err := pdb.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to list files for checkpoint: %w", err)
	}

	for _, file := range files {
		_, err = pdb.db.Exec(
			"INSERT INTO checkpoint_files (checkpoint_id, file_id, content, hash) VALUES (?, ?, ?, ?)",
			checkpointID, file.ID, file.Content, file.Hash,
		)
		if err != nil {
			return fmt.Errorf("failed to save file to checkpoint: %w", err)
		}
	}

	return nil
}

// RestoreCheckpoint restores project to a checkpoint state
func (pdb *ProjectDatabase) RestoreCheckpoint(checkpointID int) error {
	// Get all files from checkpoint
	rows, err := pdb.db.Query(
		"SELECT cf.file_id, cf.content, cf.hash, pf.path FROM checkpoint_files cf JOIN project_files pf ON cf.file_id = pf.id WHERE cf.checkpoint_id = ?",
		checkpointID,
	)
	if err != nil {
		return fmt.Errorf("failed to get checkpoint files: %w", err)
	}
	defer rows.Close()

	// Update files to checkpoint state
	for rows.Next() {
		var fileID int
		var content, hash, path string
		err := rows.Scan(&fileID, &content, &hash, &path)
		if err != nil {
			return fmt.Errorf("failed to scan checkpoint file: %w", err)
		}

		// Update file content
		_, err = pdb.db.Exec(
			"UPDATE project_files SET content = ?, hash = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?",
			content, hash, fileID,
		)
		if err != nil {
			return fmt.Errorf("failed to restore file: %w", err)
		}
	}

	return nil
}

// ListCheckpoints returns all checkpoints
func (pdb *ProjectDatabase) ListCheckpoints() ([]Checkpoint, error) {
	rows, err := pdb.db.Query(
		"SELECT id, name, description, timestamp, file_count FROM checkpoints ORDER BY timestamp DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list checkpoints: %w", err)
	}
	defer rows.Close()

	var checkpoints []Checkpoint
	for rows.Next() {
		var cp Checkpoint
		err := rows.Scan(&cp.ID, &cp.Name, &cp.Description, &cp.Timestamp, &cp.FileCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan checkpoint: %w", err)
		}
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

// AddMCPIntegration adds an MCP server integration
func (pdb *ProjectDatabase) AddMCPIntegration(name, endpoint, authToken, config string) error {
	_, err := pdb.db.Exec(
		"INSERT INTO mcp_integrations (name, endpoint, auth_token, config) VALUES (?, ?, ?, ?)",
		name, endpoint, authToken, config,
	)
	if err != nil {
		return fmt.Errorf("failed to add MCP integration: %w", err)
	}
	return nil
}

// GetMCPIntegrations returns all active MCP integrations
func (pdb *ProjectDatabase) GetMCPIntegrations() ([]map[string]interface{}, error) {
	rows, err := pdb.db.Query(
		"SELECT id, name, endpoint, auth_token, config FROM mcp_integrations WHERE is_active = 1",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP integrations: %w", err)
	}
	defer rows.Close()

	var integrations []map[string]interface{}
	for rows.Next() {
		var id int
		var name, endpoint, authToken, config string
		err := rows.Scan(&id, &name, &endpoint, &authToken, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to scan MCP integration: %w", err)
		}

		integration := map[string]interface{}{
			"id":         id,
			"name":       name,
			"endpoint":   endpoint,
			"auth_token": authToken,
			"config":     config,
		}
		integrations = append(integrations, integration)
	}

	return integrations, nil
}

// RecordMultiplayerAction records a multiplayer session action
func (pdb *ProjectDatabase) RecordMultiplayerAction(sessionID, userID, action, data string) error {
	_, err := pdb.db.Exec(
		"INSERT INTO multiplayer_sessions (session_id, user_id, action, data) VALUES (?, ?, ?, ?)",
		sessionID, userID, action, data,
	)
	if err != nil {
		return fmt.Errorf("failed to record multiplayer action: %w", err)
	}
	return nil
}

// GetMultiplayerHistory returns multiplayer session history
func (pdb *ProjectDatabase) GetMultiplayerHistory(sessionID string) ([]map[string]interface{}, error) {
	rows, err := pdb.db.Query(
		"SELECT user_id, action, data, timestamp FROM multiplayer_sessions WHERE session_id = ? ORDER BY timestamp DESC",
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get multiplayer history: %w", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var userID, action, data string
		var timestamp time.Time
		err := rows.Scan(&userID, &action, &data, &timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan multiplayer action: %w", err)
		}

		actionData := map[string]interface{}{
			"user_id":   userID,
			"action":    action,
			"data":      data,
			"timestamp": timestamp,
		}
		history = append(history, actionData)
	}

	return history, nil
}

// BackupProject creates a full backup of the project
func (pdb *ProjectDatabase) BackupProject(backupPath string) error {
	// Export all data as JSON
	data := map[string]interface{}{}

	// Get all files
	files, err := pdb.ListFiles()
	if err != nil {
		return fmt.Errorf("failed to get files for backup: %w", err)
	}
	data["files"] = files

	// Get all checkpoints
	checkpoints, err := pdb.ListCheckpoints()
	if err != nil {
		return fmt.Errorf("failed to get checkpoints for backup: %w", err)
	}
	data["checkpoints"] = checkpoints

	// Get MCP integrations
	mcpIntegrations, err := pdb.GetMCPIntegrations()
	if err != nil {
		return fmt.Errorf("failed to get MCP integrations for backup: %w", err)
	}
	data["mcp_integrations"] = mcpIntegrations

	// Write to file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup data: %w", err)
	}

	if err := os.WriteFile(backupPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// RestoreProject restores project from backup
func (pdb *ProjectDatabase) RestoreProject(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var backupData map[string]interface{}
	if err := json.Unmarshal(data, &backupData); err != nil {
		return fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	// Restore files
	if filesData, ok := backupData["files"]; ok {
		filesJSON, _ := json.Marshal(filesData)
		var files []ProjectFile
		if err := json.Unmarshal(filesJSON, &files); err == nil {
			for _, file := range files {
				pdb.SaveFile(file.Path, file.Content, file.Hash)
			}
		}
	}

	// Restore checkpoints
	if checkpointsData, ok := backupData["checkpoints"]; ok {
		checkpointsJSON, _ := json.Marshal(checkpointsData)
		var checkpoints []Checkpoint
		if err := json.Unmarshal(checkpointsJSON, &checkpoints); err == nil {
			for _, cp := range checkpoints {
				pdb.CreateCheckpoint(cp.Name, cp.Description)
			}
		}
	}

	return nil
}
