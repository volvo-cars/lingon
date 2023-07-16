// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package score

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/eidolon/wordwrap"
	"github.com/fatih/color"
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

func Score(r io.Reader) (*scorecard.Scorecard, error) {
	reader := &inputReader{Reader: r}

	cnf := config.Configuration{AllFiles: []domain.NamedReader{reader}}
	p, err := parser.New()
	if err != nil {
		return nil, err
	}
	parsed, err := p.ParseFiles(cnf)
	if err != nil {
		return nil, err
	}

	card, err := score.Score(parsed, cnf)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func RenderScoreCard2(
	card *scorecard.Scorecard,
	w io.Writer,
	useColor bool,
) error {
	color.NoColor = !useColor

	for _, key := range orderedKeys(*card) {
		scoredObject := (*card)[key]
		if _, err := fmt.Fprintf(
			w,
			"%s/%s %s\n",
			scoredObject.TypeMeta.APIVersion,
			scoredObject.TypeMeta.Kind,
			scoredObject.ObjectMeta.Name,
		); err != nil {
			return err
		}

		if scoredObject.ObjectMeta.Namespace != "" {
			if _, err := fmt.Fprintf(
				w,
				" in %s\n",
				scoredObject.ObjectMeta.Namespace,
			); err != nil {
				return err
			}
		}

		checks := scoredObject.Checks
		sort.SliceStable(
			checks, func(i, j int) bool {
				return checks[i].Check.ID < checks[j].Check.ID
			},
		)
		for _, c := range checks {
			r := writeStep(c, 0)
			if _, err := io.Copy(w, r); err != nil {
				return fmt.Errorf("copy output: %w", err)
			}
		}
	}
	return nil
}

func writeStep(card scorecard.TestScore, verboseOutput int) io.Reader {
	bs := bytes.NewBufferString("")

	if card.Skipped && verboseOutput < 2 {
		return bs
	}

	var col color.Attribute

	switch {
	case card.Skipped || card.Grade >= scorecard.GradeAllOK:
		// Higher than or equal to --threshold-ok
		col = color.FgGreen

		// If verbose output is disabled, skip OK items in the output
		if verboseOutput == 0 {
			return bs
		}

	case card.Grade >= scorecard.GradeWarning:
		// Higher than or equal to --threshold-warning
		col = color.FgYellow
	default:
		// All lower than both --threshold-ok and --threshold-warning are critical
		col = color.FgRed
	}

	if card.Skipped {
		_, _ = color.New(col).Fprintf(
			bs,
			"    [SKIPPED] %s\n",
			card.Check.Name,
		)
	} else {
		_, _ = color.New(col).Fprintf(
			bs,
			"    [%s] %s\n",
			card.Grade.String(),
			card.Check.Name,
		)
	}

	for _, comment := range card.Comments {
		_, _ = fmt.Fprintf(bs, "        · ")

		if len(comment.Path) > 0 {
			_, _ = fmt.Fprintf(bs, "%s -> ", comment.Path)
		}

		_, _ = fmt.Fprint(bs, comment.Summary)

		if len(comment.Description) > 0 {
			wrapWidth := 100
			// if wrapWidth < 40 {
			// 	wrapWidth = 40
			// }
			wrapper := wordwrap.Wrapper(wrapWidth, false)
			wrapped := wrapper(comment.Description)
			_, _ = fmt.Fprintln(bs)
			_, _ = fmt.Fprint(
				bs,
				wordwrap.Indent(wrapped, strings.Repeat(" ", 12), false),
			)
		}

		if len(comment.DocumentationURL) > 0 {
			_, _ = fmt.Fprintln(bs)
			_, _ = fmt.Fprintf(
				bs,
				"%sMore information: %s",
				strings.Repeat(" ", 12),
				comment.DocumentationURL,
			)
		}

		_, _ = fmt.Fprintln(bs)
	}

	return bs
}

func orderedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
