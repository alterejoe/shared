package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type FieldInfo struct {
	Name      string
	Type      string
	IsPointer bool
	IsStruct  bool
	Package   string
}

func ParseStructs(structsPath string) (map[string]StructInfo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, structsPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structs := make(map[string]StructInfo)

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok {
					return true
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return true
				}

				structInfo := StructInfo{
					Name:   typeSpec.Name.Name,
					Fields: []FieldInfo{},
				}

				for _, field := range structType.Fields.List {
					for _, name := range field.Names {
						fieldInfo := parseField(name.Name, field.Type)
						structInfo.Fields = append(structInfo.Fields, fieldInfo)
					}

					// Handle embedded structs (no name)
					if len(field.Names) == 0 {
						fieldInfo := parseField("", field.Type)
						if fieldInfo.IsStruct {
							// This is an embedded struct, we'll expand its fields
							if embeddedStruct, exists := structs[fieldInfo.Type]; exists {
								structInfo.Fields = append(structInfo.Fields, embeddedStruct.Fields...)
							}
						}
					}
				}

				structs[structInfo.Name] = structInfo
				return true
			})
		}
	}

	return structs, nil
}

func parseField(name string, fieldType ast.Expr) FieldInfo {
	field := FieldInfo{Name: name}

	switch t := fieldType.(type) {
	case *ast.Ident:
		field.Type = t.Name
		field.IsStruct = isStructType(t.Name)

	case *ast.StarExpr:
		field.IsPointer = true
		if ident, ok := t.X.(*ast.Ident); ok {
			field.Type = ident.Name
			field.IsStruct = isStructType(ident.Name)
		}

	case *ast.SelectorExpr:
		if pkg, ok := t.X.(*ast.Ident); ok {
			field.Package = pkg.Name
			field.Type = t.Sel.Name
			field.IsStruct = true
		}
	}

	return field
}

func isStructType(typeName string) bool {
	// Basic Go types
	basicTypes := map[string]bool{
		"string": true, "int": true, "int64": true, "int32": true,
		"bool": true, "float64": true, "float32": true,
	}
	return !basicTypes[typeName]
}
