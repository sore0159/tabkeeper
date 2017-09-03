package main

import (
	"net/http"
	"os"
	"sync"
)

func MakeMux(f *SafeFiler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/write/", f.WritePage)
	mux.HandleFunc("/", f.ReadPage)
	return mux
}

type SafeFiler struct {
	FileName string
	*sync.RWMutex
}

type SafeFile struct {
	*os.File
	s *SafeFiler
	r bool
}

func NewSafeFiler(fileName string) *SafeFiler {
	return &SafeFiler{
		FileName: fileName,
		RWMutex:  new(sync.RWMutex),
	}
}

func (s *SafeFiler) GetR() (*SafeFile, error) {
	s.RLock()
	f, err := os.OpenFile(s.FileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &SafeFile{
		File: f,
		s:    s,
		r:    true,
	}, nil
}

func (s *SafeFiler) GetW() (*SafeFile, error) {
	s.Lock()
	f, err := os.OpenFile(s.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &SafeFile{
		File: f,
		s:    s,
		r:    false,
	}, nil
}

// This doesn't deal with an error during closing
// 'poisoned' data
func (sf *SafeFile) Close() error {
	if sf.r {
		sf.s.RUnlock()
	} else {
		sf.s.Unlock()
	}
	return sf.File.Close()
}
