package main

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"
)

//go:embed badge.svg.tpl
var tplFS embed.FS

func CoverP() {
	coverOutput := "cover.out"
	coverMode := "count" // see `go help testflag` for more info
	iferr(
		Go(
			"test",
			recDir,
			"-coverprofile="+coverOutput,
			"-covermode="+coverMode,
		))

	var buf bytes.Buffer
	iferr(coverPct(&buf, coverOutput))
	pct, err := parseCovPct(&buf)
	iferr(err)
	s, err := badgetpl(pct)
	iferr(err)
	iferr(writeBadge(".github/coverage.svg", []byte(s)))
	fmt.Printf("✅ coverage: %v%%\n", pct)
}

const (
	brightgreen = "#4c1"
	green       = "#97ca00"
	yellow      = "#dfb317"
	yellowgreen = "#a4a61d"
	orange      = "#fe7d37"
	red         = "#e05d44"
)

func badgetpl(pct float64) (string, error) {
	color := ""
	switch {
	case pct < 50:
		color = orange
	case pct < 60:
		color = yellow
	case pct < 70:
		color = yellowgreen
	case pct < 80:
		color = green
	case pct < 90:
		color = brightgreen
	default:
		color = red
	}
	t := template.Must(template.ParseFS(tplFS, "badge.svg.tpl"))
	data := struct {
		Percentage float64
		Color      string
	}{Percentage: pct, Color: color}

	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func writeBadge(path string, b []byte) error {
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), 0o755)
	}
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	_, err = f.Write(b)
	return err
}

func parseCovPct(r io.Reader) (float64, error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		if strings.HasPrefix(t, "total:") {
			f := strings.Fields(t)
			if f[0] == "total:" && len(f) == 3 {
				return strconv.ParseFloat(strings.TrimRight(f[2], "%"), 64)
			}
			break
		}
	}
	if err := s.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "scan:", err)
	}
	return 0.0, fmt.Errorf("no coverage found")
}

func coverPct(w io.Writer, o string) error {
	cmd := exec.Command("go", "tool", "cover", "-func="+o) //nolint:gosec
	slog.Info("exec", slog.String("cmd", cmd.String()))
	defer slog.Info("done", slog.String("cmd", cmd.String()))

	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.Stderr.Sync()
		return fmt.Errorf("go: %s", err)
	}
	return nil
}
