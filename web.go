package main

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

//go:embed files
var efs embed.FS

type server struct {
	hs  *http.Server
	cli *client

	tplIndex *template.Template

	acts map[int]func(context.Context) error
	devs []Device
	mu   sync.Mutex
}

func newServer(cli *client) *server {
	f, err := fs.Sub(efs, "files")
	if err != nil {
		panic(err)
	}

	s := &server{
		tplIndex: template.Must(template.ParseFS(f, "index.html")),
		cli:      cli,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/act/", s.act)
	mux.Handle("/assets/", http.FileServer(http.FS(f)))
	mux.Handle("/favicon.svg", http.FileServer(http.FS(f)))

	addr, ok := os.LookupEnv("ADDR")
	if !ok || addr == "" {
		addr = ":8080"
	}

	s.hs = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s
}

func (s *server) start() {
	log.Fatal(s.hs.ListenAndServe())
}

func (s *server) updateDevices(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// TODO: update sometimes?
	if s.devs != nil {
		return
	}

	devs := s.cli.list(ctx)

	id := 1
	acts := make(map[int]func(context.Context) error)
	for i := range devs {
		for j := range devs[i].Actions {
			devs[i].Actions[j].ID = id
			acts[id] = devs[i].Actions[j].Act
			id++
		}
	}

	s.devs = devs
	s.acts = acts
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	s.updateDevices(context.TODO())

	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.tplIndex.Execute(w, map[string]any{
		"Devs": s.devs,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
	}
}

func (s *server) act(w http.ResponseWriter, r *http.Request) {
	output := func(format string, a ...any) {
		msg := fmt.Sprintf(format, a...)
		log.Printf(msg)
		w.Write([]byte(msg))
	}

	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		output("error on parsing action ID: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.devs == nil {
		w.WriteHeader(http.StatusBadRequest)
		output("not initialized")
		return
	}

	act, ok := s.acts[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		output("unknown action ID: %d", id)
		return
	}
	if act == nil {
		w.WriteHeader(http.StatusInternalServerError)
		output("no action is registered for action ID: %d", id)
		return
	}
	if err := act(context.TODO()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		output("failed to act: %v", err)
		return
	}
}
