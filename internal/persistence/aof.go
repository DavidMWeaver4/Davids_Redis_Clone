package persistence

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type FsyncPolicy string

// AOF Policy defaults to FsyncEverySec
const (
	FsyncAlways   FsyncPolicy = "always"
	FsyncEverySec FsyncPolicy = "everysec"
	FsyncNo       FsyncPolicy = "no"
)
const aofBaseDir = "data"

var (
	ErrInvalidFsync = errors.New("invalid fsync policy")
	ErrInvalidPath  = errors.New("invalid file path")
	ErrDirTrav      = errors.New("error directory traversal")
	ErrAOFClosed    = errors.New("AOF is closed")
)

type AOFConfig struct {
	Path        string
	FsyncPolicy FsyncPolicy
}
type AOF struct {
	mu     sync.Mutex
	file   *os.File
	policy FsyncPolicy
}

func NewAOF(config AOFConfig) (*AOF, error) {
	var aof AOF
	if config.Path == "" {
		return nil, ErrInvalidPath
	}
	if config.FsyncPolicy == "" {
		config.FsyncPolicy = FsyncEverySec
	}
	switch config.FsyncPolicy {
	case FsyncAlways, FsyncEverySec, FsyncNo:
		aof.policy = config.FsyncPolicy
	default:
		return nil, ErrInvalidFsync
	}
	fp, err := validatePath(config.Path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	aof.file = file
	return &aof, nil
}
func (a *AOF) Append(command protocol.Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file == nil {
		return ErrAOFClosed
	}
	err := protocol.Write(a.file, command)
	if err != nil {
		return err
	}
	if a.policy == FsyncAlways {
		return a.file.Sync()
	}
	return nil
}
func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file == nil {
		return nil
	}

	if err := a.file.Sync(); err != nil {
		return err
	}

	err := a.file.Close()
	a.file = nil
	return err
}
func validatePath(fp string) (string, error) {

	if strings.TrimSpace(fp) == "" {
		return "", ErrInvalidPath
	}
	if strings.ContainsRune(fp, '\x00') {
		return "", ErrInvalidPath
	}
	cleanedInput := filepath.Clean(fp)
	if !filepath.IsLocal(cleanedInput) {
		return "", ErrDirTrav
	}
	return filepath.Join(aofBaseDir, cleanedInput), nil
}
