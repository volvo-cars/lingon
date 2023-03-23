package kube

// Exporter interfaces for kubernetes objects defined in a Go structs
type Exporter interface {
	Klamydia()
}

var _ Exporter = (*App)(nil)

// App struct is meant to be embedded in other structs
// to specify that they are a set of kubernetes manifests
type App struct{}

// IamGroot is a dummy method to make sure that App implements Exporter
func (a *App) Klamydia() {}
