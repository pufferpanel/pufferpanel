package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type ImageDownload struct {
	Status         string
	ProgressDetail ProgressDetail
	Id             string
	Progress       string
}

type ProgressDetail struct {
	Current int64
	Total   int64
}

type ImageWriter struct {
	io.Writer
	Parent io.Writer
}

func (w *ImageWriter) Write(data []byte) (n int, err error) {
	//what happens is the writer should send us a json blob, which has the status we want
	//to handle this, we will parse it, and then format it to our upper-writer
	//NOTE: WE CAN GET MULTIPLE LINES IN 1 BLOCK

	var buf bytes.Buffer
	n, err = buf.Write(data)
	if err != nil {
		return
	}

	for err != io.EOF {
		d, e := buf.ReadBytes(byte('\n'))
		if e != nil && e != io.EOF {
			err = e
			return
		}

		if d == nil || len(d) == 0 {
			break
		}

		var imageDownload ImageDownload
		err = json.Unmarshal(d, &imageDownload)
		if err != nil {
			return
		}

		message := fmt.Sprintf("%s %s %s", imageDownload.Status, imageDownload.Id, strings.ReplaceAll(imageDownload.Progress, "\u003e", ""))
		message = strings.TrimSpace(message)
		_, err = w.Parent.Write([]byte(message + "\n"))
	}

	if err == io.EOF {
		err = nil
	}

	return
}
