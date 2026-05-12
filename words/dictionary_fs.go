package words

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type dictionaryFs interface {
	Open(wordLen int) (io.ReadCloser, error)
}

type embeddedDictionaryFs struct {
	embed embed.FS
}

func (e *embeddedDictionaryFs) Open(wordLen int) (io.ReadCloser, error) {
	return e.embed.Open(useDictionary + "-" + strconv.Itoa(wordLen) + "-letters.txt")
}

func newExternalDictionaryFs(dir string) (dictionaryFs, error) {
	var result *externalDictionaryFs
	if s, err := os.Stat(dir); err == nil {
		if !s.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dir)
		}
		result = &externalDictionaryFs{basePath: dir}
	} else {
		// try stripping off the last part
		base := filepath.Dir(dir)
		if s, err = os.Stat(base); err != nil || !s.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", dir)
		}
		result = &externalDictionaryFs{basePath: base, prefix: filepath.Base(dir)}
	}
	return result, result.checkFiles()
}

type externalDictionaryFs struct {
	basePath string
	prefix   string
}

func (e *externalDictionaryFs) Open(wordLen int) (io.ReadCloser, error) {
	return os.Open(filepath.Join(e.basePath, e.prefix+strconv.Itoa(wordLen)+"-letters.txt"))
}

func (e *externalDictionaryFs) checkFiles() error {
	if e.prefix == "" {
		rx := regexp.MustCompile(`^(.*-)\d+-letters\.txt$`)
		_ = filepath.WalkDir(e.basePath, func(path string, de os.DirEntry, err error) error {
			if !de.IsDir() {
				if m := rx.FindStringSubmatch(de.Name()); m != nil {
					e.prefix = m[1]
					return io.EOF
				}
			}
			return nil
		})
	}
	for wl := 2; wl <= 15; wl++ {
		if s, err := os.Stat(filepath.Join(e.basePath, e.prefix+strconv.Itoa(wl)+"-letters.txt")); err != nil || s.IsDir() {
			return fmt.Errorf("%s/%s%d-letters.txt does not exist or not a file", e.basePath, e.prefix, wl)
		}
	}
	return nil
}
