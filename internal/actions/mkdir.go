package actions

import "os"

// Mkdir creates a directory at the given path.
func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
