//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// Compiles adapter shared library (lib/liblinuxcncadapter.so)
func (Build) Adapter() error {
	linuxCNCRoot := strings.TrimRight(os.Getenv("EMC2_HOME"), "/")

	// Run C++ compiler
	i := "-I" + linuxCNCRoot
	if err := sh.RunV("g++", "-c", "adapter/linuxcnc.cpp", "-fPIC",
		"-o", "adapter/linuxcnc.o",
		i+"/lib", i+"/include",
		i+"/src", i+"/src/libnml",
		i+"/src/rtapi", i+"/src/emc/tooldata",
	); err != nil {
		return fmt.Errorf("compiling adapter/linuxcnc.o: %w", err)
	}

	// Create shared library
	if err := sh.RunV("gcc", "-shared",
		"-o", "lib/liblinuxcncadapter.so",
		"adapter/linuxcnc.o",
		"-lnml", "-llinuxcnc", "-llinuxcnchal", "-ltooldata",
		"-L"+linuxCNCRoot+"/lib",
	); err != nil {
		return fmt.Errorf("linking adapter/liblinuxcncadapter.so: %w", err)
	}

	return nil
}

func (Build) Datadump() error {
	mg.Deps(Build.Adapter)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get workdir: %w", err)
	}

	env := map[string]string{
		"LD_LIBRARY_PATH": os.Getenv("LD_LIBRARY_PATH") + ":" + wd + "/lib",
	}
	if err := sh.RunWithV(env, "go", "build", "-v", "-o", "bin/datadump", "./cmd/datadump"); err != nil {
		return fmt.Errorf("compiling cmd/datadump: %w", err)
	}
	return nil
}

type Run mg.Namespace

func (Run) Datadump() error {
	mg.Deps(Build.Datadump)

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get workdir: %w", err)
	}

	env := map[string]string{
		"LD_LIBRARY_PATH": os.Getenv("LD_LIBRARY_PATH") + ":" + wd + "/lib",
	}
	if err := sh.RunWithV(env, "go", "run", "./cmd/datadump"); err != nil {
		return fmt.Errorf("running cmd/datadump: %w", err)
	}
	return nil
}
