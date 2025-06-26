package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/frijoles524/Blacklight/pkg"
)

func downloadRuntime() error {
	const (
		repoZipURL     = "https://codeload.github.com/frijoles524/Blacklight/zip/refs/heads/master"
		tempZipPath    = "temp_repo.zip"
		tempExtractDir = "temp_repo"
		runtimeSubDir  = "src/pkg/python312runtime"
		destDir        = "python312runtime"
	)
	if _, err := os.Stat(destDir); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if runtime exists: %w", err)
	}
	resp, err := http.Get(repoZipURL)
	if err != nil {
		return fmt.Errorf("failed to download repo: %w", err)
	}
	defer resp.Body.Close()
	out, err := os.Create(tempZipPath)
	if err != nil {
		return fmt.Errorf("failed to create temp zip: %w", err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	if err := unzip(tempZipPath, tempExtractDir); err != nil {
		return fmt.Errorf("failed to unzip repo: %w", err)
	}
	srcRuntimePath := filepath.Join(tempExtractDir, runtimeSubDir)
	if err := os.RemoveAll(destDir); err != nil {
		return fmt.Errorf("failed to clear existing runtime: %w", err)
	}
	if err := os.Rename(srcRuntimePath, destDir); err != nil {
		return fmt.Errorf("failed to move runtime folder: %w", err)
	}
	os.RemoveAll(tempZipPath)
	os.RemoveAll(tempExtractDir)
	return nil
}
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ParseFlags()
	pkg.FetchPackageList()
	if err := downloadRuntime(); err != nil {
		panic(err)
	}
	switch *command {
	case "run":
		RunApp()
	}
}
