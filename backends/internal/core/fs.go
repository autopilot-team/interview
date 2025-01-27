package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ZeroFS is an empty FS implementation
var ZeroFS FS

// FS is an interface that abstracts filesystem operations.
// It can be implemented by both real and embedded filesystems.
type FS interface {
	// Open opens the named file for reading
	Open(name string) (fs.File, error)

	// ReadDir reads the contents of the directory and returns a slice of directory entries
	ReadDir(name string) ([]fs.DirEntry, error)

	// ReadFile reads the named file and returns its contents
	ReadFile(name string) ([]byte, error)
}

// LocalFS implements FS interface for the local filesystem
type LocalFS struct {
	root string
}

// NewLocalFS creates a new LocalFS instance with the given root directory.
// The root path is resolved relative to the project root.
func NewLocalFS(root string) (*LocalFS, error) {
	// Find project root
	projectRoot, err := FindProjectRoot()
	if err != nil {
		return nil, err
	}

	// Clean the input path and remove any leading ".." segments
	cleanPath := filepath.Clean(root)
	if strings.HasPrefix(cleanPath, "..") {
		cleanPath = strings.TrimPrefix(strings.TrimPrefix(cleanPath, ".."), string(filepath.Separator))
	}

	// Join with project root
	fullPath := filepath.Join(projectRoot, cleanPath)

	// Verify the directory exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path %s: %w", fullPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path %s is not a directory", fullPath)
	}

	return &LocalFS{root: fullPath}, nil
}

// resolvePath joins and cleans the root path with the given name
func (l *LocalFS) resolvePath(name string) string {
	return filepath.Join(l.root, filepath.Clean(name))
}

// Open implements FS.Open
func (l *LocalFS) Open(name string) (fs.File, error) {
	return os.Open(l.resolvePath(name))
}

// ReadDir implements FS.ReadDir
func (l *LocalFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(l.resolvePath(name))
}

// ReadFile implements FS.ReadFile
func (l *LocalFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(l.resolvePath(name))
}

// FindProjectRoot walks up the directory tree until it finds the project root
// (identified by the presence of go.mod or .git)
func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	dir := cwd
	for {
		// Check for markers of project root
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (no go.mod or .git found)")
		}
		dir = parent
	}
}
