package main

import (
	"html/template"
	"net/http"
	"strings"
)

const TEMPLATE_FILE_NAME = FILE_DIR_NAME + "template.html"

func MakeMux(sf *SafeFiler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", sf)
	return mux
}

func (s *SafeFiler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if false { // TODO: get this straightened out
		rIP := r.Header.Get("x-forwarded-for")
		if !strings.HasPrefix(rIP, "192.168.1.") && !strings.HasPrefix(rIP, "127.0.0.1") && !strings.HasPrefix(rIP, "127.0.0.") {
			http.Error(w, "Does not support nonlocal connections", 400)
			return
		}
	}
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", 301)
		return
	}
	if r.Method == "POST" {
		s.HandlePost(w, r)
		return
	}

	tab, err := s.ReadTab()
	if err != nil {
		LOG.ServerErr("Failed to read to tab: %s", err.Error())
		http.Error(w, "TAB READ ERROR", 500)
		return
	}

	tp, err := template.ParseFiles(TEMPLATE_FILE_NAME)
	if err != nil {
		LOG.ServerErr("Failed to read to template: %s", err.Error())
		http.Error(w, "TEMPLATE READ ERROR", 500)
		return
	}
	tp.ExecuteTemplate(w, "frame", tab)
}

func (sf *SafeFiler) HandlePost(w http.ResponseWriter, r *http.Request) {
	entry, err := EntryFromPost(r)
	if err != nil {
		http.Error(w, "BAD USER DATA: "+err.Error(), 400)
		return
	}
	if err = sf.AppendToTab(entry); err != nil {
		LOG.ServerErr("Failed to write to tab: %s", err.Error())
		http.Error(w, "TAB WRITE ERROR", 500)
		return
	}
	http.Redirect(w, r, "/", 301)
}
