// Copyright 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var commentSecret = "TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!" //nolint:gosec

func returnTypeAlias(
	v reflect.Value,
	typename string,
	stmt *jen.Statement,
) *jen.Statement {
	if v.Type().String() != typename {
		return jen.Qual(
			v.Type().PkgPath(),
			v.Type().Name(),
		).Call(
			stmt,
		)
	}
	return stmt
}

func (j *jamel) convertValue(v reflect.Value) *jen.Statement {
	if v.IsZero() {
		return nil
	}
	switch v.Type().Kind() {
	case reflect.String:
		s := v.String()
		if strings.Contains(s, "\n") {
			return returnTypeAlias(
				v, reflect.String.String(), rawString(s),
			)
		}
		return returnTypeAlias(
			v, reflect.String.String(), jen.Lit(v.String()),
		)
	case reflect.Bool:
		return returnTypeAlias(
			v,
			reflect.Bool.String(),
			jen.Lit(v.Bool()),
		)
	case reflect.Int:
		return returnTypeAlias(
			v,
			reflect.Int.String(),
			jen.Lit(int(v.Int())),
		)
	case reflect.Int64:
		return returnTypeAlias(
			v,
			reflect.Int64.String(),
			jen.Lit(v.Int()),
		)
	case reflect.Int32:
		return returnTypeAlias(
			v,
			reflect.Int32.String(),
			jen.Lit(int32(v.Int())),
		)
	case reflect.Int16:
		return returnTypeAlias(
			v,
			reflect.Int16.String(),
			jen.Lit(int16(v.Int())),
		)
	case reflect.Int8:
		return returnTypeAlias(
			v,
			reflect.Int8.String(),
			jen.Lit(int8(v.Int())),
		)
	case reflect.Uint:
		return returnTypeAlias(
			v,
			reflect.Uint.String(),
			jen.Lit(v.Uint()),
		)
	case reflect.Uint64:
		return returnTypeAlias(
			v,
			reflect.Uint64.String(),
			jen.Lit(v.Uint()),
		)
	case reflect.Uint32:
		return returnTypeAlias(
			v,
			reflect.Uint32.String(),
			jen.Lit(uint32(v.Uint())),
		)
	case reflect.Uint16:
		return returnTypeAlias(
			v,
			reflect.Uint16.String(),
			jen.Lit(uint16(v.Uint())),
		)
	case reflect.Uint8:
		return returnTypeAlias(
			v,
			reflect.Uint8.String(),
			jen.Lit(uint8(v.Uint())),
		)
	case reflect.Float32:
		return returnTypeAlias(
			v,
			reflect.Float32.String(),
			jen.Lit(float32(v.Float())),
		)
	case reflect.Float64:
		return returnTypeAlias(v, reflect.Float64.String(), jen.Lit(v.Float()))
	//
	// Map types
	//
	case reflect.Map:
		pk := j.prefixKind(v)
		vf := jen.DictFunc(
			func(d jen.Dict) {
				for _, key := range v.MapKeys() {
					k := j.convertValue(key)
					d[k] = j.convertValue(v.MapIndex(key))
				}
			},
		)
		return pk.Values(vf)
	//
	// Array and Slice types
	//
	case reflect.Array, reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			r := []rune{}
			var i int
			for i = 0; i < v.Len(); i++ {
				bla := v.Index(i).Interface().(uint8)
				if strconv.IsPrint(rune(bla)) {
					r = append(r, rune(bla))
				}
			}
			if i == v.Len() {
				s := string(r)
				return jen.Index().Byte().Params(jen.Lit(s))
			}
		}

		pk := j.prefixKind(v)
		if isEmptyValue(v) {
			return pk.Block()
		}

		vf := pk.ValuesFunc(
			func(g *jen.Group) {
				for i := 0; i < v.Len(); i++ {
					g.Add(j.convertValue(v.Index(i)))
				}
			},
		)

		return vf
	//
	// Struct types
	//
	case reflect.Struct:
		switch v.Type().Name() {
		case "Quantity":
			return convertQuantity(v)
		case "Secret":
			return j.convertSecret(v).Comment(commentSecret)
		}

		pk := j.prefixKind(v)
		vf := jen.DictFunc(
			func(d jen.Dict) {
				// code from [json.Encode](https://go.dev/src/encoding/json/encode.go)
				// func typeFields(t reflect.Type) structFields
				for i := 0; i < v.NumField(); i++ {
					vtf := v.Type().Field(i)
					if vtf.Anonymous {
						if !vtf.IsExported() && v.Type().Kind() != reflect.Struct {
							// Ignore embedded fields of unexported non-struct types.
							continue
						}
					} else if !vtf.IsExported() {
						// Ignore unexported fields.
						continue
					}

					d[jen.Id(vtf.Name)] = j.convertValue(v.Field(i))
				}
			},
		)

		return pk.Values(vf)
	//
	// Pointer types
	//
	case reflect.Ptr:
		if v.Elem().IsZero() {
			return nil
		}

		//
		// we cannot use a point of basic types
		// therefore we need to use a function to return the pointer
		// ❌: &int
		// ✅: func P[T any](t T) *T { return &t }
		//
		switch v.Elem().Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16,
			reflect.Int8, reflect.Uint, reflect.Uint64, reflect.Uint32,
			reflect.Uint16, reflect.Uint8, reflect.Float32, reflect.Float64,
			reflect.Bool, reflect.String:
			return jen.Id("P").Call(j.convertValue(v.Elem()))
		default:
			return jen.Op("&").Add(j.convertValue(v.Elem()))
		}

	case reflect.Interface:
		if isEmptyValue(v) {
			return jen.Nil()
		}
		return jen.Interface(j.convertValue(v.Elem()))

	default:
		slog.Error(
			"unsupported",
			slog.String("type", v.String()),
			slog.String("kind", v.Kind().String()),
		)
		return jen.Nil()
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func (j *jamel) configMapComment(
	obj *corev1.ConfigMap,
	data []byte,
) (*jen.Statement, error) {
	// parse YAML to extract comments
	var d yaml.Node
	err := yaml.Unmarshal(data, &d)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	found := false
	mc := make(map[string]string)

outer:
	for _, c := range d.Content {
		if c.Kind == yaml.MappingNode && len(c.Content) > 0 {
			for _, c2 := range c.Content {
				// the data value is right after the "data" key
				if found && c2.Kind == yaml.MappingNode {
					for _, cdata := range c2.Content {
						if cdata.HeadComment != "" {
							// remove the "# " from the comment
							mc[cdata.Value] = strings.ReplaceAll(
								cdata.HeadComment,
								"# ",
								"",
							)
						}
					}
					// we have the comments, break out
					break outer
				}
				// the data value is right after the "data" key
				if c2.Kind == yaml.ScalarNode && c2.Value == "data" {
					found = true
				}
			}
		}
	}
	rv := reflect.ValueOf(obj)
	return j.convertConfigMap(rv, mc), nil
}

