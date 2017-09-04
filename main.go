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
const DEFAULT_TAB_FILE_NAME = FILE_DIR_NAME + "TAB_FILE.json"
const DEFAULT_PROXY_ADDR = "/tab"

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

	LOG.Record("Starting server on port %s", HTTP_PORT)
	// Creating my own server var to have access to server.Shutdown()
	var tabFile, proxyAddr string
	var flagT, flagP bool
	for _, str := range os.Args[1:] {
		if flagT {
			flagT = false
			tabFile = str
			LOG.Inform("Using %s for tabfile\n", str)
		} else if flagP {
			flagP = false
			proxyAddr = str
			LOG.Inform("Using %s for proxy address\n", str)
		} else if str == "-p" {
			flagP = true
		} else if str == "-t" {
			flagT = true
		}
	}
	if proxyAddr == "" {
		proxyAddr = DEFAULT_PROXY_ADDR
	}
	if tabFile == "" {
		tabFile = DEFAULT_TAB_FILE_NAME
	}

	sf := NewSafeFiler(tabFile, proxyAddr)
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
		LOG.Record("Termination signal recieved, stopping server...")
		ctx := context.TODO()
		err := server.Shutdown(ctx)
		if err != nil {
			LOG.ServerErr("shutdown failure: %v", err)
		}
	case <-dn:
		LOG.NewLine()
		LOG.Record("Exiting program...")
	}
}
