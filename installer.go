package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func downloadAndUntar(url, dest, folderName string, wg *sync.WaitGroup) error {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dest, folderName, header.Name[8:])
		info := header.FileInfo()
		// // fmt.Println(target)

		if info.IsDir() {
			err = os.MkdirAll(target, info.Mode())
			if err != nil {
				return err
			}
			continue
		}
		dirPath := filepath.Dir(target)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			os.MkdirAll(dirPath, info.Mode())
		}
		file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
}

// func untarGz(src, dest string) error {
// 	file, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	gzipReader, err := gzip.NewReader(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer gzipReader.Close()

// 	tarReader := tar.NewReader(gzipReader)

// 	for {
// 		header, err := tarReader.Next()

// 		switch {
// 		case err == io.EOF:
// 			return nil
// 		case err != nil:
// 			return err
// 		case header == nil:
// 			continue
// 		}

// 		target := filepath.Join(dest, header.Name)
// 		info := header.FileInfo()
// 		// fmt.Println(info)

// 		if info.IsDir() {
// 			err = os.MkdirAll(target, info.Mode())
// 			if err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		dirPath := filepath.Dir(target)
// 		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
// 			err = os.MkdirAll(dirPath, info.Mode())
// 		}
// 		file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, info.Mode())
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		_, err = io.Copy(file, tarReader)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// func installPackage(downloadUrl, dest string) error {
// 	resp, err := http.Get(downloadUrl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// fmt.Println(resp.Body)
// 	defer resp.Body.Close()

// 	out, err := os.Create("archive.tar.gz")
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()
// 	_, err = io.Copy(out, resp.Body)
// 	// fmt.Println(untarGz("archive.tar.gz", MODULES))
// 	return err
// }
