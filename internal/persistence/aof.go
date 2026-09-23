package persistence

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

	stop      chan struct{}
	wg        sync.WaitGroup
	closeOnce sync.Once
	errors    chan error
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
	aof.errors = make(chan error, 1)

	if aof.policy == FsyncEverySec {
		aof.stop = make(chan struct{})

		ticker := time.NewTicker(1 * time.Second)
		aof.wg.Add(1)

		go func() {
			defer ticker.Stop()
			defer aof.wg.Done()

			for {
				select {
				case <-aof.stop:
					return

				case <-ticker.C:
					aof.mu.Lock()

					if aof.file != nil {
						syncErr := aof.file.Sync()
						if syncErr != nil {
							select {
							case aof.errors <- syncErr:
							default:
							}
						}
					}
					aof.mu.Unlock()
				}
			}
		}()
	}
	return &aof, nil
}
func (a *AOF) Errors() <-chan error {
	return a.errors
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
	var closeErr error

	a.closeOnce.Do(func() {
		if a.policy == FsyncEverySec {
			close(a.stop)
			a.wg.Wait()
		}

		a.mu.Lock()
		defer a.mu.Unlock()

		if a.file == nil {
			close(a.errors)
			return
		}

		syncErr := a.file.Sync()
		fileCloseErr := a.file.Close()
		a.file = nil

		if syncErr != nil {
			closeErr = syncErr
		} else if fileCloseErr != nil {
			closeErr = fileCloseErr
		}

		close(a.errors)
	})

	return closeErr
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
