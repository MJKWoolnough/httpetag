// Package httpetag wraps a http.FileStstem to provide ETag headers
package httpetag // import "vimagination.zapto.org/httpetag"

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"io"
	"net/http"
	"path"
	"sync"
	"time"
)

type file struct {
	LastModified time.Time
	Etag         string
}

type fileServer struct {
	root    http.FileSystem
	handler http.Handler

	mu    sync.RWMutex
	etags map[string]file
}

var hashPool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

// New creates a new Handler around an http.FileSystem which generates sha256
// hash's of the file contents as ETags
func New(root http.FileSystem) http.Handler {
	return &fileServer{
		root:    root,
		handler: http.FileServer(root),
	}
}

// NewWithHandler creates a new handler, as in New, but with a custom handler
// which is passed to after setting the ETag header
func NewWithHandler(root http.FileSystem, handler http.Handler) http.Handler {
	return &fileServer{
		root:    root,
		handler: handler,
	}
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Clean(r.URL.Path)
	f, err := fs.root.Open(p)
	if err == nil {
		s, err := f.Stat()
		if err == nil {
			fs.mu.RLock()
			fe, ok := fs.etags[p]
			fs.mu.RUnlock()
			mt := s.ModTime()
			if !ok || !fe.LastModified.Equal(mt) {
				h := hashPool.Get().(hash.Hash)
				_, err = io.Copy(h, f)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				fe.LastModified = mt
				fe.Etag = base64.StdEncoding.EncodeToString(h.Sum(nil))
				h.Reset()
				hashPool.Put(h)
				fs.mu.Lock()
				fs.etags[p] = fe
				fs.mu.Unlock()
			}
			w.Header().Set("ETag", fe.Etag)
		}
	}
	fs.handler.ServeHTTP(w, r)
}
