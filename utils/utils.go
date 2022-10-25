package utils

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"sort"
	"sync"
)

// Walk walks the filesystem rooted at root, calling walkFn for each file or
// directory in the filesystem, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order.
func Walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := stat(fs, root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

// #nosec G304
// ReadFile returns the raw content of a file
func ReadFile(path string, fs http.FileSystem) ([]byte, error) {
	if fs != nil {
		file, err := fs.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return ioutil.ReadAll(file)
	}
	return ioutil.ReadFile(path)
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := readDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(fis))
	for i := range fis {
		names[i] = fis[i].Name()
	}
	sort.Strings(names)
	return names, nil
}

// walk recursively descends path, calling walkFn.
func walk(fs http.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(fs, path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := pathpkg.Join(path, name)
		fileInfo, err := stat(fs, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// readDir reads the contents of the directory associated with file and
// returns a slice of FileInfo values in directory order.
func readDir(fs http.FileSystem, name string) ([]os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdir(0)
}

// stat returns the FileInfo structure describing file.
func stat(fs http.FileSystem, name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

var copyBufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4096)
	},
}

func SaveFile(fh *multipart.FileHeader, path string) (err error) {
	var (
		f  multipart.File
		ff *os.File
	)
	f, err = fh.Open()
	if err != nil {
		return
	}

	var ok bool
	if ff, ok = f.(*os.File); ok {
		// Windows can't rename files that are opened.
		if err = f.Close(); err != nil {
			return
		}

		// If renaming fails we try the normal copying method.
		// Renaming could fail if the files are on different devices.
		if os.Rename(ff.Name(), path) == nil {
			return nil
		}

		// Reopen f for the code below.
		if f, err = fh.Open(); err != nil {
			return
		}
	}

	defer func() {
		e := f.Close()
		if err == nil {
			err = e
		}
	}()

	if ff, err = os.Create(path); err != nil {
		return
	}
	defer func() {
		e := ff.Close()
		if err == nil {
			err = e
		}
	}()
	_, err = copyZeroAlloc(ff, f)
	return
}

func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	vbuf := copyBufPool.Get()
	buf := vbuf.([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(vbuf)
	return n, err
}