func (j *jamel) convertConfigMap(
	v reflect.Value,
	comment map[string]string,
) *jen.Statement {
	if v.IsZero() {
		return jen.Nil()
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	pk := j.prefixKind(v)
	vf := jen.DictFunc(
		func(d jen.Dict) {
			for i := 0; i < v.NumField(); i++ {
				switch v.Type().Field(i).Name {
				case "Data":
					d[jen.Id(v.Type().Field(i).Name)] = j.convertConfigMapData(
						v.Field(i),
						comment,
					)
				default:
					d[jen.Id(v.Type().Field(i).Name)] = j.convertValue(v.Field(i))
				}
			}
		},
	)
	return jen.Op("&").Add(pk.Values(vf))
}

func (j *jamel) convertConfigMapData(
	field reflect.Value,
	comment map[string]string,
) *jen.Statement {
	if field.IsZero() {
		return jen.Nil()
	}

	keys := field.MapKeys()
	sort.SliceStable(
		keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		},
	)

	return jen.Map(jen.String()).String().ValuesFunc(
		func(v *jen.Group) {
			for _, k := range keys {
				v.Add(
					jen.CustomFunc(
						jen.Options{
							Open:      "",
							Close:     "",
							Separator: ",",
							Multi:     true,
						}, func(g *jen.Group) {
							c, ok := comment[k.String()]
							if ok {
								g.Add(
									jen.Comment(c).Line(),
									jen.Lit(k.String()),
									jen.Op(":"),
									rawString(field.MapIndex(k).String()),
								)
							} else {
								g.Add(
									jen.Lit(k.String()),
									jen.Op(":"),
									rawString(field.MapIndex(k).String()),
								)
							}
						},
					),
				)
			}
		},
	)
}

