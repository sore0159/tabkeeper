package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const HTTP_PORT = ":8001"
const FILE_DIR_NAME = "FILES/"
const TAB_FILE_NAME = FILE_DIR_NAME + "TAB_FILE.json"

var LOG *Logger

func main() {
	var err error
	if LOG, err = NewLogger(); err != nil {
		log.Printf("ABORTING: Failed to initialize logging: %v\n", err)
		return
	}
	dn := make(chan byte)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	LOG.Inform("Starting server on port", HTTP_PORT)
	// Creating my own server var to have access to server.Shutdown()
	sf := NewSafeFiler(TAB_FILE_NAME)
	m := MakeMux(sf)
	server := &http.Server{Addr: HTTP_PORT, Handler: m}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			LOG.ServerErr("Listen and Serve Error: %v", err)
			dn <- 0
		}
	}()
	select {
	case <-ch:
		LOG.NewLine()
		LOG.Inform("Termination signal recieved, stopping server...")
		ctx := context.TODO()
		err := server.Shutdown(ctx)
		if err != nil {
			LOG.ServerErr("shutdown failure: %v", err)
		}
	case <-dn:
		LOG.NewLine()
		LOG.Inform("Exiting program...")
	}
}
