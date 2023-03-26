// Copyright (c) Volvo Car AB
// SPDX-License-Identifier: Apache-2.0

package kube

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func returnTypeAlias[X any](
	v reflect.Value,
	typename string,
	x X,
) *jen.Statement {
	if v.Type().String() != typename {
		return jen.Qual(
			v.Type().PkgPath(),
			v.Type().Name(),
		).Call(
			jen.Lit(x),
		)
	}
	return jen.Lit(x)
}

func (j *jamel) convertValue(v reflect.Value) *jen.Statement {
	if v.IsZero() {
		return nil
	}
	switch v.Type().Kind() {
	case reflect.String:
		return returnTypeAlias(v, reflect.String.String(), v.String())
	case reflect.Bool:
		return returnTypeAlias(v, reflect.Bool.String(), v.Bool())
	case reflect.Int:
		return returnTypeAlias(v, reflect.Int.String(), int(v.Int()))
	case reflect.Int64:
		return returnTypeAlias(v, reflect.Int64.String(), v.Int())
	case reflect.Int32:
		return returnTypeAlias(v, reflect.Int32.String(), int32(v.Int()))
	case reflect.Int16:
		return returnTypeAlias(v, reflect.Int16.String(), int16(v.Int()))
	case reflect.Int8:
		return returnTypeAlias(v, reflect.Int8.String(), int8(v.Int()))
	case reflect.Uint:
		return returnTypeAlias(v, reflect.Uint.String(), v.Uint())
	case reflect.Uint64:
		return returnTypeAlias(v, reflect.Uint64.String(), v.Uint())
	case reflect.Uint32:
		return returnTypeAlias(v, reflect.Uint32.String(), uint32(v.Uint()))
	case reflect.Uint16:
		return returnTypeAlias(v, reflect.Uint16.String(), uint16(v.Uint()))
	case reflect.Uint8:
		return returnTypeAlias(v, reflect.Uint8.String(), uint8(v.Uint()))
	case reflect.Float32:
		return returnTypeAlias(v, reflect.Float32.String(), float32(v.Float()))
	case reflect.Float64:
		return returnTypeAlias(v, reflect.Float64.String(), v.Float())
	// ----------------------------------------
	//
	// Map types
	//
	case reflect.Map:
		pk := prefixKind(v)
		vf := jen.DictFunc(
			func(d jen.Dict) {
				for _, key := range v.MapKeys() {
					k := j.convertValue(key)
					d[k] = j.convertValue(v.MapIndex(key))
				}
			},
		)
		return pk.Values(vf)
	// ----------------------------------------
	//
	// Array and Slice types
	//
	case reflect.Array, reflect.Slice:
		pk := prefixKind(v)
		if v.Len() == 0 {
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
	// ----------------------------------------
	//
	// Struct types
	//
	case reflect.Struct:
		name := v.Type().Name()
		switch name {
		case "Quantity":
			return convertQuantity(v)
		// case "IntOrString":
		// 	return convertIntOrString(v)
		case "Secret":
			return j.convertSecret(v).
				Comment("TODO: SECRETS SHOULD BE STORED ELSEWHERE THAN IN THE CODE!!!!")
		}

		pk := prefixKind(v)
		vf := jen.DictFunc(
			func(d jen.Dict) {
				for i := 0; i < v.NumField(); i++ {
					d[jen.Id(v.Type().Field(i).Name)] = j.convertValue(v.Field(i))
				}
			},
		)

		return pk.Values(vf)
	// ----------------------------------------
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
		case reflect.Int:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Int.String(), v.Elem().Interface().(int),
				),
			)
		case reflect.Int64:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Int64.String(), v.Elem().Int(),
				),
			)
		case reflect.Int32:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Int32.String(), v.Elem().Interface().(int32),
				),
			)
		case reflect.Int16:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Int16.String(), v.Elem().Interface().(int16),
				),
			)
		case reflect.Int8:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Int8.String(), v.Elem().Interface().(int8),
				),
			)
		case reflect.Uint:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Uint.String(), v.Elem().Interface().(uint),
				),
			)
		case reflect.Uint64:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Uint64.String(), v.Elem().Interface().(uint64),
				),
			)
		case reflect.Uint32:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Uint32.String(), v.Elem().Interface().(uint32),
				),
			)
		case reflect.Uint16:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Uint16.String(), v.Elem().Interface().(uint16),
				),
			)
		case reflect.Uint8:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Uint8.String(), v.Elem().Interface().(uint8),
				),
			)
		case reflect.Float32:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Float32.String(), v.Elem().Interface().(float32),
				),
			)
		case reflect.Float64:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Float64.String(), v.Elem().Interface().(float64),
				),
			)
		case reflect.Bool:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.Bool.String(), v.Elem().Bool(),
				),
			)
		case reflect.String:
			return jen.Id("P").Call(
				returnTypeAlias(
					v.Elem(),
					reflect.String.String(),
					v.Elem().String(),
				),
			)
		default:
			return jen.Op("&").Add(j.convertValue(v.Elem()))
		}

	default:
		slog.Info("unsupported", slog.String("kind", v.Type().Kind().String()))
		return jen.Nil()
	}
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

	pk := prefixKind(v)
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
	pk := prefixKind(v)
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