func rawString(s string) *jen.Statement {
	if len(s) == 0 {
		return jen.Lit("")
	}
	newS := s
	if !strings.Contains(newS, "\n") { // single line string
		return jen.Lit(newS)
	}
	if strings.Contains(newS, "`") {
		newS = strings.ReplaceAll(newS, "`", "\"")
	}
	return jen.Custom(
		jen.Options{Open: "`", Close: "`", Multi: true},
		jen.Op(newS),
	)
}

// convertSecret converts a Secret to a jen statement.
func (j *jamel) convertSecret(v reflect.Value) *jen.Statement {
	if v.IsZero() {
		return jen.Nil()
	}
	pk := j.prefixKind(v)
	vf := jen.DictFunc(
		func(d jen.Dict) {
			for i := 0; i < v.NumField(); i++ {
				switch v.Type().Field(i).Name {
				case "Data":
					d[jen.Id(v.Type().Field(i).Name)] = j.replaceSecretData(v.Field(i))
				// TODO: TLS cert ?
				default:
					d[jen.Id(v.Type().Field(i).Name)] = j.convertValue(v.Field(i))
				}
			}
		},
	)

	return pk.Values(vf)
}

func (j *jamel) replaceSecretData(field reflect.Value) jen.Code {
	if field.IsZero() {
		return jen.Nil()
	}
	return jen.Map(jen.String()).Index().Byte().ValuesFunc(
		func(v *jen.Group) {
			for _, k := range field.MapKeys() {
				v.Add(
					jen.DictFunc(
						func(d jen.Dict) {
							if j.o.RedactSecrets {
								d[jen.Lit(k.String())] = jen.Index().Byte().Call(jen.Lit("<REDACTED>"))
							} else {
								d[jen.Lit(k.String())] = jen.Index().Byte().Call(jen.Lit(string(field.MapIndex(k).Bytes())))
							}
						},
					),
				)
			}
		},
	)
}

// convertQuantity converts a [resource.Quantity] to a jen statement.
// A Quantity is a struct with unexported fields.
// To fill it, we must use the [resource.MustParse] function.
// Refer to kubernetes code base for more details.
func convertQuantity(field reflect.Value) *jen.Statement {
	if field.IsZero() {
		return jen.Empty()
	}

	qty := ""

	q, ok := field.Interface().(resource.Quantity)
	if ok {
		qty = q.String()
	} else {
		qty = ">>> FAILED TO PARSE QUANTITY <<<"
	}

	return jen.Qual(
		"k8s.io/apimachinery/pkg/api/resource",
		"MustParse",
	).Call(jen.Lit(qty))
}

var repl = strings.NewReplacer(
	"-",
	"",
	".",
	"",
	"pkg/api/",
	"",
	"pkg/apis/",
	"",
)

// storePkgPath stores the first non-native kubernetes package it finds,
// hoping it will be the CRD Go package instead of the APIVersion URL.
func (j *jamel) storePkgPath(t reflect.Type) (string, string) {
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return "", t.Name()
	}
	name := t.Name()
	split := strings.Split(pkgPath, "/")
	if split[0] == "k8s.io" {
		return pkgPath, name
	}
	if len(split) >= 2 && j.crdCurrent == "" {
		j.crdCurrent = pkgPath
		// package alias
		switch split[0] {
		case "sigs.k8s.io":
			j.crdPkgAlias[pkgPath] = repl.Replace(strings.Join(split[1:], ""))
		case "github.com":
			j.crdPkgAlias[pkgPath] = repl.Replace(strings.Join(split[1:], ""))
		default:
			j.crdPkgAlias[pkgPath] = strings.Join(split[len(split)-2:], "")
		}

	}
	return pkgPath, name
}

