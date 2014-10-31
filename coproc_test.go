// originally based on goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package coproc

import (
	"testing"
)

func TestPidfile(t *testing.T) {
	c := new(Group)
	c.Add(&Process{
		Name:    "test",
		Pidfile: "test.pid",
	})
	p := c.Get("test")
	err := p.Pidfile.write(100)
	if err != nil {
		t.Errorf("Error: %s.", err)
		return
	}
	ex := 100
	r := p.Pidfile.read()
	if ex != r {
		t.Errorf("Expected %#v. Result %#v\n", ex, r)
	}

	s := p.Pidfile.delete()
	if s != true {
		t.Error("Failed to remove pidfile.")
		return
	}
}

func TestProcessStart(t *testing.T) {
	c := new(Group)
	c.Add(&Process{
		Name:    "bash",
		Command: "/bin/bash",
		Args:    []string{"foo", "bar"},
		Pidfile: "echo.pid",
		Logfile: "debug.log",
		Errfile: "error.log",
		Respawn: 3,
	})
	p := c.Get("bash")
	p.start("bash")
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	p.stop()
}