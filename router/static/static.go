package static

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	// DefaultMaxAge is 60 days.
	DefaultMaxAge = 86400
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
	MaxAge() int
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
	maxAge  int
}

func LocalFile(root string, maxAge int, indexes bool) *localFileSystem {
	if maxAge == 0 {
		maxAge = DefaultMaxAge
	}
	return &localFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
		maxAge:     maxAge,
	}
}

func (l *localFileSystem) MaxAge() int {
	return l.maxAge
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if !l.indexes && stats.IsDir() {
			return false
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, 0, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", fs.MaxAge()))
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
