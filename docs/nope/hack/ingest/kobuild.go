package ingest

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/exp/slog"
	"golang.org/x/tools/txtar"
)

type Schema struct {
	// Name of the schema
	Name string
	// Txtar of the protobuf schema
	Txtar []byte
}

type KoBuildConfig struct {
	// Workdir is the directory where the build will be executed.
	Workdir string
	// GoServiceTxtar is the txtar of the Go service.
	GoServiceTxtar []byte
	// Schema is the schema to build.
	Schema Schema
}

func KoBuild(cfg KoBuildConfig) error {
	if err := UnpackTxtar(cfg.Workdir, cfg.GoServiceTxtar); err != nil {
		return fmt.Errorf("unpacking txtar: %w", err)
	}

	if err := protoBuild(protoBuildConfig{
		Schema:  cfg.Schema,
		Workdir: filepath.Join(cfg.Workdir, "proto"),
		Out:     filepath.Join(cfg.Workdir, "schema"),
		Pkg:     "ingestion/schema",
	}); err != nil {
		return fmt.Errorf("building proto: %w", err)
	}

	cmd := exec.Command("ko", "build", "--push=false", "--local", "--bare")
	// cmd.Env = []string{fmt.Sprintf("KO_DOCKER_REPO=%s", "platypus/ingestion")}
	cmd.Dir = cfg.Workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running ko: %w", err)
	}

	return nil
}

type protoBuildConfig struct {
	Schema Schema
	// Workdir is the directory where the protobuf files from the schema will be
	// written to.
	Workdir string
	// Out is the directory where the generated Go code will be written to.
	Out string
	// Pkg is the package name of the generated Go code.
	Pkg string
}

func protoBuild(cfg protoBuildConfig) error {

	if err := UnpackTxtar(cfg.Workdir, cfg.Schema.Txtar); err != nil {
		return fmt.Errorf("unpacking txtar: %w", err)
	}

	bufTemplate := filepath.Join(cfg.Workdir, "buf.gen.yaml")
	if err := os.WriteFile(
		bufTemplate,
		[]byte(bufGenerateYAML(cfg.Pkg, cfg.Out)),
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("writing buf.gen.yaml: %w", err)
	}

	// Compile the protobuf files
	cmd := exec.Command("buf", "generate", "--template", bufTemplate, cfg.Workdir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	slog.Info("generating protobuf files", "exec", cmd.String())

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running buf: %w", err)
	}

	return nil
}

// PackTxtar takes a list of files and creates a [txtar.Archive] for them.
func PackTxtar(dir string, files []string) (*txtar.Archive, error) {
	var archive txtar.Archive
	for _, file := range files {
		path := filepath.Join(dir, file)
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading proto file %s: %w", path, err)
		}
		ttFile := txtar.File{
			Name: file,
			Data: contents,
		}
		archive.Files = append(archive.Files, ttFile)
	}
	return &archive, nil
}

func UnpackTxtar(dir string, contents []byte) error {
	ar := txtar.Parse(contents)
	if len(ar.Files) == 0 {
		return fmt.Errorf("txtar archive contains no files")
	}

	// Write the files to disk
	for _, tFile := range ar.Files {
		path := filepath.Join(dir, tFile.Name)
		fileDir := filepath.Dir(path)
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return fmt.Errorf("creating directory for file %s: %w", path, err)
		}
		if err := os.WriteFile(path, tFile.Data, os.ModePerm); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
	}
	return nil
}

func bufGenerateYAML(pkg string, out string) string {
	return fmt.Sprintf(`version: v1
managed:
  enabled: true
  go_package_prefix:
    default: %s

plugins:
- plugin: buf.build/protocolbuffers/go:v1.31.0
  out: %s
  opt:
  - paths=source_relative
`, pkg, out)
}
