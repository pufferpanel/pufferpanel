package pufferpanel

import (
	"io"
	"net/http"
)

func Close(closer io.Closer) {
	//at this point, i give up trying to get this to not fail, so we'll go with the brute force
	defer func() {
		recover()
	}()

	if closer != nil {
		_ = closer.Close()
	}
}

func CloseResponse(response *http.Response) {
	if response != nil {
		Close(response.Body)
	}
}
