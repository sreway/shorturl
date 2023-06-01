package cache

// Option describes an option for repository.
type Option func(*repo) error

// File implements an option that sets the file path.
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
