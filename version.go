
package installer

import (
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func VersionFromString(data string)(v Version, err error){
	v.Major = 0
	v.Minor = 0
	v.Patch = 0
	var i int
	if i = strings.IndexByte(data, '.'); i < 0 {
		v.Major, err = strconv.Atoi(data)
		return
	}
	if v.Major, err = strconv.Atoi(data[:i]); err != nil {
		return
	}
	data = data[i + 1:]
	if i = strings.IndexByte(data, '.'); i < 0 {
		v.Minor, err = strconv.Atoi(data)
		return
	}
	if v.Minor, err = strconv.Atoi(data[:i]); err != nil {
		return
	}
	if v.Patch, err = strconv.Atoi(data[i + 1:]); err != nil {
		return
	}
	return
}

func (v Version)String()(s string){
	s = strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
	return
}

func (v Version)Less(o Version)(bool){
	return v.Major < o.Major || v.Minor < o.Minor || v.Patch < o.Patch
}
