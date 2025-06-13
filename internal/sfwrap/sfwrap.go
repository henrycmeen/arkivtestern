package sfwrap

import (
	"os"
	"path/filepath"
	"github.com/richardlehane/siegfried"
)

// IdentifyFile returns Siegfried's identification result for a given file path.
func IdentifyFile(path string) (interface{}, error) {
	tmpfile, err := os.CreateTemp("", "arkivtestern-*.sig")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())
	_, err = tmpfile.Write(sig)
	if err != nil {
		tmpfile.Close()
		return nil, err
	}
	err = tmpfile.Close()
	if err != nil {
		return nil, err
	}
	sf, err := siegfried.Load(tmpfile.Name())
	if err != nil {
		return nil, err
	}
	return identifyFileOrDir(sf, path)
}

func safeIdentify(sf *siegfried.Siegfried, f *os.File, fp string) (interface{}, error) {
	defer func() {
		recover()
	}()
	match, err := sf.Identify(f, fp, "")
	if err != nil {
		return nil, err
	}
	return match, nil
}

func identifyFileOrDir(sf *siegfried.Siegfried, path string) (interface{}, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		var results []map[string]interface{}
		err := filepath.Walk(path, func(fp string, fi os.FileInfo, err error) error {
			if err != nil {
				results = append(results, map[string]interface{}{
					"file": fp,
					"error": err.Error(),
				})
				return nil // fortsett videre
			}
			if fi.IsDir() || !fi.Mode().IsRegular() {
				return nil
			}
			if fi.Size() == 0 {
				results = append(results, map[string]interface{}{
					"file": fp,
					"error": "empty file (not scanned)",
				})
				return nil
			}
			f, ferr := os.Open(fp)
			if ferr != nil {
				results = append(results, map[string]interface{}{
					"file": fp,
					"error": ferr.Error(),
				})
				return nil
			}
			buf := make([]byte, 1)
			n, rerr := f.Read(buf)
			if rerr != nil || n == 0 {
				f.Close()
				results = append(results, map[string]interface{}{
					"file": fp,
					"error": "unreadable or empty file (not scanned)",
				})
				return nil
			}
			f.Seek(0, 0) // reset file pointer
			match, errOrPanic := safeIdentify(sf, f, fp)
			f.Close()
			if errOrPanic != nil {
				results = append(results, map[string]interface{}{
					"file": fp,
					"error": errOrPanic.Error(),
				})
				return nil
			}
			results = append(results, map[string]interface{}{
				"file": fp,
				"result": match,
			})
			return nil

		})
		if err != nil {
			return nil, err
		}
		return results, nil
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		match, err := sf.Identify(f, path, "")
		if err != nil {
			return nil, err
		}
		return match, nil
	}
}
