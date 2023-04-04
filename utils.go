
package installer

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

var TargetAlreadyExistErr = errors.New("target already exist")

func getHttpJson(url string, obj any)(err error){
	var res *http.Response
	if res, err = http.DefaultClient.Get(url); err != nil {
		return
	}
	defer res.Body.Close()
	var data []byte
	if data, err = io.ReadAll(res.Body); err != nil {
		return
	}
	return json.Unmarshal(data, obj)
}


func renameIfNotExist(src, dst string)(error){
	if _, err := os.Stat(dst); !os.IsNotExist(err) {
		return &os.LinkError{
			Op: "rename",
			Old: src,
			New: dst,
			Err: TargetAlreadyExistErr,
		}
	}
	return os.Rename(src, dst)
}

func safeDownload(reader io.Reader, path string)(err error){
	var fd *os.File
	if fd, err = os.OpenFile(path + ".downloading", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644); err != nil {
		return
	}
	_, err = io.Copy(fd, reader)
	fd.Close()
	if err != nil {
		return
	}
	if err = renameIfNotExist(path + ".downloading", path); err != nil {
		return
	}
	return nil
}

func lookJavaPath()(string, error){
	javahome := os.Getenv("JAVA_HOME")
	if len(javahome) > 0 {
		return exec.LookPath(filepath.Join(javahome, "bin", "java"))
	}
	return exec.LookPath("java")
}
