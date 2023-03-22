package cache

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type fs struct {
	Data map[uuid.UUID]storageURL `json:"data"`
}

func (r *repo) fileOpen(path string) error {
	flag := os.O_RDWR | os.O_CREATE
	file, err := os.OpenFile(path, flag, 0o644)
	if err != nil {
		return err
	}
	r.file = file
	return nil
}

func (r *repo) fileClose() error {
	err := r.fileStore()
	if err != nil {
		r.logger.Error("failed store data to file", err)
		return err
	}
	r.logger.Info("trigger close repository file")
	return r.file.Close()
}

func (r *repo) fileLoad() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	store := new(fs)

	if err := json.NewDecoder(r.file).Decode(store); err != nil {
		return err
	}

	r.data = store.Data
	r.logger.Info("success load url data from file")

	return nil
}

func (r *repo) fileStore() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = r.file.Seek(0, 0)
	if err != nil {
		return err
	}

	store := new(fs)
	store.Data = r.data

	if err = json.NewEncoder(r.file).Encode(store); err != nil {
		return err
	}

	r.logger.Info("success store url data to file")

	return nil
}
