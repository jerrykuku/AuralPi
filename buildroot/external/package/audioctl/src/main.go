package main

import (
	"embed"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
)

const (
	serviceAirPlay = "shairport"
	serviceRoon    = "roonbridge"
)

//go:embed static/*
var staticFiles embed.FS

type ServiceRunner interface {
	Start(name string) error
	Stop(name string) error
	Status(name string) (string, error)
}

type Controller struct {
	runner ServiceRunner
}

type StatusResponse struct {
	AirPlay    string `json:"airplay"`
	RoonBridge string `json:"roonbridge"`
}

func NewController(runner ServiceRunner) *Controller {
	return &Controller{runner: runner}
}

func (c *Controller) StartAirPlay() error {
	if err := c.runner.Stop(serviceRoon); err != nil {
		return err
	}
	if err := c.runner.Stop(serviceAirPlay); err != nil {
		return err
	}
	return c.runner.Start(serviceAirPlay)
}

func (c *Controller) StopAirPlay() error {
	return c.runner.Stop(serviceAirPlay)
}

func (c *Controller) StartRoon() error {
	if err := c.runner.Stop(serviceAirPlay); err != nil {
		return err
	}
	if err := c.runner.Stop(serviceRoon); err != nil {
		return err
	}
	return c.runner.Start(serviceRoon)
}

func (c *Controller) StopRoon() error {
	return c.runner.Stop(serviceRoon)
}

func (c *Controller) Status() (StatusResponse, error) {
	airplay, err := c.runner.Status(serviceAirPlay)
	if err != nil {
		return StatusResponse{}, err
	}
	roon, err := c.runner.Status(serviceRoon)
	if err != nil {
		return StatusResponse{}, err
	}
	return StatusResponse{AirPlay: airplay, RoonBridge: roon}, nil
}

func main() {
	ctl := NewController(initScriptRunner{})
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFiles)))
	mux.HandleFunc("/api/status", handleStatus(ctl))
	mux.HandleFunc("/api/airplay/start", postOnly(handleAction(ctl.StartAirPlay)))
	mux.HandleFunc("/api/airplay/stop", postOnly(handleAction(ctl.StopAirPlay)))
	mux.HandleFunc("/api/roon/start", postOnly(handleAction(ctl.StartRoon)))
	mux.HandleFunc("/api/roon/stop", postOnly(handleAction(ctl.StopRoon)))

	log.Println("audioctl listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleStatus(ctl *Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		status, err := ctl.Status()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, status)
	}
}

func handleAction(action func() error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := action(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	}
}

func postOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

type initScriptRunner struct{}

func (initScriptRunner) Start(name string) error {
	return runInitScript(name, "start")
}

func (initScriptRunner) Stop(name string) error {
	return runInitScript(name, "stop")
}

func (initScriptRunner) Status(name string) (string, error) {
	err := runInitScript(name, "status")
	if err == nil {
		return "running", nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return "stopped", nil
	}
	return "unknown", err
}

func runInitScript(name string, action string) error {
	path := "/etc/init.d/S50shairport"
	if name == serviceRoon {
		path = "/etc/init.d/S55roonbridge"
	}
	cmd := exec.Command(path, action)
	return cmd.Run()
}
