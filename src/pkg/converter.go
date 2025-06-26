package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Download(p *Package) error {
	const tempZip = "repo_temp.zip"
	resp, err := http.Get(fmt.Sprintf("https://github.com/ravendevteam/%s/archive/refs/tags/%s.zip", p.Name, p.Version))
	if err != nil {
		return fmt.Errorf("failed to download repo: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(tempZip)
	if err != nil {
		return fmt.Errorf("failed to create temp zip: %w", err)
	}
	defer os.Remove(tempZip)
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write zip: %w", err)
	}

	// Extract the ZIP
	r, err := zip.OpenReader(tempZip)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, strings.Join(strings.SplitN(f.Name, "/", 2)[1:], "/"))

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		in, err := f.Open()
		if err != nil {
			return err
		}
		defer in.Close()

		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, in); err != nil {
			return err
		}
	}

	return nil
}
