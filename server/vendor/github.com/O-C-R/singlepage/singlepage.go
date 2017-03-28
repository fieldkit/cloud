package singlepage

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func init() {
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".map", "application/json")
}

var (
	directoryError = errors.New("directory error")
)

type FileInfo struct {
	Path          string
	Etag          string
	IsDir         bool
	LongtermCache bool
}

type SinglePageApplicationOptions struct {
	Root                       string
	Application, LongtermCache *regexp.Regexp
}

type SinglePageApplication struct {
	options  SinglePageApplicationOptions
	fileInfo map[string]*FileInfo
}

func NewSinglePageApplication(options SinglePageApplicationOptions) (*SinglePageApplication, error) {
	s := &SinglePageApplication{
		options:  options,
		fileInfo: make(map[string]*FileInfo),
	}

	root := filepath.Clean(options.Root)
	if err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		rootPath := strings.TrimPrefix(path, root)
		s.fileInfo[rootPath] = &FileInfo{
			Path:          path,
			IsDir:         fileInfo.IsDir(),
			LongtermCache: s.options.LongtermCache != nil && s.options.LongtermCache.MatchString(rootPath),
		}

		if fileInfo.IsDir() {
			return nil
		}

		// Attempt to open the file.
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		// Defer closing the file.
		defer file.Close()

		// Get the SHA1 hash value for the file.
		fileHash := sha1.New()
		if _, err := io.Copy(fileHash, file); err != nil {
			return err
		}

		// Set an ETag header based on the SHA1 hash.
		s.fileInfo[rootPath].Etag = `"` + base64.StdEncoding.EncodeToString(fileHash.Sum(nil)) + `"`

		return nil
	}); err != nil {
		return nil, err
	}

	return s, nil
}

func acceptEncodingGzip(req *http.Request) bool {
	encodings := strings.Split(req.Header.Get("accept-encoding"), ",")
	for _, encoding := range encodings {
		if strings.TrimSpace(encoding) == "gzip" {
			return true
		}
	}

	return false
}

func (s *SinglePageApplication) openFileInfo(req *http.Request, name string) (*FileInfo, bool, error) {
	fileInfo, ok := s.fileInfo[name]
	if !ok {
		return nil, false, os.ErrNotExist
	}

	if fileInfo.IsDir {
		return nil, false, directoryError
	}

	if acceptEncodingGzip(req) {
		fileInfo, ok := s.fileInfo[name+".gz"]
		if ok {
			return fileInfo, true, nil
		}
	}

	return fileInfo, false, nil
}

func (s *SinglePageApplication) serveFile(w http.ResponseWriter, req *http.Request, name string) error {
	fileInfo, gzipEncoded, err := s.openFileInfo(req, name)
	if err != nil {
		return err
	}

	// Attempt to open the file.
	file, err := os.Open(fileInfo.Path)
	if err != nil {
		return err
	}

	// Defer closing the file.
	defer file.Close()

	// Set an ETag header.
	w.Header().Set("etag", fileInfo.Etag)

	// If the file is gzip-encoded, set a Content-Encoding header.
	if gzipEncoded {
		w.Header().Set("content-encoding", "gzip")
		w.Header().Set("vary", "accept-encoding")
	}

	// Set a Cache-Control header.
	if fileInfo.LongtermCache {
		w.Header().Set("cache-control", "max-age=31536000")
	} else {
		w.Header().Set("cache-control", "max-age=60")
	}

	http.ServeContent(w, req, name, time.Time{}, file)
	return nil
}

func (s *SinglePageApplication) serveAny(w http.ResponseWriter, req *http.Request, name string) error {

	// Attempt to serve the file.
	err := s.serveFile(w, req, name)

	// If the file is a directory, attempt to serve an index.
	if err == directoryError {
		return s.serveFile(w, req, name+"/index.html")
	}

	// If the file doesn't exist, but the URL path matches the application
	// regexp, attempt to serve the application index.
	if os.IsNotExist(err) && (s.options.Application == nil || s.options.Application.MatchString(req.URL.Path)) {
		return s.serveFile(w, req, "/index.html")
	}

	// Return any other error.
	return err
}

func (s *SinglePageApplication) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !(req.Method == http.MethodGet || req.Method == http.MethodHead) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := req.URL.Path

	// Handle a request for index.html on the root level.
	if name == "/index.html" {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	// Handle trailing slashes.
	if name != "/" && strings.HasSuffix(name, "/") {
		http.Redirect(w, req, strings.TrimSuffix(req.URL.Path, "/"), http.StatusFound)
		return
	}

	// Handle a request for a sub-directory index.
	if strings.HasSuffix(name, "/index.html") {
		http.Redirect(w, req, strings.TrimSuffix(req.URL.Path, "/index.html"), http.StatusFound)
		return
	}

	var err error
	if name == "/" {
		err = s.serveFile(w, req, "/index.html")
	} else {
		err = s.serveAny(w, req, name)
	}

	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
