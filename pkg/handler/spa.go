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

func NewSpaHandler(dist fs.FS, index string) *SpaHandler {
	return &SpaHandler{Dist: dist, Index: index}
}

//ServeHTTP handles incoming http requests. Always responds with the index.html if the request resource does not exist
func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s%s", "dist", r.URL.Path)

	_, err := h.Dist.Open(path)
	//Respond with index.html
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		bytesWritten, err := fmt.Fprint(w, h.Index)
		if err != nil {
			logrus.WithError(err).WithField("Bytes written", bytesWritten).Error("Unable to send index.html")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	//serve the requested file from the dist directory
	r.URL.Path = path
	http.FileServer(http.FS(h.Dist)).ServeHTTP(w, r)
}
