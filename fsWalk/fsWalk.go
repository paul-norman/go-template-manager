package fsWalk

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
)

type DirEntry struct {
	info fs.FileInfo
}
func (d *DirEntry) Name() string				{ return d.info.Name() }
func (d *DirEntry) IsDir() bool					{ return d.info.IsDir() }
func (d *DirEntry) Type() fs.FileMode			{ return d.info.Mode().Type() }
func (d *DirEntry) Info() (fs.FileInfo, error)	{ return d.info, nil }

// Walks the embedded filesystem calling `fn` for each file or directory found.
// All errors that arise visiting files and directories are filtered by `fn`.
// The files are walked in lexical order.
func WalkDir(fileSystem http.FileSystem, directory string, fn fs.WalkDirFunc) error {
	DirEntry, err := dirEntry(fileSystem, directory)

	if err != nil {
		err = fn(directory, nil, err)
	} else {
		err = walkDir(fileSystem, directory, DirEntry, fn)
	}

	if err == fs.SkipDir {
		return nil
	}

	return err
}

// Returns the content of a file that is embedded in the file system.
func ReadFile(filename string, fileSystem http.FileSystem) ([]byte, error) {
	if fileSystem != nil {
		file, err := fileSystem.Open(filename)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		return io.ReadAll(file)
	}

	return os.ReadFile(filename)
}

// Returns the `fs.DirEntry` structure describing the directory.
func dirEntry(fileSystem http.FileSystem, directory string) (fs.DirEntry, error) {
	dir, err := fileSystem.Open(directory)
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	info, err := dir.Stat()
	if err != nil {
		return nil, err
	}

	return &DirEntry{info}, nil
}

// Reads `directory` and returns a sorted list of its contents.
func scanDir(fileSystem http.FileSystem, directory string) ([]string, error) {	
	dir, err := fileSystem.Open(directory)
	if err != nil {
		return nil, err
	}

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(fileInfos))
	for i := range fileInfos {
		names[i] = fileInfos[i].Name()
	}

	sort.Strings(names)

	return names, nil
}

// Recursively descends directories, calling the `fs.WalkDirFunc`: `fn`.
func walkDir(fileSystem http.FileSystem, directory string, info fs.DirEntry, fn fs.WalkDirFunc) error {
	err := fn(directory, info, nil)
	if err != nil {
		if info.IsDir() && err == fs.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := scanDir(fileSystem, directory)
	if err != nil {
		return fn(directory, info, err)
	}

	for _, name := range names {
		filename := path.Join(directory, name)
		DirEntry, err := dirEntry(fileSystem, filename)
		if err != nil {
			if err := fn(filename, DirEntry, err); err != nil && err != fs.SkipDir {
				return err
			}
		} else {
			err = walkDir(fileSystem, filename, DirEntry, fn)
			if err != nil {
				if !DirEntry.IsDir() || err != fs.SkipDir {
					return err
				}
			}
		}
	}

	return nil
}