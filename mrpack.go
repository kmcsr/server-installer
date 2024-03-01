package installer

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const currentMrpackVersion = 1

var supportVersions = []int{currentMrpackVersion}

const (
	MrpackEnvRequired    = "required"
	MrpackEnvOptional    = "optional"
	MrpackEnvUnsupported = "unsupported"
)

type MrpackVerisonErr struct {
	Version  int
	Supports []int
}

func (e *MrpackVerisonErr) Error() string {
	return fmt.Sprintf("Unsupport mrpack format version %d, supports %v", e.Version, e.Supports)
}

type (
	Mrpack struct {
		r *zip.ReadCloser

		MrpackMeta

		overrides       []*zip.File
		clientOverrides []*zip.File
		serverOverrides []*zip.File
	}
	MrpackMeta struct {
		FormatVersion int    `json:"formatVersion"`
		Game          string `json:"game"`

		VersionId string `json:"versionId"`
		Name      string `json:"name"`
		Summary   string `json:"summary,optional"`

		Files []MrpackFileMeta `json:"files"`
		Deps  StringMap        `json:"dependencies"`
	}
	MrpackFileMeta struct {
		Path      string    `json:"path"`
		Hashes    StringMap `json:"hashes"`
		Env       StringMap `json:"env"`
		Downloads []string  `json:"downloads"`
		Size      int64     `json:"fileSize"`
	}
)

func OpenMrpack(filename string) (pack *Mrpack, err error) {
	pack = new(Mrpack)
	pack.r, err = zip.OpenReader(filename)
	if err != nil {
		return
	}
	if err = pack.decodeIndex(); err != nil {
		return
	}
	for _, f := range pack.r.File {
		if strings.HasPrefix(f.Name, "overrides/") {
			pack.overrides = append(pack.overrides, f)
		} else if strings.HasPrefix(f.Name, "client-overrides/") {
			pack.clientOverrides = append(pack.clientOverrides, f)
		} else if strings.HasPrefix(f.Name, "server-overrides/") {
			pack.serverOverrides = append(pack.serverOverrides, f)
		}
	}
	return
}

func (p *Mrpack) Close() (err error) {
	return p.r.Close()
}

func (p *Mrpack) decodeIndex() (err error) {
	var indexFd fs.File
	indexFd, err = p.r.Open("modrinth.index.json")
	if err != nil {
		return
	}
	defer indexFd.Close()
	if err = json.NewDecoder(indexFd).Decode(&p.MrpackMeta); err != nil {
		return
	}
	if p.FormatVersion != currentMrpackVersion {
		return &MrpackVerisonErr{
			Version:  p.FormatVersion,
			Supports: supportVersions,
		}
	}
	return
}

type MrpackOptionalChecker func(f MrpackFileMeta) bool

func (p *Mrpack) InstallClient(target string) (err error) {
	return p.InstallClientWithOptional(target, func(f MrpackFileMeta) bool { return true })
}

func (p *Mrpack) installWithEnv(env string, target string, optionalChecker MrpackOptionalChecker) (err error) {
	loger.Infof("Installing [%s]modpack %s(%s) to %q ...", p.Game, p.Name, p.VersionId, target)
	if len(p.Summary) > 0 {
		loger.Infof("  Summary: %s", p.Summary)
	}
	if p.Game != "minecraft" {
		return &UnsupportGameErr{
			Game: p.Game,
		}
	}
	var wg sync.WaitGroup
	for _, f := range p.Files {
		required := true
		if f.Env != nil {
			env := f.Env[env]
			if env == MrpackEnvUnsupported {
				continue
			}
			if env == MrpackEnvOptional {
				if !optionalChecker(f) {
					continue
				}
				required = false
			}
		}
		wg.Add(1)
		go func(f MrpackFileMeta) {
			defer wg.Done()
			if !filepath.IsLocal(f.Path) {
				err = &NotLocalPathErr{f.Path}
				goto checkErr
			}
			if err = downloadAnyAndCheckHashes(f.Downloads, filepath.Join(target, f.Path), f.Hashes, f.Size); err != nil {
				goto checkErr
			}
		checkErr:
			if err != nil {
				if required {
					return
				}
				loger.Warnf("Skipped to install optional mod %q due %v", f.Path, err)
				err = nil
			}
		}(f)
	}
	wg.Wait()
	return
}

func (p *Mrpack) InstallClientWithOptional(target string, optionalChecker MrpackOptionalChecker) (err error) {
	if err = p.installWithEnv("client", target, optionalChecker); err != nil {
		return
	}
	return p.OverrideClient(target)
}

func (p *Mrpack) InstallServer(target string) (err error) {
	return p.InstallServerWithOptional(target, func(f MrpackFileMeta) bool { return true })
}

func (p *Mrpack) InstallServerWithOptional(target string, optionalChecker MrpackOptionalChecker) (err error) {
	if err = p.installWithEnv("server", target, optionalChecker); err != nil {
		return
	}
	return p.OverrideServer(target)
}

func trimLeftDir(path string) string {
	i := strings.IndexByte(path, '/')
	if i < 0 {
		return ""
	}
	return path[i+1:]
}

func (p *Mrpack) override(target string, f *zip.File) (err error) {
	name := trimLeftDir(f.Name)
	if len(name) == 0 {
		return
	}
	path := filepath.Join(target, name)
	if f.FileInfo().IsDir() {
		return os.MkdirAll(path, f.Mode()|0111)
	}
	var (
		r  io.ReadCloser
		fd *os.File
	)
	if r, err = f.Open(); err != nil {
		return
	}
	defer r.Close()
	if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	if fd, err = os.OpenFile(path,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
		return
	}
	if _, err = io.Copy(fd, r); err != nil {
		return
	}
	return
}

func (p *Mrpack) overrideGlobal(target string) (err error) {
	for _, f := range p.overrides {
		if err = p.override(target, f); err != nil {
			return
		}
	}
	return
}

func (p *Mrpack) OverrideClient(target string) (err error) {
	if err = p.overrideGlobal(target); err != nil {
		return
	}
	for _, f := range p.clientOverrides {
		if err = p.override(target, f); err != nil {
			return
		}
	}
	return
}

func (p *Mrpack) OverrideServer(target string) (err error) {
	if err = p.overrideGlobal(target); err != nil {
		return
	}
	for _, f := range p.serverOverrides {
		if err = p.override(target, f); err != nil {
			return
		}
	}
	return
}
