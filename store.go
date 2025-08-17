package main

import (
    "encoding/json"
    "log"
    "os"
    "sync"
)

type FileMeta struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Size int64  `json:"size"`
    Type string `json:"type"`
    URL  string `json:"url"`
    Path string `json:"-"`
}

type Store struct {
    sync.Mutex
    Files map[string]FileMeta `json:"files"`
}

var (
    dataDir   = "data"
    indexFile = "data/index.json"
    store     = &Store{Files: make(map[string]FileMeta)}
)

func saveIndex() {
    tmp := indexFile + ".tmp"
    f, err := os.Create(tmp)
    if err != nil {
        log.Println("failed to save index:", err)
        return
    }
    json.NewEncoder(f).Encode(store)
    f.Close()
    os.Rename(tmp, indexFile)
}

func loadIndex() {
    f, err := os.Open(indexFile)
    if err != nil {
        return
    }
    defer f.Close()
    json.NewDecoder(f).Decode(store)
}
