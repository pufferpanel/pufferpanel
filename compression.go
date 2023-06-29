/*
 Copyright 2022 PufferPanel
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
 	http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package pufferpanel

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pault.ag/go/debian/deb"
	"strings"
)

func ExtractTar(stream io.Reader, directory string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	var tarReader *tar.Reader
	if r, isGood := stream.(*tar.Reader); isGood {
		tarReader = r
	} else {
		tarReader = tar.NewReader(stream)
	}

	var header *tar.Header
	for {
		header, err = tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(filepath.Join(directory, header.Name), 0755); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err = os.MkdirAll(filepath.Join(directory, filepath.Dir(header.Name)), 0755); err != nil {
				return err
			}

			//symlinks suck... so much
			sourceFile := filepath.Join(directory, header.Name)
			targetFile := header.Linkname
			if strings.HasPrefix(header.Linkname, "/") {
				targetFile = filepath.Join(directory, strings.TrimPrefix(header.Linkname, "/"))
			}

			if err = os.Symlink(targetFile, sourceFile); err != nil {
				return err
			}
		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Join(directory, filepath.Dir(header.Name)), 0755); err != nil {
				return err
			}
			var outFile *os.File
			outFile, err = os.Create(filepath.Join(directory, header.Name))
			if err != nil {
				return err
			}
			if _, err = io.Copy(outFile, tarReader); err != nil {
				_ = outFile.Close()
				return err
			}
			_ = outFile.Close()
			err = os.Chmod(filepath.Join(directory, header.Name), header.FileInfo().Mode())
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("uknown type: %s in %s", string([]byte{header.Typeflag}), header.Name)
		}
	}

	return nil
}

func ExtractTarGz(gzipStream io.Reader, directory string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()
	return ExtractTar(uncompressedStream, directory)
}

func ExtractDeb(stream io.ReaderAt, directory string) error {
	file, err := deb.Load(stream, directory)
	if err != nil {
		return err
	}

	defer file.Close()

	return ExtractTar(file.Data, directory)
}

func ExtractZip(name, directory string) error {
	file, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, f := range file.File {
		err = unzipFile(f, directory, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractZipIgnoreSingleDir(name, directory string) error {
	file, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileList []string
	for _, f := range file.File {
		fileList = append(fileList, f.Name)
	}

	dirs := make(map[string]bool)
	for _, f := range fileList {
		folderName := filepath.Dir(f)
		if folderName == ".." || folderName == "." || folderName == "/" {
			folderName = "."
		}
		dirs[folderName] = true
	}

	isSingleDir := true
	if len(dirs) > 0 {
		var rootDir string
		for k := range dirs {
			if rootDir == "" {
				firstPath := strings.SplitN(k, string(os.PathSeparator), 2)
				rootDir = firstPath[0]
			} else if rootDir != k && !strings.HasPrefix(k, rootDir+string(os.PathSeparator)) {
				isSingleDir = false
				break
			}
		}
	}

	for _, f := range file.File {
		err = unzipFile(f, directory, isSingleDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(f *zip.File, destination string, skipLevel bool) error {
	// 4. Check if file paths are not vulnerable to Zip Slip

	fileName := f.Name
	if skipLevel {
		parts := strings.SplitN(fileName, "/", 2)
		if len(parts) != 2 {
			return nil
		}
		fileName = parts[1]
		if fileName == "" {
			return nil
		}
	}

	filePath := filepath.Join(destination, fileName)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
