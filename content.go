package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (s *SafeFiler) WritePage(w http.ResponseWriter, r *http.Request) {
	str := strings.TrimPrefix(r.URL.Path, "/write/")
	if str == "" {
		http.Error(w, "NO WRITE DATA FOUND", 400)
		return
	}
	str = fmt.Sprintf("%s: %s\n", time.Now(), str)
	f, err := s.GetW()
	if err == nil {
		defer f.Close()
		_, err = fmt.Fprintf(f, str)
	}
	if err != nil {
		LogServerErr("Failed to write to file: %s", err.Error())
		http.Error(w, "DATA WRITE ERROR", 500)
		return
	}
	fmt.Fprintf(w, "WROTE: %s", str)
}
func (s *SafeFiler) ReadPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", 301)
		return
	}
	f, err := s.GetR()
	if err == nil {
		defer f.Close()
		_, err = io.Copy(w, f)
	}
	if err != nil {
		LogServerErr("Failed to read to file: %s", err.Error())
		http.Error(w, "DATA READ ERROR", 500)
		return
	}
}
