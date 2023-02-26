package url

import "fmt"

type Option func(*repo) error

func Counter(counter uint64) Option {
	return func(r *repo) error {
		if r.counter != 0 {
			r.logger.Warn(fmt.Sprintf("counter value is overwritten by the value %d", counter))
		}
		r.counter = counter
		return nil
	}
}

func File(path string) Option {
	return func(r *repo) error {
		if len(path) == 0 {
			return ErrEmptyPath
		}

		err := r.fileOpen(path)
		if err != nil {
			r.logger.Error("failed open file path", err)
			return err
		}

		r.fileUse = true

		err = r.fileLoad()
		if err != nil {
			r.logger.Error("failed load url data from file", err)
			return err
		}
		return nil
	}
}