// prefixKind returns the jen statement for the type of the value.
// It is used to set the proper import package for the type.
// For instance, `v1.ServiceAccount`, renamed [corev1.ServiceAccount]
// with `import corev1 "k8s.io/api/core/v1"`
func (j *jamel) prefixKind(v reflect.Value) *jen.Statement {
	if v.IsZero() {
		return nil
	}
	switch v.Kind() {

	case reflect.Ptr:
		if v.IsNil() {
			return jen.Nil()
		}

		return jen.Op("&").Add(j.prefixKind(v.Elem()))

	case reflect.Array, reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.Ptr:
			// arrays of pointers
			pkgPath, name := j.storePkgPath(v.Type().Elem().Elem())
			return jen.Index().Op("*").Add(jen.Qual(pkgPath, name))

		case reflect.Slice, reflect.Array:
			// arrays of arrays
			pkgPath, name := j.storePkgPath(v.Type().Elem().Elem())
			if pkgPath == "" {
				// built-in type

				// []uint8 is the representation of []byte
				if name == "uint8" {
					name = "byte"
				}
				return jen.Index().Index().Id(name)
			}
			return jen.Index().Index().Qual(pkgPath, name)

		// array of maps .... TODO

		case reflect.Interface:
			return jen.Index().Interface()
		default:
			pkgPath, name := j.storePkgPath(v.Type().Elem())
			if pkgPath == "" {
				// built-in type
				return jen.Index().Id(name)
			}
			return jen.Index().Qual(pkgPath, name)
		}

	case reflect.Map:
		// Resolve Key type first
		kPkgPath, kname := j.storePkgPath(v.Type().Key())
		keyName := jen.Qual(kPkgPath, kname)
		if kPkgPath == "" {
			// built-in type
			keyName = jen.Id(kname)
		}

		switch v.Type().Elem().Kind() {
		case reflect.Ptr:
			// map of pointers
			pkgPath, name := j.storePkgPath(v.Type().Elem().Elem())
			return jen.Map(keyName).Op("*").Add(jen.Qual(pkgPath, name))

		case reflect.Slice, reflect.Array:
			// map of array
			pkgPath, name := j.storePkgPath(v.Type().Elem().Elem())
			if pkgPath == "" {
				// built-in type
				return jen.Map(keyName).Index().Id(name)
			}
			return jen.Map(keyName).Index().Qual(pkgPath, name)

		case reflect.Map:
			// map of map
			ekPkgPath, ekname := j.storePkgPath(v.Type().Elem().Key())
			ekeyName := jen.Qual(ekPkgPath, ekname)
			if ekPkgPath == "" {
				// built-in type
				ekeyName = jen.Id(kname)
			}
			return jen.Map(keyName).Map(ekeyName).Add(
				jen.Qual(
					ekPkgPath,
					ekname,
				),
			)
		case reflect.Interface:
			return jen.Map(keyName).Interface()
		default:

			// Resolve Value type
			valuePkgPath, vname := j.storePkgPath(v.Type().Elem())
			if vname == "" {
				// If the value type is not a named type, use string as the default type.
				// In some configMap in YAML, `data: {}` is used to represent an empty map
				// which yields `map[string]{}` and cause the rendering to fail.
				vname = "string"
			}
			valueName := jen.Qual(valuePkgPath, vname)
			if valuePkgPath == "" {
				// built-in type
				valueName = jen.Id(vname)
			}

			return jen.Map(keyName).Add(valueName)
		}

	case reflect.Struct:
		pkgPath, name := j.storePkgPath(v.Type())
		return jen.Qual(pkgPath, name)

	default:
		slog.Info(
			"unknown kind",
			slog.String("v.Kind()", v.Kind().String()),
			slog.String("v.Type().PkgPath()", v.Type().PkgPath()),
		)
		return jen.Empty()
	}
}
