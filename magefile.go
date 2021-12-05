//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	protocVersion               = "3.19.1"
	protocGenGoVersion          = "1.25.0"
	protocGenGoGRPCVersion      = "1.0.1"
	protocGenGRPCGatewayVersion = "2.1.0"
)

var (
	workDir        string
	depsDir        string
	depsBinDir     string
	depsIncludeDir string
)

func init() {
	workDir, _ = os.Getwd()
	depsDir = workDir + "/.deps"
	depsBinDir = depsDir + "/bin"
	depsIncludeDir = depsDir + "/include"
}

type Dependency mg.Namespace

func (Dependency) All() {
	mg.Deps(
		Dependency.Protoc,
		Dependency.GRPCGatewayIncludes,
		mg.F(Dependency.GoInstall, "google.golang.org/protobuf/cmd/protoc-gen-go@v"+protocGenGoVersion),
		mg.F(Dependency.GoInstall, "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v"+protocGenGoGRPCVersion),
		mg.F(Dependency.GoInstall, "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v"+protocGenGRPCGatewayVersion),
	)
}

func (Dependency) Dirs() error {
	// ensure deps dir
	if err := os.MkdirAll(".cache", os.ModePerm); err != nil {
		return fmt.Errorf("creating deps dir: %w", err)
	}
	if err := os.MkdirAll(".deps/bin", os.ModePerm); err != nil {
		return fmt.Errorf("creating deps/bin dir: %w", err)
	}
	if err := os.MkdirAll(".deps/include/google", os.ModePerm); err != nil {
		return fmt.Errorf("creating deps/include dir: %w", err)
	}
	return nil
}

func (Dependency) Protoc() error {
	mg.Deps(Dependency.Dirs)

	// Tempdir
	tempDir, err := os.MkdirTemp(".cache", "")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download
	if err := sh.Run(
		"curl", "-L", "--fail",
		"-o", tempDir+"/protoc.zip",
		fmt.Sprintf(
			// https://github.com/protocolbuffers/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip
			"https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-linux-x86_64.zip",
			protocVersion, protocVersion,
		),
	); err != nil {
		return fmt.Errorf("downloading protoc: %w", err)
	}
	// Unzip
	if err := sh.Run(
		"unzip", "-qq", tempDir+"/protoc.zip", "-d", tempDir); err != nil {
		return fmt.Errorf("unzipping protoc: %w", err)
	}

	// Move
	if err := os.RemoveAll(depsDir + "/include/google/protobuf"); err != nil {
		return fmt.Errorf("clean protobuf imports: %w", err)
	}
	if err := os.Rename(tempDir+"/include/google/protobuf", depsIncludeDir+"/google/protobuf"); err != nil {
		return fmt.Errorf("move include: %w", err)
	}
	if err := os.Rename(tempDir+"/bin/protoc", depsDir+"/bin/protoc"); err != nil {
		return fmt.Errorf("move protoc: %w", err)
	}
	return nil
}

func (Dependency) GRPCGatewayIncludes() error {
	mg.Deps(Dependency.Dirs)

	// Tempdir
	tempDir, err := os.MkdirTemp(".cache", "")
	if err != nil {
		return fmt.Errorf("temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Git clone
	if err := sh.RunV(
		"git", "clone", "https://github.com/grpc-ecosystem/grpc-gateway",
		"--depth=1", "--branch=v"+protocGenGRPCGatewayVersion,
		tempDir,
	); err != nil {
		return fmt.Errorf("git checkout: %w", err)
	}

	// Move
	if err := os.RemoveAll(depsIncludeDir + "/google/api"); err != nil {
		return fmt.Errorf("clean protobuf gateway imports: %w", err)
	}
	if err := os.RemoveAll(depsIncludeDir + "/google/rpc"); err != nil {
		return fmt.Errorf("clean protobuf imports: %w", err)
	}
	if err := os.Rename(tempDir+"/third_party/googleapis/google/api", depsIncludeDir+"/google/api"); err != nil {
		return fmt.Errorf("move includes: %w", err)
	}
	if err := os.Rename(tempDir+"/third_party/googleapis/google/rpc", depsIncludeDir+"/google/rpc"); err != nil {
		return fmt.Errorf("move includes: %w", err)
	}
	return nil
}

func (Dependency) GoInstall(pack string) error {
	mg.Deps(Dependency.Dirs)

	if err := sh.RunWithV(map[string]string{
		"GOBIN": depsBinDir,
	}, mg.GoCmd(),
		"install", pack,
	); err != nil {
		return fmt.Errorf("install %s: %w", pack, err)
	}
	return nil
}

type Build mg.Namespace

func (Build) Proto() error {
	// mg.Deps(Dependency.All)

	matches, err := filepath.Glob("api/**/*.proto")
	if err != nil {
		return fmt.Errorf("glob *.proto files: %w", err)
	}

	for _, match := range matches {
		err := sh.RunWithV(
			map[string]string{
				"PATH": depsBinDir + ":$PATH",
			}, "protoc",
			"--go_out="+filepath.Dir(match), "--go_opt=paths=source_relative",
			"--go-grpc_out="+filepath.Dir(match), "--go-grpc_opt=paths=source_relative",
			"--grpc-gateway_out="+filepath.Dir(match),
			"--grpc-gateway_opt", "paths=source_relative",
			"--grpc-gateway_opt", "generate_unbound_methods=true",
			"-I"+depsIncludeDir,
			"-Iapi/v1",
			"-I"+filepath.Dir(match),
			match,
		)
		if err != nil {
			return fmt.Errorf("running protoc for %q: %w", match, err)
		}
	}

	return nil
}

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
