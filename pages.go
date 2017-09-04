package main

import (
	"html/template"
	"net/http"
)

const TEMPLATE_FILE_NAME = FILE_DIR_NAME + "template.html"

func MakeMux(sf *SafeFiler) *http.ServeMux {
	mux := http.NewServeMux()

	const STATIC_DIR = "FILES"
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, STATIC_DIR+"/img/yd32.ico")
	})
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(STATIC_DIR+"/img"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(STATIC_DIR+"/css"))))
	mux.Handle("/", sf)
	return mux
}

func (s *SafeFiler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	var assume int // todo
	pTab, err := ProcessTab(tab, assume)
	if err != nil {
		LOG.ServerErr("Failed to process tab: %s", err.Error())
		http.Error(w, "TAB PROCESS ERROR: "+err.Error(), 500)
		return
	}
	tp.ExecuteTemplate(w, "frame", pTab)
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
