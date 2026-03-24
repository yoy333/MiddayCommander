package vfs

import (
	"io"
	"io/fs"
)

// FS is the base read-only filesystem interface.
type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.StatFS
}

// WritableFS extends FS with mutation operations.
type WritableFS interface {
	FS
	Create(name string) (WriteFile, error)
	Mkdir(name string, perm fs.FileMode) error
	MkdirAll(name string, perm fs.FileMode) error
	Remove(name string) error
	RemoveAll(name string) error
	Rename(oldname, newname string) error
}

// WriteFile extends fs.File with write capability.
type WriteFile interface {
	fs.File
	io.Writer
}
