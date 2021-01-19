package handler

import (
	"embed"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
)

type SpaHandler struct {
	Dist  embed.FS
	Index string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = fmt.Sprintf("%s%s", "dist", path)
	r.URL.Path = path

	_, err = h.Dist.Open(path)
	if os.IsNotExist(err) {
		w.Header().Set("Content-Type", "text/html")
		bytesWritten, err := fmt.Fprint(w, h.Index)
		if err != nil {
			logrus.WithError(err).WithField("Bytes written", bytesWritten).Error("Unable to send index.html")
		}
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.FS(h.Dist)).ServeHTTP(w, r)
}
