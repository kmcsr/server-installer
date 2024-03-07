package installer

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HTTPClient struct {
	http.Client

	UserAgent string
}

var DefaultHTTPClient = &HTTPClient{
	Client: http.Client{
		Timeout: time.Second * 10,
	},
	UserAgent: "github.com/kmcsr/server-installer/" + PkgVersion,
}

func (c *HTTPClient) NewRequest(method string, url string, body io.Reader) (req *http.Request, err error) {
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	req.Header.Set("User-Agent", c.UserAgent)
	return
}

func (c *HTTPClient) Do(req *http.Request) (res *http.Response, err error) {
	if ua := req.Header.Get("User-Agent"); ua == "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if res, err = c.Client.Do(req); err != nil {
		return
	}
	return
}

func (c *HTTPClient) Get(url string) (res *http.Response, err error) {
	var req *http.Request
	if req, err = c.NewRequest("GET", url, nil); err != nil {
		return
	}
	return c.Do(req)
}

func (c *HTTPClient) GetJson(url string, obj any) (err error) {
	var req *http.Request
	if req, err = c.NewRequest("GET", url, nil); err != nil {
		return
	}
	req.Header.Set("Accept", "application/json, */*;q=0.1")
	var res *http.Response
	if res, err = c.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &HttpStatusError{
			Code: res.StatusCode,
		}
	}
	return json.NewDecoder(res.Body).Decode(obj)
}

func (c *HTTPClient) GetXml(url string, obj any) (err error) {
	var req *http.Request
	if req, err = c.NewRequest("GET", url, nil); err != nil {
		return
	}
	req.Header.Set("Accept", "application/xml, */*;q=0.1")
	var res *http.Response
	if res, err = c.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &HttpStatusError{
			Code: res.StatusCode,
		}
	}
	return xml.NewDecoder(res.Body).Decode(obj)
}

type progressReader struct {
	Reader io.Reader
	read   int64
	size   int64
	cb     DlCallback
}

func newProgressReader(r io.Reader, size int64, cb DlCallback) *progressReader {
	cb(0, size)
	return &progressReader{
		Reader: r,
		read:   0,
		size:   size,
		cb:     cb,
	}
}

func (r *progressReader) Read(buf []byte) (n int, err error) {
	n, err = r.Reader.Read(buf)
	if n > 0 {
		r.read += (int64)(n)
		r.cb(r.read, r.size)
	}
	return
}

func (c *HTTPClient) DownloadTmp(url string, pattern string, mode os.FileMode, hashes StringMap, size int64, cb DlCallback) (path string, err error) {
	var res *http.Response
	if res, err = c.Get(url); err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = &HttpStatusError{
			Code: res.StatusCode,
		}
		return
	}
	if size < 0 {
		size = res.ContentLength
	} else if res.ContentLength < 0 && res.ContentLength != size {
		err = &ContentLengthNotMatchErr{
			ContentLength: res.ContentLength,
			Expect:        size,
		}
		return
	}
	var r io.Reader = res.Body
	if cb != nil {
		r = newProgressReader(r, size, cb)
	}
	dir, base := filepath.Split(pattern)
	var fd *os.File
	if fd, err = os.CreateTemp(dir, base); err != nil {
		return
	}
	defer func(fd *os.File) {
		fd.Close()
		if err != nil {
			os.Remove(fd.Name())
		}
	}(fd)
	if _, err = checkHashStream(r, hashes, fd); err != nil {
		return
	}
	if mode != 0 {
		if err = fd.Chmod(mode); err != nil {
			return
		}
	}
	path = fd.Name()
	return
}

func (c *HTTPClient) Download(url string, path string, mode os.FileMode, hashes StringMap, size int64, cb DlCallback) (err error) {
	var tmppath string
	tmppath, err = c.DownloadTmp(url, path+".*.downloading", mode, hashes, size, cb)
	if err = renameIfNotExist(tmppath, path, 0644); err != nil {
		return
	}
	return
}

func (c *HTTPClient) Head(url string) (res *http.Response, err error) {
	var req *http.Request
	if req, err = c.NewRequest("HEAD", url, nil); err != nil {
		return
	}
	return c.Do(req)
}

func (c *HTTPClient) Post(url string, contentType string, body io.Reader) (res *http.Response, err error) {
	var req *http.Request
	if req, err = c.NewRequest("POST", url, body); err != nil {
		return
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *HTTPClient) PostForm(url string, form url.Values) (res *http.Response, err error) {
	formStr := form.Encode()
	return c.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader(formStr))
}

func (c *HTTPClient) DownloadDirect(url string, ExactDownloadeName string, cb DlCallback) (installed string, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	filename := filepath.Base(url)
	flags := os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	buf := make([]byte, 16*1024)
	_, err = io.CopyBuffer(f, resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
	}
	cpath, err := os.Getwd()
	if err != nil {
		return
	}
	installed = filepath.Join(cpath, ExactDownloadeName)
	return
}
