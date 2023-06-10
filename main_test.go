package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestGee(t *testing.T) {
	cmd := exec.Command("./gee", "testdata/input.go")
	got, err := cmd.Output()
	if err != nil {
		t.Fatal("failed to run command:", err)
	}
	want, err := ioutil.ReadFile("testdata/output.go")
	if err != nil {
		t.Fatal("failed to read expected output from file:", err)
	}
	if bytes.Compare(got, want) != 0 {
		t.Errorf("FAILED: want %s, got %s", want, got)
	}
}

func TestGeeIgnoreGenFiles(t *testing.T) {
	cmd := exec.Command("./gee", "testdata/input.gen.go")
	got, err := cmd.Output()

	if err != nil {
		t.Fatal("failed to run command:", err)
	}

	if bytes.Compare(got, []byte("")) != 0 {
		t.Errorf("FAILED: want %s, got %s", "", got)
	}
}
