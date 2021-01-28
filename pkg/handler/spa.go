package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"net/http"
)

type SpaHandler struct {
	Dist  fs.FS
	Index string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s%s", "dist", r.URL.Path)

	_, err := h.Dist.Open(path)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		bytesWritten, err := fmt.Fprint(w, h.Index)
		if err != nil {
			logrus.WithError(err).WithField("Bytes written", bytesWritten).Error("Unable to send index.html")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	r.URL.Path = path
	http.FileServer(http.FS(h.Dist)).ServeHTTP(w, r)
}
