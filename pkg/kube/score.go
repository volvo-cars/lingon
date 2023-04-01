package kube

import (
	"bytes"
	"io"

	"github.com/rogpeppe/go-internal/txtar"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/scorecard"
)

type inputReader struct {
	io.Reader
}

func (inputReader) Name() string {
	return "input"
}

func Txtar2Reader(b *bytes.Buffer) *bytes.Buffer {
	ar := txtar.Parse(b.Bytes())
	txtar.Format(ar)
	var buf bytes.Buffer
	for _, f := range ar.Files {
		buf.WriteString("\n\n---\n")
		buf.WriteString("# " + f.Name + "\n")
		buf.Write(f.Data)
	}
	return &buf
}

func Score(r io.Reader) (*scorecard.Scorecard, error) {
	reader := &inputReader{
		Reader: r,
	}

	cnf := config.Configuration{
		AllFiles: []domain.NamedReader{reader},
	}
	p, err := parser.New()
	if err != nil {
		return nil, err
	}
	allObjs, err := p.ParseFiles(cnf)
	if err != nil {
		return nil, err
	}

	card, err := score.Score(allObjs, cnf)
	if err != nil {
		return nil, err
	}
	return card, nil
}
