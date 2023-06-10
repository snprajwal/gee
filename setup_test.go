package main

import (
	"os"
	"os/exec"
	"testing"
)

/*
	Setup Test
	This file will run before any other test will run.

	It must named be as `setup_test.go`

    We are going to build the `gee` binary here for testing
*/

func TestMain(m *testing.M) {

    workingDir, err := os.Getwd();

    if err != nil {
        panic(err)
    }

    cmd := exec.Command("go", "build", "-o", "gee", workingDir)

    _, err = cmd.Output()

    if err != nil {
        panic(err)
    }

    os.Exit(m.Run()) // Run all test now
}
