package main

import "testing"

func TestStartAirPlayStopsRoonFirst(t *testing.T) {
	recorder := &recordingServiceRunner{}
	ctl := NewController(recorder)

	if err := ctl.StartAirPlay(); err != nil {
		t.Fatalf("StartAirPlay returned error: %v", err)
	}

	want := []string{"roonbridge:stop", "shairport:stop", "shairport:start"}
	assertCommands(t, recorder.commands, want)
}

func TestStartRoonStopsAirPlayFirst(t *testing.T) {
	recorder := &recordingServiceRunner{}
	ctl := NewController(recorder)

	if err := ctl.StartRoon(); err != nil {
		t.Fatalf("StartRoon returned error: %v", err)
	}

	want := []string{"shairport:stop", "roonbridge:stop", "roonbridge:start"}
	assertCommands(t, recorder.commands, want)
}

func TestStatusReportsBothServices(t *testing.T) {
	recorder := &recordingServiceRunner{
		status: map[string]string{
			"shairport":  "running",
			"roonbridge": "stopped",
		},
	}
	ctl := NewController(recorder)

	status, err := ctl.Status()
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}

	if status.AirPlay != "running" {
		t.Fatalf("AirPlay status = %q, want running", status.AirPlay)
	}
	if status.RoonBridge != "stopped" {
		t.Fatalf("RoonBridge status = %q, want stopped", status.RoonBridge)
	}
}

type recordingServiceRunner struct {
	commands []string
	status   map[string]string
}

func (r *recordingServiceRunner) Start(name string) error {
	r.commands = append(r.commands, name+":start")
	return nil
}

func (r *recordingServiceRunner) Stop(name string) error {
	r.commands = append(r.commands, name+":stop")
	return nil
}

func (r *recordingServiceRunner) Status(name string) (string, error) {
	if r.status == nil {
		return "unknown", nil
	}
	if value, ok := r.status[name]; ok {
		return value, nil
	}
	return "unknown", nil
}

func assertCommands(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("commands = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("commands = %#v, want %#v", got, want)
		}
	}
}
