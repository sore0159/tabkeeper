package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type SafeFiler struct {
	FileName  string
	ProxyAddr string
	*sync.RWMutex
}

func NewSafeFiler(fileName, proxy string) *SafeFiler {
	return &SafeFiler{
		FileName:  fileName,
		ProxyAddr: proxy,
		RWMutex:   new(sync.RWMutex),
	}
}

func (sf *SafeFiler) ReadTab() ([]*Entry, error) {
	sf.RLock()
	defer sf.RUnlock()
	f, err := os.OpenFile(sf.FileName, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Entry{}, nil
		}
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var tab []*Entry
	if err = dec.Decode(&tab); err != nil {
		return nil, fmt.Errorf("failed to decode file data: %v", err)
	}
	return tab, nil
}

func (sf *SafeFiler) AppendToTab(entry *Entry) error {
	sf.Lock()
	defer sf.Unlock()
	var tab []*Entry
	f, err := os.OpenFile(sf.FileName, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		tab = make([]*Entry, 0, 1)
	} else if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	} else {
		dec := json.NewDecoder(f)
		if err = dec.Decode(&tab); err != nil {
			f.Close()
			return fmt.Errorf("failed to decode tab: %v", err)
		}
		f.Close()
	}

	tab = append(tab, entry)

	f, err = os.OpenFile(sf.FileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %v", err)
	}
	defer f.Close()
	var out bytes.Buffer
	data, err := json.Marshal(&tab)
	if err != nil {
		return fmt.Errorf("failed to encode tab: %v", err)
	}
	json.Indent(&out, data, "", "  ")
	_, err = out.WriteTo(f)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}
