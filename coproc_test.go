// originally based on goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package coproc

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
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

func TestWorkdir(t *testing.T) {
	wd, err := ioutil.TempDir("/tmp", "TestWorkdir")
	if err != nil {
		t.Errorf("Unexpected error creating temporary directory")
	}
	logFile := newFile(t, "/tmp")
	c := new(Group)
	c.Add(&Process{
		Name:    "pwd",
		Command: "/bin/pwd",
		Workdir: wd,
		Pidfile: "TestWorkdir.pid",
		Logfile: logFile,
	})
	p := c.Get("pwd")
	p.start("pwd")
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	p.stop()
	assertFileContains(t, logFile, wd)
}

func TestEnv(t *testing.T) {
	logFile := newFile(t, "/tmp")
	c := new(Group)
	c.Add(&Process{
		Name:    "echofoo",
		Command: "/bin/bash",
		Args:    []string{"-c", "echo $FOO"},
		Env:     []string{"FOO=BAR"},
		Pidfile: "TestEnv.pid",
		Logfile: logFile,
	})
	p := c.Get("echofoo")
	p.start("echofoo")
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	time.Sleep(100 * time.Millisecond)
	p.stop()
	time.Sleep(100 * time.Millisecond)
	assertFileContains(t, logFile, "BAR")
}

func assertFileContains(t *testing.T, path string, substr string) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Unexpected error reading file %s: %#v", path, err)
	}
	s := string(dat)
	if substr != "" && !strings.Contains(s, substr) {
		t.Fatalf("File %s contents\n%s\n  does not contain\n%s\n", path, s, substr)
	}
}

func newFile(t *testing.T, dir string) string {
	tf, err := ioutil.TempFile(dir, "")
	if err != nil {
		t.Errorf("Unexpected error creating temporary file")
	}
	return tf.Name()
}
