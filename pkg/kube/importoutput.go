package kube

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/veggiemonk/strcase"
	"github.com/volvo-cars/go-terriyaki/pkg/meta"
)

const (
	embeddedStructName = "App"
	kubeAppPkgPath     = "github.com/volvo-cars/go-terriyaki/pkg/kube"
)

func (j *jamel) render() error {
	if err := os.MkdirAll(j.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// convert to go code
	if err := j.generateGo(); err != nil {
		return fmt.Errorf("generate go: %w", err)
	}

	// render all the kubernetes objects
	if j.o.GroupByKind {
		if err := j.renderFileByKind(); err != nil {
			return fmt.Errorf("save by kind: %w", err)
		}
	} else {
		if err := j.renderFileByName(); err != nil {
			return fmt.Errorf("save by name: %w", err)
		}
	}

	// app.go with kubeapp struct
	appFile := j.appFile()
	filename := filepath.Join(j.o.OutputDir, "app.go")
	if _, err := j.o.GoCodeWriter.Write([]byte("-- " + filename + " --\n")); err != nil {
		return err
	}
	if err := appFile.Render(j.o.GoCodeWriter); err != nil {
		return fmt.Errorf("render app.go: %w", err)
	}

	return nil
}

// renderFileByKind renders all the kubernetes objects
// to each file containing all the objects of the same kind
func (j *jamel) renderFileByKind() error {
	kindFileMap, err := j.fileMap()
	if err != nil {
		return fmt.Errorf("filemap: %w", err)
	}

	for _, kind := range orderedKeys(kindFileMap) {
		file := kindFileMap[kind]
		filename := strcase.Kebab(kind) + ".go"
		if j.o.OutputDir != "" {
			filename = filepath.Join(j.o.OutputDir, filename)
		}
		if _, err = j.o.GoCodeWriter.Write([]byte("-- " + filename + " --\n")); err != nil {
			return err
		}
		err = file.Render(j.o.GoCodeWriter)
		if err != nil {
			return fmt.Errorf("render: %w", err)
		}
	}
	return nil
}

func (j *jamel) renderFileByName() error {
	for _, k := range orderedKeys(j.objectsCode) {
		nameVar, stmt := k, j.objectsCode[k]
		objMeta, ok := j.objectsMeta[nameVar]
		if !ok {
			return fmt.Errorf("no object meta for %s", nameVar)
		}
		nameVarObj := j.o.NameVarFunc(objMeta)
		if j.o.RemoveAppName {
			nameVarObj = RemoveAppName(nameVarObj, j.o.AppName)
		}
		file := stmtKubeObjectFile(j.o.OutputPkgName, nameVarObj, stmt)

		filename := j.o.NameFileObjFunc(objMeta)
		if j.o.RemoveAppName {
			filename = RemoveAppName(filename, j.o.AppName)
		}
		if j.o.OutputDir != "" {
			filename = filepath.Join(j.o.OutputDir, filename)
		}

		_, err := j.o.GoCodeWriter.Write([]byte("-- " + filename + " --\n"))
		if err != nil {
			return err
		}
		err = file.Render(j.o.GoCodeWriter)
		if err != nil {
			return fmt.Errorf("render: %w", err)
		}
	}
	return nil
}

func (j *jamel) fileMap() (map[string]*jen.File, error) {
	keepTrack := make(map[string]struct{}, 0)
	kindFileMap := make(map[string]*jen.File, 0)

	// create a file for each kind
	for _, nameVar := range orderedKeys(j.objectsCode) {
		stmt := j.objectsCode[nameVar]
		// for nameVar, stmt := range j.objectsCode {
		objMeta, ok := j.objectsMeta[nameVar]
		if !ok {
			return nil, fmt.Errorf("no object meta for %s", nameVar)
		}
		nameVarObj := j.o.NameVarFunc(objMeta)
		if j.o.RemoveAppName {
			nameVarObj = RemoveAppName(nameVarObj, j.o.AppName)
		}

		// if last letter of nameVar is a number, it is a duplicate
		// we add that number to the nameVarObj
		if lastChar := nameVar[len(nameVar)-1]; lastChar >= '0' && lastChar <= '9' {
			nameVarObj += string(lastChar)
		}

		// check if already a file exists for this kind
		if _, ok := keepTrack[objMeta.Kind]; ok {
			kindFileMap[objMeta.Kind].Line().
				Var().Id(nameVarObj).Op("=").Add(stmt)
			continue
		}
		// no file exists for this kind, create one
		keepTrack[objMeta.Kind] = struct{}{}
		kindFileMap[objMeta.Kind] = stmtKubeObjectFile(
			j.o.OutputPkgName,
			nameVarObj,
			stmt,
		)
	}
	return kindFileMap, nil
}

func (j *jamel) save() error {
	if err := os.MkdirAll(j.o.OutputDir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// convert to go code
	if err := j.generateGo(); err != nil {
		return fmt.Errorf("generate go: %w", err)
	}

	// render all the kubernetes objects
	if j.o.GroupByKind {
		if err := j.saveFileByKind(); err != nil {
			return fmt.Errorf("save by kind: %w", err)
		}
	} else {
		if err := j.saveFileByName(); err != nil {
			return fmt.Errorf("save by name: %w", err)
		}
	}

	// app.go with kubeapp struct
	appPath := filepath.Join(j.o.OutputDir, "app.go")
	appFile := j.appFile()
	if err := appFile.Save(appPath); err != nil {
		return err
	}
	return nil
}

func (j *jamel) saveFileByName() error {
	for _, k := range orderedKeys(j.objectsCode) {
		nameVar, stmt := k, j.objectsCode[k]
		objMeta, ok := j.objectsMeta[nameVar]
		if !ok {
			return fmt.Errorf("no object meta for %s", nameVar)
		}
		nameVarObj := j.o.NameVarFunc(objMeta)
		if j.o.RemoveAppName {
			nameVarObj = RemoveAppName(nameVarObj, j.o.AppName)
		}
		file := stmtKubeObjectFile(j.o.OutputPkgName, nameVarObj, stmt)

		filename := j.o.NameFileObjFunc(objMeta)
		if j.o.RemoveAppName {
			filename = RemoveAppName(filename, j.o.AppName)
		}
		outputPath := filepath.Join(j.o.OutputDir, filename)
		if err := file.Save(outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (j *jamel) saveFileByKind() error {
	kindFileMap, err := j.fileMap()
	if err != nil {
		return fmt.Errorf("filemap: %w", err)
	}
	// write each file contains all the objects of the same kind
	for kind, file := range kindFileMap {
		filename := strcase.Kebab(kind)
		outputPath := filepath.Join(
			j.o.OutputDir,
			filename+".go",
		)
		if err := file.Save(outputPath); err != nil {
			return fmt.Errorf("saving file: %w", err)
		}
	}
	return nil
}

func stmtKubeObjectFile(
	pkgName string,
	nameVarObj string,
	stmt *jen.Statement,
) *jen.File {
	f := jen.NewFile(pkgName)
	f.HeaderComment(headerComment)
	meta.ImportKubernetesPkgAlias(f)
	f.Line()
	f.Var().Id(nameVarObj).Op("=").Add(stmt)
	return f
}

func (j *jamel) appFile() *jen.File {
	appFile := jen.NewFile(j.o.OutputPkgName)
	appFile.HeaderComment(headerComment)
	meta.ImportKubernetesPkgAlias(appFile)
	appFile.Line()

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

	// add Apply function
	appFile.Add(stmtApplyFunc())

	return appFile
}

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

func addMethods(f *jen.File, nameStruct string) *jen.File {
	// Apply
	f.Comment("Apply applies the kubernetes objects to the cluster").Line().
		Func().Params(jen.Id("a").Op("*").Id(nameStruct)).Id("Apply").
		Params(jen.Id("ctx").Qual("context", "Context")).
		Params(jen.Error()).
		BlockFunc(
			func(g *jen.Group) {
				g.Return(
					// jen.Qual(kubeAppPkgPath, "Apply").
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

func orderedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
