{{ range . -}}
{{ .Name }},{{ .Version }},{{ .LicenseName }},{{ .LicenseURL }},{{ .LicensePath }}
{{ end -}}