package pkg

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Package struct {
	Name     string
	Version  string
	Checksum string
}

const cacheFile = "installed_cache.json"

var cache = map[string]Package{}

func LoadCache() {
	file, err := os.Open(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error loading cache:", err)
		return
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&cache)
}

func saveCache() {
	file, err := os.Create(cacheFile)
	if err != nil {
		fmt.Println("Error saving cache:", err)
		return
	}
	defer file.Close()
	_ = json.NewEncoder(file).Encode(cache)
}

func IsInstalled(p *Package) bool {
	if cached, ok := cache[p.Name]; ok {
		if p.Version == "" {
			return true
		}
		return cached.Version == p.Version
	}
	return false
}

func Download(p *Package) error {
	if IsInstalled(p) {
		fmt.Printf("Package %s@%s is already installed.\n", p.Name, p.Version)
		return nil
	}

	const tempZip = "repo_temp.zip"
	url := fmt.Sprintf("https://github.com/ravendevteam/%s/archive/refs/tags/%s.zip", p.Name, p.Version)
	fmt.Println("Downloading:", url)
	resp, err := http.Get(url)
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

	r, err := zip.OpenReader(tempZip)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := f.Name
		if fpath == "" {
			continue
		}

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

		outFile, err := os.Create(fpath)
		if err != nil {
			in.Close()
			return err
		}

		if _, err := io.Copy(outFile, in); err != nil {
			outFile.Close()
			in.Close()
			return err
		}

		outFile.Close()
		in.Close()
	}
	CleanDirectory(fmt.Sprintf("%s-%s", p.Name, p.Version), []string{fmt.Sprintf("%s.py", p.Name), fmt.Sprintf("%s.png", p.Name), "style.css", "requirements.txt", "ICON.ico"})
	cache[p.Name] = *p
	saveCache()
	fmt.Printf("Package %s@%s installed.\n", p.Name, p.Version)
	return nil
}

func CleanDirectory(dir string, keepFiles []string) error {
	keepMap := make(map[string]bool)
	for _, f := range keepFiles {
		keepMap[f] = true
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == dir {
			return nil
		}

		if keepMap[info.Name()] {
			return nil
		}

		if info.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", path, err)
			}
			return filepath.SkipDir
		} else {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", path, err)
			}
		}
		return nil
	})
}

func GetInstalledVersions(name string) []string {
	var versions []string
	for _, p := range cache {
		if p.Name == name {
			versions = append(versions, p.Version)
		}
	}
	return versions
}

func parseSemVer(v string) ([]int, error) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	intParts := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		intParts[i] = n
	}
	return intParts, nil
}

func semVerLess(a, b string) bool {
	aParts, errA := parseSemVer(a)
	bParts, errB := parseSemVer(b)
	if errA != nil || errB != nil {
		return a < b
	}
	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		if aParts[i] < bParts[i] {
			return true
		} else if aParts[i] > bParts[i] {
			return false
		}
	}
	return len(aParts) < len(bParts)
}

func GetHighestVersion(name string) (string, error) {
	versions := GetInstalledVersions(name)
	if len(versions) == 0 {
		return "", fmt.Errorf("no installed versions for %s", name)
	}

	sort.Slice(versions, func(i, j int) bool {
		return semVerLess(versions[i], versions[j])
	})

	return versions[len(versions)-1], nil
}
