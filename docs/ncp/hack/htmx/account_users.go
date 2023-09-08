package htmx

import (
	_ "embed"
)

//go:embed account_users.html
var accountUsersHTML string

const (
	accountUsersTmplName = "account_users"
)

// var accountUsersTmpl = template.Must(
// 	baseTmpl.New(accountUsersTmplName).Parse(accountUsersHTML),
// )
