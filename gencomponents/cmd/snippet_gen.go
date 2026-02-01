package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type SnippetData struct {
	Prefix        string
	ComponentName string
	StructName    string
	Fields        []SnippetField
}

type SnippetField struct {
	Name         string
	DefaultValue string
	TabOrder     int
}

func generateSnippets(comp ComponentConfig, structInfo StructInfo, config *Config) error {
	snippets := []SnippetData{}

	for _, variant := range comp.Variants {
		snippet := SnippetData{
			Prefix:        variant.Prefix,
			ComponentName: variant.Name + comp.Name,
			StructName:    comp.Struct,
			Fields:        buildFieldList(structInfo, config),
		}
		snippets = append(snippets, snippet)
	}

	return generateLuaFile(snippets, comp.Name, config)
}

func buildFieldList(structInfo StructInfo, config *Config) []SnippetField {
	fields := []SnippetField{}
	tabOrder := 1

	for _, field := range structInfo.Fields {
		defaultVal := getDefaultValue(field, config)
		fields = append(fields, SnippetField{
			Name:         field.Name,
			DefaultValue: defaultVal,
			TabOrder:     tabOrder,
		})
		tabOrder++
	}

	return fields
}

func getDefaultValue(field FieldInfo, config *Config) string {
	// Check field-specific defaults first
	if val, ok := config.Global.FieldDefaults[field.Name]; ok {
		return val
	}

	// Check type defaults
	if val, ok := config.Global.SnippetDefaults[field.Type]; ok {
		return val
	}

	// Fallback defaults - ALWAYS return quoted strings
	switch field.Type {
	case "string":
		return `""`
	case "bool":
		return `"false"`
	case "int", "int64", "int32":
		return `"0"`
	case "float64", "float32":
		return `"0.0"`
	default:
		return `""`
	}
}

func generateLuaFile(snippets []SnippetData, componentName string, config *Config) error {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"openBrace":  func() string { return "{{" },
		"closeBrace": func() string { return "}}" },
	}

	tmplStr := `local ls = require("luasnip")
local s = ls.snippet
local i = ls.insert_node
local fmt = require("luasnip.extras.fmt").fmt

return {
{{range .}}
  s("{{.Prefix}}", fmt([[
@{{.ComponentName}}(&structs.{{.StructName}}{{openBrace}}
{{range .Fields}}	{{.Name}}: {},
{{end}}{{closeBrace}})
  ]], {
{{range .Fields}}    i({{.TabOrder}}, {{.DefaultValue}}),
{{end}}  })),

{{end}}
}
`

	t := template.Must(template.New("snippets").Funcs(funcMap).Parse(tmplStr))

	outputPath := filepath.Join(config.Global.Output.SnippetsDir, strings.ToLower(componentName)+".lua")

	if err := os.MkdirAll(config.Global.Output.SnippetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create snippets directory: %w", err)
	}

	fmt.Printf("   Creating snippet: %s\n", outputPath)

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer f.Close()

	if err := t.Execute(f, snippets); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("   ✅ Wrote snippet: %s\n", outputPath)

	return nil
}
