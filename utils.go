package installer

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var TargetAlreadyExistErr = errors.New("Target file already exists, please clean the install directory and retry")
var EmptyLinkArrayErr = errors.New("Link array is empty")

type StringMap = map[string]string

func osCopy(src, dst string, mode os.FileMode) (err error) {
	srcFd, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFd.Close()
	dstFd, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return
	}
	_, err = io.Copy(dstFd, srcFd)
	if er := dstFd.Close(); err == nil && er != nil {
		err = er
	}
	if err != nil {
		os.Remove(dst)
		return
	}
	return
}

func renameIfNotExist(src, dst string, mode os.FileMode) (err error) {
	if _, e := os.Stat(dst); os.IsNotExist(e) {
		if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return
		}
	} else {
		return &os.LinkError{
			Op:  "rename",
			Old: src,
			New: dst,
			Err: TargetAlreadyExistErr,
		}
	}
	if err = os.Rename(src, dst); err != nil {
		if crossDevice(err) {
			if err = osCopy(src, dst, mode); err != nil {
				return
			}
			os.Remove(src)
			return
		}
		return
	}
	os.Chmod(dst, mode)
	return
}

func safeDownload(reader io.Reader, path string) (err error) {
	var fd *os.File
	if fd, err = os.OpenFile(path+".downloading", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return
	}
	defer os.Remove(fd.Name())
	_, err = io.Copy(fd, reader)
	fd.Close()
	if err != nil {
		return
	}
	if err = renameIfNotExist(fd.Name(), path, 0644); err != nil {
		return
	}
	return nil
}

func lookJavaPath() (string, error) {
	javahome := os.Getenv("JAVA_HOME")
	if len(javahome) > 0 {
		if path, err := exec.LookPath(filepath.Join(javahome, "bin", "java")); err == nil {
			return path, nil
		}
	}
	return exec.LookPath("java")
}

var hashesNewer = map[string]func() hash.Hash{
	"md5":    md5.New,
	"sha1":   sha1.New,
	"sha224": sha256.New224,
	"sha256": sha256.New,
	"sha384": sha512.New384,
	"sha512": sha512.New,
}

func checkHashStream(r io.Reader, hashes StringMap, w io.Writer) (n int64, err error) {
	hashers := make([]hash.Hash, 0, len(hashes))
	expects := make([][2]string, 0, len(hashes))
	for h, sum := range hashes {
		n, ok := hashesNewer[h]
		if ok {
			hashers = append(hashers, n())
			expects = append(expects, [2]string{h, sum})
		}
	}
	writers := make([]io.Writer, len(hashers), len(hashers)+1)
	for i, h := range hashers {
		writers[i] = h
	}
	if w != nil {
		writers = append(writers, w)
	}
	if n, err = io.Copy(io.MultiWriter(writers...), r); err != nil {
		fmt.Printf("err is not nil: %T", err)
		return
	}
	for i, h := range hashers {
		sum := hex.EncodeToString(h.Sum(nil))
		if expect := expects[i]; expect[1] != sum {
			err = &HashErr{
				Hash:   expect[0],
				Sum:    sum,
				Expect: expect[1],
			}
			return
		}
	}
	return
}

func matchHashes(path string, hashes StringMap) (ok bool) {
	fd, err := os.Open(path)
	if err != nil {
		return false
	}
	defer fd.Close()
	_, err = checkHashStream(fd, hashes, nil)
	return err == nil
}

func downloadAnyAndCheckHashes(links []string, path string, hashes StringMap, size int64) (err error) {
	if matchHashes(path, hashes) {
		return
	}
	if len(links) == 0 {
		err = EmptyLinkArrayErr
		return
	}
	for _, l := range links {
		var tmp string
		if tmp, err = DefaultHTTPClient.DownloadTmp(l, "*.downloading", 0644, hashes, size,
			downloadingCallback(l)); err != nil {
			continue
		}
		defer os.Remove(tmp)
		if err = renameIfNotExist(tmp, path, 0644); err != nil {
			return
		}
		break
	}
	return
}

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func formatSize(bytes int64, format string) string {
	b := (float32)(bytes)
	var unit string
	for _, u := range sizeUnits {
		unit = u
		if b <= 1000 {
			break
		}
		b /= 1024
	}
	return fmt.Sprintf(format, b) + unit
}

type DlCallback = func(n int64, size int64)

func downloadingCallback(url string) DlCallback {
	var last time.Time
	var start time.Time
	return func(n int64, size int64) {
		if n == 0 {
			start = time.Now()
			last = start
			loger.Infof("Downloading %q...", url)
			return
		}
		if n == size {
			loger.Infof("Downloaded %q (%v)", url, time.Since(start))
			return
		}
		if time.Since(last) > time.Second {
			last = time.Now()
			loger.Infof("Downloading %q [%s/%s %.2f%%]\n", url, formatSize(n, "%.2f"), formatSize(size, "%.2f"), (float32)(n)/(float32)(size)*100)
		}
	}
}
