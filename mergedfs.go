package pufferpanel

import (
	"io/fs"
	"os"
)

type MergedFS struct {
	priorityFs, defFs fs.FS
}

func NewMergedFS(a, b fs.FS) *MergedFS {
	return &MergedFS{
		priorityFs: a,
		defFs:      b,
	}
}

type Filesystem interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

func (m *MergedFS) Open(name string) (fs.File, error) {
	f, e := m.priorityFs.Open(name)
	if e != nil && os.IsNotExist(e) {
		return m.defFs.Open(name)
	}

	return f, e
}

func (m *MergedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	results := make([]fs.DirEntry, 0)
	var err error

	if dirFs, ok := m.priorityFs.(fs.ReadDirFS); ok {
		results, err = dirFs.ReadDir(name)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		err = nil
	}

	if dirFs, ok := m.defFs.(fs.ReadDirFS); ok {
		secondary, err := dirFs.ReadDir(name)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		for _, v := range secondary {
			inPrimary := false
			for _, z := range results {
				if v.Name() == z.Name() {
					inPrimary = true
				}
			}
			if !inPrimary {
				results = append(results, v)
			}
		}
	}

	return results, nil
}

func (m *MergedFS) ReadFile(name string) (data []byte, err error) {
	var primaryFile, secondaryFile fs.File

	if readFs, ok := m.priorityFs.(fs.ReadFileFS); ok {
		data, err = readFs.ReadFile(name)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	} else {
		//primary FS does not expose the read file endpoint. We must read it directly
		primaryFile, err = m.priorityFs.Open(name)
		defer Close(primaryFile)

		if err != nil && !os.IsNotExist(err) {
			return
		}
		if err == nil {
			//to get the buffer... we have to make it the right size
			var fi fs.FileInfo
			fi, err = primaryFile.Stat()
			if err != nil {
				return
			}
			data = make([]byte, fi.Size())
			_, err = primaryFile.Read(data)
			return
		}
	}

	//if we got here, primary could not fulfill request. We must do it through the default
	if readFs, ok := m.defFs.(fs.ReadFileFS); ok {
		data, err = readFs.ReadFile(name)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	} else {
		secondaryFile, err = m.defFs.Open(name)
		defer Close(secondaryFile)

		if err != nil && !os.IsNotExist(err) {
			return
		}
		if err == nil {
			var fi fs.FileInfo
			fi, err = secondaryFile.Stat()
			if err != nil {
				return
			}
			data = make([]byte, fi.Size())
			_, err = secondaryFile.Read(data)
			return
		}
	}

	return nil, os.ErrNotExist
}
