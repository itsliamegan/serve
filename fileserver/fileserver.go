package fileserver

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type FileServer struct {
	rootDir string
}

func New(rootDir string) *FileServer {
	return &FileServer{rootDir}
}

func (server *FileServer) Listen(addr string) error {
	return http.ListenAndServe(addr, server)
}

func (server *FileServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := filepath.Join(server.rootDir, req.URL.Path)

	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			indexFile := filepath.Join(path, "index.html")
			_, err = os.Stat(indexFile)

			if err == nil {
				file, err := os.Open(indexFile)
				if err != nil {
					res.Header().Set("Content-Type", "text/plain; charset=utf-8")
					res.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(res, "500 internal server error")
					return
				}

				res.Header().Set("Content-Type", "text/html; charset=utf-8")
				res.WriteHeader(http.StatusOK)
				io.Copy(res, file)
			} else if errors.Is(err, os.ErrNotExist) {
				name := filepath.Base(path)
				entries, err := os.ReadDir(path)
				if err != nil {
					res.Header().Set("Content-Type", "text/plain; charset=utf-8")
					res.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintln(res, "500 internal server error")
					return
				}

				res.Header().Set("Content-Type", "text/html; charset=utf-8")
				res.WriteHeader(http.StatusOK)

				fmt.Fprint(res, "<!doctype html><html><head>")
				fmt.Fprint(res, "<meta charset=\"utf-8\">")
				fmt.Fprint(res, "<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\">")
				fmt.Fprint(res, "<style>")
				fmt.Fprint(res, "html { font-family: monospace; }")
				fmt.Fprint(res, "h1 { font-size: 1rem; font-weight: normal; }")
				fmt.Fprint(res, "ul { list-style: none; padding: 0; }")
				fmt.Fprint(res, "</style></head><body>")
				fmt.Fprintf(res, "<h1>%s/</h1>", name)
				fmt.Fprint(res, "<ul>")

				for _, entry := range entries {
					name := entry.Name()
					if entry.IsDir() {
						name += "/"
					}

					fmt.Fprint(res, "<li>")
					fmt.Fprintf(res, "<a href=\"%s\">%s</a>", name, name)
					fmt.Fprint(res, "</li>")
				}

				fmt.Fprint(res, "</ul></body>")
			} else {
				res.Header().Set("Content-Type", "text/plain; charset=utf-8")
				res.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(res, "500 internal server error")
			}
		} else {
			file, err := os.Open(path)
			if err != nil {
				res.Header().Set("Content-Type", "text/plain; charset=utf-8")
				res.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(res, "500 internal server error")
				return
			}

			extension := filepath.Ext(path)
			mimeType := mime.TypeByExtension(extension)
			if mimeType == "" {
				mimeType = "text/plain; charset=utf-8"
			}

			res.Header().Set("Content-Type", mimeType)
			res.WriteHeader(http.StatusOK)
			io.Copy(res, file)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		htmlFile := path + ".html"
		_, err = os.Stat(htmlFile)

		if err == nil {
			file, err := os.Open(htmlFile)
			if err != nil {
				res.Header().Set("Content-Type", "text/plain; charset=utf-8")
				res.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(res, "500 internal server error")
				return
			}

			res.Header().Set("Content-Type", "text/html; charset=utf-8")
			res.WriteHeader(http.StatusOK)
			io.Copy(res, file)
		} else if errors.Is(err, os.ErrNotExist) {
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(res, "404 not found")
		} else {
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(res, "500 internal server error")
		}
	} else {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(res, "500 internal server error")
	}
}
