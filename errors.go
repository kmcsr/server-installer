
package installer

import (
	"fmt"
	"net/http"
)

type UnsupportGameErr struct {
	Game string
}

func (e *UnsupportGameErr)Error()(string){
	return fmt.Sprintf("Unsupport game type %q", e.Game)
}

type NotLocalPathErr struct {
	Path string
}

func (e *NotLocalPathErr)Error()(string){
	return fmt.Sprintf("%q is not a local path", e.Path)
}

type HashErr struct {
	Hash   string
	Sum    string
	Expect string
}

func (e *HashErr)Error()(string){
	return fmt.Sprintf("Unexpect %s hash %s, expect %s", e.Hash, e.Sum, e.Expect)
}

type HttpStatusError struct{
	Code int
}

func (e *HttpStatusError)Error()(string){
	return fmt.Sprintf("Unexpect http status %d %s", e.Code, http.StatusText(e.Code))
}

type ContentLengthNotMatchErr struct {
	ContentLength int64
	Expect        int64
}

func (e *ContentLengthNotMatchErr)Error()(string){
	return fmt.Sprintf("Unexpect content length %d, expect %d", e.ContentLength, e.Expect)
}