// convertQuantity converts a Quantity to a jen statement.
// A Quantity is a struct with unexported fields.
// To fill it, we must use the resource.MustParse function.
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

// prefixKind returns the jen statement for the type of the value.
// It is used to set the proper import package for the type.
// For instance, `v1.ServiceAccount`, renamed `corev1.ServiceAccount`
// with `import corev1 "k8s.io/api/core/v1"`
func prefixKind(v reflect.Value) *jen.Statement {
	if v.IsZero() {
		return nil
	}
	switch v.Kind() {
	case reflect.String:
		return jen.String()
	case reflect.Bool:
		return jen.Bool()
	case reflect.Int:
		return jen.Int()
	case reflect.Int64:
		return jen.Int64()
	case reflect.Int32:
		return jen.Int32()
	case reflect.Int16:
		return jen.Int16()
	case reflect.Int8:
		return jen.Int8()
	case reflect.Uint:
		return jen.Uint()
	case reflect.Uint64:
		return jen.Uint64()
	case reflect.Uint32:
		return jen.Uint32()
	case reflect.Uint16:
		return jen.Uint16()
	case reflect.Uint8:
		return jen.Uint8()
	case reflect.Float32:
		return jen.Float32()
	case reflect.Float64:
		return jen.Float64()
		// ----------------------------------------
	case reflect.Ptr:
		if v.IsNil() {
			return jen.Nil()
		}

		return jen.Op("&").Add(prefixKind(v.Elem()))
		// ----------------------------------------
	case reflect.Array, reflect.Slice:
		pkgPath := v.Type().Elem().PkgPath()
		name := v.Type().Elem().Name()
		if pkgPath == "" {
			// built-in type
			return jen.Index().Id(name)
		}
		return jen.Index().Qual(pkgPath, name)
		// ----------------------------------------
	case reflect.Map:
		// Resolve Key type first
		kname := v.Type().Key().Name()
		keyPkgPath := v.Type().Key().PkgPath()
		keyName := jen.Qual(keyPkgPath, kname)
		if keyPkgPath == "" {
			// built-in type
			keyName = jen.Id(kname)
		}

		// Resolve Value type
		vname := v.Type().Elem().Name()
		if vname == "" {
			// If the value type is not a named type, use string as the default type.
			// In some configMap in YAML, `data: {}` is used to represent an empty map
			// which yields `map[string]{}` and cause the rendering to fail.
			vname = "string"
		}
		valuePkgPath := v.Type().Elem().PkgPath()
		valueName := jen.Qual(valuePkgPath, vname)
		if valuePkgPath == "" {
			// built-in type
			valueName = jen.Id(vname)
		}

		return jen.Map(keyName).Add(valueName)
		// ----------------------------------------
	case reflect.Struct:
		pkgPath := v.Type().PkgPath()
		name := v.Type().Name()
		return jen.Qual(pkgPath, name)
		// ----------------------------------------
	default:
		slog.Info(
			"unknown kind",
			slog.String("v.Kind()", v.Kind().String()),
			slog.String("v.Type().PkgPath()", v.Type().PkgPath()),
		)
		return jen.Empty()
	}
}
