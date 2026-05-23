package files

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v3"
	"github.com/pufferpanel/pufferpanel/v3/utils"
)

const PathSeparator = "/"

func Compress(targetFS FileServer, targetFile string, sourceFS FileServer, files []string) error {
	if len(files) == 0 {
		return errors.New("no files to compress")
	}

	c, err := archiver.ByExtension(targetFile)
	if err != nil {
		return err
	}
	var compressor archiver.Writer
	var ok bool
	if compressor, ok = c.(archiver.Writer); !ok {
		return archiver.ErrFormatNotRecognized
	}

	for _, file := range files {
		err = writeFileToArchive(compressor, sourceFS, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func Extract(sourceFS FileServer, sourceFile string, targetFS FileServer, targetPath string, filter string) error {
	a, err := archiver.ByExtension(sourceFile)
	if err != nil {
		return err
	}

	var extractor archiver.Reader
	var ok bool
	if extractor, ok = a.(archiver.Reader); !ok {
		return archiver.ErrFormatNotRecognized
	}

	source, err := sourceFS.OpenFile(sourceFile, os.O_RDONLY, 0644)
	defer utils.Close(source)
	defer utils.Close(extractor)
	if err != nil {
		return err
	}

	err = extractor.Open(source, 0)
	if err != nil {
		return err
	}

	var file archiver.File
	for {
		file, err = extractor.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = walk(file, targetFS, targetPath, filter)
		if err != nil {
			return err
		}
	}
	return nil
}

func walk(file archiver.File, fs FileServer, targetPath, filter string) (err error) {
	path := getCompressedItemName(file)

	if !utils.CompareWildcard(file.Name(), filter) {
		return
	}

	parent := filepath.Join(targetPath, filepath.Dir(path))
	path = filepath.Join(targetPath, path)

	if file.Mode().IsDir() {
		if err = fs.MkdirAll(path, 0755); err != nil {
			return err
		}
	} else if file.Mode().IsRegular() {
		if err = fs.MkdirAll(parent, 0755); err != nil {
			return err
		}

		var outFile *os.File
		outFile, err = fs.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Mode()|0600)
		if err != nil {
			return err
		}
		defer utils.Close(outFile)
		_, err = io.Copy(outFile, file.ReadCloser)
	} else if file.Mode()&os.ModeSymlink != 0 {
		target, err := getLinkTarget(file)
		if err != nil {
			return err
		}

		if err = fs.MkdirAll(parent, 0755); err != nil {
			return err
		}
		if err = fs.Symlink(target, path); err != nil {
			return err
		}
	}
	return
}

// getCompressedItemName Resolves headers in the event the wrapped interface fails
func getCompressedItemName(file archiver.File) string {
	//For certain headers, the actual File interface uses the wrong value
	//Example, ZIP gives the filename, not the full path

	switch v := file.Header.(type) {
	case zip.FileHeader:
		return v.Name
	case *tar.Header:
		return v.Name
	default:
		return file.Name()
	}
}

func getLinkTarget(file archiver.File) (string, error) {
	switch v := file.Header.(type) {
	case *tar.Header:
		return v.Linkname, nil
	case zip.FileHeader:
		buffer := make([]byte, file.Size())
		size, err := file.Read(buffer)
		if err != nil {
			return "", err
		}
		return string(buffer[:size]), nil
	default:
		return "", archiver.ErrFormatNotRecognized
	}
}

func writeFileToArchive(writer archiver.Writer, fs FileServer, path string) error {
	fi, err := fs.Stat(path)
	if err != nil {
		return err
	}

	if fi.Mode().IsRegular() {
		file, err := fs.Open(path)
		defer utils.Close(file)
		if err != nil {
			return err
		}
		err = writer.Write(archiver.File{
			ReadCloser: file,
		})
		return err
	} else if fi.Mode().IsDir() {
		files, err := fs.ReadDir(path)
		if err != nil {
			return err
		}
		for _, v := range files {
			err = writeFileToArchive(writer, fs, filepath.Join(path, v.Name()))
			if err != nil {
				return err
			}
		}
	} //TODO: Add Symlink support...?
	return nil
}
