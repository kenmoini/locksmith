package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// DirectoryListingNames lists just the name of files in a certain directory
func DirectoryListingNames(path string) []string {
	if path == "" {
		path = "."
	}

	file, err := os.Open(path)
	check(err)
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	var fileNames []string
	for _, name := range list {
		fileNames = append(fileNames, name)
		//fmt.Println(name)
	}
	return fileNames
}

// FileExists checks if a file exists and returns a boolean or an erro
func FileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err == nil {
		// path/to/whatever exists
		return true, nil
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

// DirectoryExists checks if a file exists and returns a boolean or an erro
func DirectoryExists(pathName string) (bool, error) {
	if _, err := os.Stat(pathName); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false, nil
	}
	if _, err := os.Stat(pathName); !os.IsNotExist(err) {
		// path/to/whatever exists
		return true, nil
	}
	return false, nil
}

// TouchFile just creates a file if it doesn't exist already
func TouchFile(fileName string, updateTime bool) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		check(err)
		defer file.Close()
	} else {
		if updateTime {
			currentTime := time.Now().Local()
			err = os.Chtimes(fileName, currentTime, currentTime)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// CopyFile copies a file
func CopyFile(src, dst string, BUFFERSIZE int64) error {
	log.Printf("Copying  %s to %s\n", src, dst)
	if BUFFERSIZE == 0 {
		BUFFERSIZE = 4096
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	log.Printf("Downloading %s to %s\n", url, filepath)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, srcFile string) error {
	log.Printf("Extracting %s to %s\n", srcFile, dst)
	r, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

// CreateDirectory is self explanitory
func CreateDirectory(path string) {
	log.Printf("Creating directory %s\n", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		check(errDir)
	}
}

// DeleteFile deletes a file
func DeleteFile(path string) {
	log.Printf("Deleting %s\n", path)
	e := os.Remove(path)
	check(e)
}
