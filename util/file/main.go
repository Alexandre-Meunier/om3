package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// IsProtected returns true if the file is too critical to alter or remove
func IsProtected(fpath string) bool {
	if fpath == "" {
		return true
	}
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return false
	}
	switch absPath {
	case "/", "/bin", "/boot", "/dev", "/dev/pts", "/dev/shm", "/home", "/opt", "/proc", "/sys", "/tmp", "/usr", "/var":
		return true
	default:
		return false
	}
}

// Exists returns true if the file path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ExistsNotDir returns true if the file path exists and is not a directory.
func ExistsNotDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ExistsAndDir returns true if the file path exists and is a directory.
func ExistsAndDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// ExistsAndRegular returns true if the file path exists and is a regular file.
func ExistsAndRegular(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode().IsRegular()
}

// ExistsAndSymlink returns true if the file path exists and is a symbolic link.
func ExistsAndSymlink(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

//
// Copy copies the file content from src file path to dst file path.
// If dst does not exist, it is created.
//
func Copy(src string, dst string) (err error) {
	var (
		r *os.File
		w *os.File
	)
	if r, err = os.Open(src); err != nil {
		return err
	}
	defer r.Close()
	if w, err = os.Create(dst); err != nil {
		return err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	r.Close()
	w.Close()
	return nil
}

// Copy2 is Identical to Copy() except that Copy2() also attempts to preserve file metadata.
func Copy2(src string, dst string) (err error) {
	if err := Copy(src, dst); err != nil {
		return err
	}
	if err := CopyMeta(src, dst); err != nil {
		return err
	}
	return nil
}

// CopyMeta clones the uid, gid, mtime and mode from src to dst
func CopyMeta(src string, dst string) (err error) {
	if info, err := os.Stat(src); err != nil {
		return err
	} else {
		if err := os.Chmod(dst, info.Mode()); err != nil {
			return err
		}
		if err := os.Chtimes(dst, time.Unix(0, 0), info.ModTime()); err != nil {
			return err
		}
	}
	return CopyOwnership(src, dst)
}

func CopyOwnership(src string, dst string) (err error) {
	if uid, gid, err := Ownership(src); err != nil {
		return err
	} else if err := os.Chown(dst, uid, gid); err != nil {
		return err
	}
	return nil
}

//
// ReadAll reads and return all content of the file at src
//
func ReadAll(src string) ([]byte, error) {
	var (
		r   *os.File
		err error
	)
	if r, err = os.Open(src); err != nil {
		return []byte{}, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

//
// ModTime returns the file modification time or a zero time.
//
func ModTime(p string) (mtime time.Time) {
	fi, err := os.Stat(p)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	return
}

//
// IsPerm returns true if the file current permissions are the same as the target.
//
func IsPerm(p string, perm os.FileMode) (bool, error) {
	currentMode, err := Mode(p)
	if err != nil {
		return false, err
	}
	return currentMode.Perm() == perm, nil
}

//
// IsMode returns true if the file current mode is the same as the target mode.
//
func IsMode(p string, mode os.FileMode) (bool, error) {
	currentMode, err := Mode(p)
	if err != nil {
		return false, err
	}
	return currentMode == mode, nil
}

//
// Mode returns the FileMode of the file.
//
func Mode(p string) (os.FileMode, error) {
	fileInfo, err := os.Lstat(p)
	if err != nil {
		return 0, err
	}
	currentMode := fileInfo.Mode()
	return currentMode, nil
}

//
// Ownership return the uid and gid owning the file
//
func Ownership(p string) (uid, gid int, err error) {
	fileInfo, err := os.Lstat(p)
	if err != nil {
		return -1, -1, err
	}
	if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
		return int(stat.Uid), int(stat.Gid), nil
	}
	// unsupported
	return -1, -1, nil
}
