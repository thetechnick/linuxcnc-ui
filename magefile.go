//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	protocVersion               = "3.19.1"
	protocGenGoVersion          = "1.25.0"
	protocGenGoGRPCVersion      = "1.0.1"
	protocGenGRPCGatewayVersion = "2.1.0"
	goimportsVersion            = "0.1.5"
	golangCILintVersion         = "1.43.0"
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

func Generate() {
	mg.Deps(Build.Proto)
}

type Dependency mg.Namespace

func (Dependency) All() {
	mg.Deps(
		Dependency.Protoc,
		Dependency.GRPCGatewayIncludes,
		mg.F(Dependency.GoInstall, "protoc-gen-go", "google.golang.org/protobuf/cmd/protoc-gen-go", protocGenGoVersion),
		mg.F(Dependency.GoInstall, "protoc-gen-go-grpc", "google.golang.org/grpc/cmd/protoc-gen-go-grpc", protocGenGoGRPCVersion),
		mg.F(Dependency.GoInstall, "protoc-gen-grpc-gateway", "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway", protocGenGRPCGatewayVersion),
		mg.F(Dependency.GoInstall, "goimports", "golang.org/x/tools/cmd/goimports", goimportsVersion),
		mg.F(Dependency.GoInstall, "golangci-lint", "github.com/golangci/golangci-lint/cmd/golangci-lint", golangCILintVersion),
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

	needsRebuild, err := checkBinDependencyNeedsRebuild("protoc", protocVersion)
	if err != nil {
		return err
	}
	if !needsRebuild {
		return nil
	}

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

	// Bump timestamp
	currentTime := time.Now().Local()
	err = os.Chtimes(depsDir+"/bin/protoc", currentTime, currentTime)
	if err != nil {
		return fmt.Errorf("bump change date: %w", err)
	}

	return nil
}

func (Dependency) GRPCGatewayIncludes() error {
	mg.Deps(Dependency.Dirs)

	// Remember version
	versionFile := depsDir + "/versions/grpc-gateway-includes/v" + protocGenGRPCGatewayVersion
	if err := ensureFile(versionFile); err != nil {
		return fmt.Errorf("ensure file: %w", err)
	}

	// Check if rebuild is needed
	rebuild, err := target.Path(".deps/include/google/api", versionFile)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}
	if !rebuild {
		return nil
	}

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

func (Dependency) GoInstall(bin, packageURl, version string) error {
	mg.Deps(Dependency.Dirs)

	needsRebuild, err := checkBinDependencyNeedsRebuild(bin, version)
	if err != nil {
		return err
	}
	if !needsRebuild {
		return nil
	}

	url := packageURl + "@v" + version
	if err := sh.RunWithV(map[string]string{
		"GOBIN": depsBinDir,
	}, mg.GoCmd(),
		"install", url,
	); err != nil {
		return fmt.Errorf("install %s: %w", url, err)
	}
	return nil
}

// ensure a file and it's file path exist.
func ensureFile(file string) error {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		f, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("creating file %s: %w", file, err)
		}
		defer f.Close()
		return nil
	}
	if err != nil {
		return fmt.Errorf("checking file %s: %w", file, err)
	}
	return nil
}

func checkBinDependencyNeedsRebuild(thing, version string) (needsRebuild bool, err error) {
	versionFile := ".deps/versions/v" + version
	if err := ensureFile(versionFile); err != nil {
		return false, fmt.Errorf("ensure file: %w", err)
	}

	rebuild, err := target.Path(".deps/bin/"+thing, versionFile)
	if err != nil {
		return false, fmt.Errorf("check: %w", err)
	}
	if !rebuild {
		return false, nil
	}

	return true, nil
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
	mg.Deps(Dependency.All)

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

func (Build) Datadump() {
	mg.SerialDeps(
		Dependency.All,
		Build.Adapter,
		mg.F(Build.Cmd, "datadump"),
	)
}

func (Build) APIServer() {
	mg.SerialDeps(
		Dependency.All,
		mg.F(Build.Cmd, "apiserver"),
	)
}

func (Build) Cmd(cmd string) error {
	os.Rename("cmd/datadump/mock.cpp", ".cache/mock.ccp")
	defer os.Rename(".cache/mock.ccp", "cmd/datadump/mock.cpp")

	env := map[string]string{
		"LD_LIBRARY_PATH": os.Getenv("LD_LIBRARY_PATH") + ":" + workDir + "/lib",
		"CGO_LDFLAGS":     "-Llib -llinuxcncadapter",
	}
	if err := sh.RunWithV(env, "go", "build", "-v", "-o", "bin/"+cmd, "./cmd/"+cmd+"/main.go"); err != nil {
		return fmt.Errorf("compiling cmd/%s: %w", cmd, err)
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
	if err := sh.RunWithV(env, "./bin/datadump"); err != nil {
		return fmt.Errorf("running cmd/datadump: %w", err)
	}
	return nil
}
