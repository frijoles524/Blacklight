package pkg

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	remoteListURL = "https://raw.githubusercontent.com/frijoles524/Blacklight/master/package_information/blacklight.list"
	cacheFilePath = "blacklight.list"
)

var (
	Packages []Package
)

type Package struct {
	Name     string
	Version  string
	URL      string
	Checksum string
}

func FetchPackageList() {
	res, err := http.Get(remoteListURL)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("failed to fetch remote list, using cache...")
		Packages, err = parseListFile(cacheFilePath)
		if err != nil {
			panic("failed to parse package index")
		}
	}
	defer res.Body.Close()

	out, err := os.Create(cacheFilePath)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(out, res.Body)
	Packages, err = parseListFile(cacheFilePath)
	if err != nil {
		panic(err)
	}
}

func parseListFile(path string) ([]Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open package list: %w", err)
	}
	defer file.Close()

	var packages []Package
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		packages = append(packages, Package{
			Name:     parts[0],
			Version:  parts[1],
			URL:      parts[2],
			Checksum: parts[3],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return packages, nil
}
