package kube

import (
	"github.com/dave/jennifer/jen"
	"github.com/veggiemonk/strcase"
	"github.com/volvo-cars/lingon/pkg/internal/api"
)

func stmtKubeObjectFile(
	pkgName string,
	nameVarObj string,
	stmt *jen.Statement,
) *jen.File {
	f := jen.NewFile(pkgName)
	f.HeaderComment(headerComment)
	api.ImportKubernetesPkgAlias(f)
	f.Line()

	// var NAME = &v1.Deployment{}
	f.Var().Id(nameVarObj).Op("=").Add(stmt)
	return f
}

// appFile generates the app.go file
func (j *jamel) appFile() *jen.File {
	appFile := jen.NewFile(j.o.OutputPkgName)
	appFile.HeaderComment(headerComment)
	api.ImportKubernetesPkgAlias(appFile)
	appFile.Line()

	// var _ kube.Exporter = (*NAME)(nil)
	nameStruct := strcase.Pascal(j.o.AppName)
	appFile.Comment("validate the struct implements the interface")
	appFile.Var().Op("_").Qual(kubeAppPkgPath, "Exporter").
		Op("=").Parens(jen.Op("*").Id(nameStruct)).Parens(jen.Nil())

	// struct
	structCode := stmtStruct(nameStruct, j.kubeAppStructCode)
	appFile.Commentf("%s contains kubernetes manifests", nameStruct)
	appFile.Add(structCode)
	appFile.Line()

	// NewApp
	newApp := stmtNewApp(nameStruct, j.nameFieldVar)
	appFile.Commentf("New creates a new %s", nameStruct)
	appFile.Add(newApp)
	appFile.Line().Line()

	if j.o.AddMethods {
		addMethods(appFile, nameStruct)
		// add Apply function
		appFile.Add(stmtApplyFunc())
	}

	// add P function to convert T to *T
	appFile.
		Commentf("P converts T to *T, useful for basic types").Line().
		Func().Id("P").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("t").Id("T"),
			// return type
		).Op("*").Id("T").Block(
		jen.Return(jen.Op("&").Id("t")),
	)
	appFile.Line().Line()

	return appFile
}

// stmtStruct generates the statements for struct containing the kubernetes objects (and kube.App)
func stmtStruct(
	nameStruct string,
	kubeAppStructCode map[string]*jen.Statement,
) *jen.Statement {
	// type NAME struct {
	return jen.Type().Id(nameStruct).StructFunc(
		func(g *jen.Group) {
			// add kube.App to the struct
			g.Qual(kubeAppPkgPath, embeddedStructName).Line()

			// add all the objects to the struct
			keys := orderedKeys(kubeAppStructCode)
			for _, k := range keys {
				v := kubeAppStructCode[k]
				g.Id(k).Add(jen.Op("*").Add(v))
			}
		},
	)
}

// stmtNewApp generates the statements for the New function instantiating the struct
func stmtNewApp(
	nameStruct string,
	nameFieldVar map[string]string,
) *jen.Statement {
	// func New() *NAME {
	return jen.Func().Id("New").Params().Op("*").Id(nameStruct).Block(
		jen.Return().Op("&").Id(nameStruct).ValuesFunc(
			func(g *jen.Group) {
				keys := orderedKeys(nameFieldVar)
				for _, nameField := range keys {
					nameVar := nameFieldVar[nameField]
					// field: kube obj var
					g.Line().Id(nameField).Op(":").Id(nameVar)
				}
			},
		),
	)
}

// addMethods adds the Apply and Export methods to the struct in the app.go file
func addMethods(f *jen.File, nameStruct string) *jen.File {
	// Apply
	f.Comment("Apply applies the kubernetes objects to the cluster").Line().
		Func().Params(jen.Id("a").Op("*").Id(nameStruct)).Id("Apply").
		Params(jen.Id("ctx").Qual("context", "Context")).
		Params(jen.Error()).
		BlockFunc(
			func(g *jen.Group) {
				g.Return(
					jen.Id("Apply").
						Call(jen.Id("ctx"), jen.Id("a")),
				)
			},
		)

	f.Line().Line()

	// Export
	f.Comment("Export exports the kubernetes objects to YAML files in the given directory").Line().
		Func().Params(jen.Id("a").Op("*").Id(nameStruct)).Id("Export").
		Params(jen.Id("dir").String()).
		Params(jen.Error()).
		BlockFunc(
			func(g *jen.Group) {
				g.Return(
					jen.Qual(kubeAppPkgPath, "Export").
						Call(jen.Id("a"), jen.Id("dir")),
				)
			},
		)

	f.Line().Line()
	return f
}

// stmtApplyFunc generates the statements for the Apply function
func stmtApplyFunc() *jen.Statement {
	return jen.
		Comment("Apply applies the kubernetes objects contained in Exporter to the cluster").Line().
		Func().Id("Apply").Params(
		jen.Id("ctx").Qual("context", "Context"),
		jen.Id("km").Qual(kubeAppPkgPath, "Exporter"),
	).Params(
		jen.Error(),
	).Block(
		jen.Id("cmd").Op(":=").Qual("os/exec", "CommandContext").Call(
			jen.Id("ctx"),
			jen.Lit("kubectl"),
			jen.Lit("apply"),
			jen.Lit("-f"),
			jen.Lit("-"),
		),
		jen.Id("cmd").Dot("Env").Op("=").Qual(
			"os",
			"Environ",
		).Call().Comment("inherit environment in case we need to use kubectl from a container"),
		jen.List(
			jen.Id("stdin"),
			jen.Err(),
		).Op(":=").Id("cmd").Dot("StdinPipe").Call().Comment("pipe to pass data to kubectl"),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
		jen.Line(),

		jen.Id("cmd").Dot("Stdout").Op("=").Qual("os", "Stdout"),
		jen.Id("cmd").Dot("Stderr").Op("=").Qual("os", "Stderr"),
		jen.Line(),
		jen.Go().Func().Params().Block(
			jen.Defer().Func().Params().Block(
				jen.Err().Op("=").Qual(
					"errors",
					"Join",
				).Call(
					jen.Err(),
					jen.Id("stdin").Dot("Close").Call(),
				),
			).Call(),
			jen.If(
				jen.Id("errEW").Op(":=").Qual(
					kubeAppPkgPath,
					"ExportWriter",
				).Call(
					jen.Id("km"),
					jen.Id("stdin"),
				).Op(";").Id("errEW").Op("!=").Nil(),
			).Block(
				jen.Err().Op("=").Qual("errors", "Join").Call(
					jen.Err(),
					jen.Id("errEW"),
				),
			),
		).Call(),
		jen.Line(),
		jen.If(jen.Id("errS").Op(":=").Id("cmd").Dot("Start").Call().Op(";").Id("errS").Op("!=").Nil()).Block(
			jen.Return(
				jen.Qual("errors", "Join").Call(
					jen.Err(),
					jen.Id("errS"),
				),
			),
		),
		jen.Line(),
		jen.Comment("waits for the command to exit and waits for any copying"),
		jen.Comment("to stdin or copying from stdout or stderr to complete"),
		jen.Return(
			jen.Qual("errors", "Join").Call(
				jen.Err(),
				jen.Id("cmd").Dot("Wait").Call(),
			),
		),
	)
}
