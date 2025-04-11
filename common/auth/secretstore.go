package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SecretStore struct {
	rootDir string
}

func MustMakeSecretStore(dir string, app AppID) SecretStore {
	ss, err := MakeSecretStore(dir, app)
	if err != nil {
		panic(err)
	}
	return ss
}

func MakeSecretStore(dir string, app AppID) (SecretStore, error) {
	if info, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return SecretStore{}, fmt.Errorf("root of secret store not found; consider running `mkdir %q`", dir)
		}
		return SecretStore{}, fmt.Errorf("check root of secret store: %w", err)
	} else if !info.IsDir() {
		return SecretStore{}, fmt.Errorf("root of secret store was not a directory: %v", dir)
	}
	// TODO: add some special file to make sure we recognize this as the right kind of directory.
	appDir := strings.Join([]string{dir, "by-app", app.id}, string(os.PathSeparator))
	if filepath.Clean(appDir) != appDir {
		return SecretStore{}, fmt.Errorf("invalid secrets path")
	}
	if info, err := os.Stat(appDir); err == nil {
		if !info.IsDir() {
			return SecretStore{}, fmt.Errorf("app dir was a file; weird")
		}
		// good to go!
	} else if !os.IsNotExist(err) {
		return SecretStore{}, err
	} else if err := os.MkdirAll(appDir, 0700); err != nil {
		return SecretStore{}, fmt.Errorf("mkdir for app-specific secret store: %w", err)
	}
	return SecretStore{appDir}, nil
}

// #encapsulation for #security
// TODO: enforce uniqueness by putting them all in the same place?
// That still doesn't prevent accidental copy-paste...
type AppID struct {
	id string
}

func NewAppID(id string) AppID {
	if strings.Contains(id, "..") || strings.ContainsAny(id, `/\`) {
		panic("invalid app ID: " + id)
	}
	return AppID{id}
}

type Secret[T any] struct {
	id string
}

func NewSecret[T any](id string) Secret[T] {
	if strings.Contains(id, "..") || strings.ContainsAny(id, `/\`) {
		panic("invalid secret ID: " + id)
	}
	return Secret[T]{id}
}

func (s Secret[T]) secretPath(ss SecretStore) string {
	return filepath.Join(ss.rootDir, s.id)
}

func (s Secret[T]) Read(src SecretStore) (T, error) {
	path := s.secretPath(src)
	data, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("Secret(%v): read: %w", s.id, err)
	}
	var ret T
	if err := json.Unmarshal(data, &ret); err != nil {
		var zero T
		// don't print err - it may leak secret info
		return zero, fmt.Errorf("Secret(%v): unmarshal: failed; is file corrupted?", s.id)
	}
	return ret, nil
}

func (s Secret[T]) Write(dst SecretStore, t T) error {
	path := s.secretPath(dst)
	data, err := json.Marshal(t)
	if err != nil {
		// don't print err - it may leak secret info
		return fmt.Errorf("Secret(%v): marshal: failed", s.id)
	}
	// Using 0400 means that we won't allow overwriting (without further
	// user actions). Seems... okay.
	if err := os.WriteFile(path, data, 0400); err != nil {
		return fmt.Errorf("Secret(%v): write: %w", s.id, err)
	}
	return nil
}
