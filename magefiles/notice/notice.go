// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package notice

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	csvPkgName     int = 0
	csvPkgVersion  int = 1
	csvLicenseName int = 2
	csvLicenseURL  int = 3
	csvLicensePath int = 4
	csvLen         int = 5
)

func GenerateNotice(notice io.Writer, csvLicenses io.Reader) error {
	licenses := []license{}
	line := 0
	scanner := bufio.NewScanner(csvLicenses)
	for scanner.Scan() {
		csvLine := strings.Split(scanner.Text(), ",")
		if len(csvLine) != csvLen {
			return fmt.Errorf(
				"csv line has incorrect number of elements"+
					": line %d: elements %d",
				line, len(csvLine),
			)
		}
		licPath := csvLine[csvLicensePath]
		noticePath := filepath.Join(filepath.Dir(licPath), "NOTICE")

		licFile, err := os.Open(licPath)
		if err != nil {
			return fmt.Errorf("opening license file %s: %w", licPath, err)
		}
		licText, err := io.ReadAll(licFile)
		if err != nil {
			return fmt.Errorf("reading license file: %w", err)
		}

		noticeFile, err := os.Open(noticePath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("opening notice file %s: %w", noticePath, err)
		}

		crText := createCopyRight(licFile, noticeFile)

		licenses = append(
			licenses, license{
				name:          csvLine[csvLicenseName],
				pkg:           csvLine[csvPkgName],
				pkgVersion:    csvLine[csvPkgVersion],
				copyrightText: crText,
				licenseText:   string(licText),
			},
		)
		line++
	}
	if err := createNotice(notice, licenses); err != nil {
		return fmt.Errorf("writing notice file: %w", err)
	}
	return nil
}

type license struct {
	// name of the license
	name          string
	pkg           string
	pkgVersion    string
	copyrightText string
	licenseText   string
}

func createNotice(notice io.Writer, lics []license) error {
	var txt bytes.Buffer

	txt.WriteString("Copyright 2023 Volvo Car Corporation\n\n")
	txt.WriteString("[github.com/volvo-cars/lingon]\n\n")
	txt.WriteString("Components:\n\n")
	for _, lic := range lics {
		txt.WriteString(
			fmt.Sprintf(
				"%s %s : %s\n", lic.pkg, lic.pkgVersion,
				lic.name,
			),
		)
	}

	txt.WriteString("\nCopyright Text:\n\n")
	for _, lic := range lics {
		txt.WriteString(fmt.Sprintf("%s %s\n", lic.pkg, lic.pkgVersion))
		txt.WriteString(lic.copyrightText + "\n")
	}

	txt.WriteString("\nLicenses:\n\n")
	for _, lic := range lics {
		txt.WriteString(lic.name + "\n")
		txt.WriteString(fmt.Sprintf("(%s %s)\n\n", lic.pkg, lic.pkgVersion))
		txt.WriteString(lic.licenseText + "\n")
		txt.WriteString("\n---\n\n")
	}

	_, err := notice.Write(txt.Bytes())
	if err != nil {
		return fmt.Errorf("writing to notice: %w", err)
	}
	return nil
}

func createCopyRight(licFile, noticeFile io.Reader) string {
	var crText strings.Builder
	licCopyrights := readCopyrightFromLicense(licFile)

	// Write copyright lines from license
	for _, cr := range licCopyrights {
		crText.WriteString("\t" + cr + "\n")
	}

	if noticeFile != nil {
		scanner := bufio.NewScanner(noticeFile)
		for scanner.Scan() {
			line := scanner.Text()
			crText.WriteString("\t" + line + "\n")
		}
	}
	return crText.String()
}

func readCopyrightFromLicense(file io.Reader) []string {
	crLines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Copyright") {
			crLines = append(crLines, line)
		}
	}
	return crLines
}
