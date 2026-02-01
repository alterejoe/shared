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
	StructType   string
	DefaultValue string
	TabOrder     int
	IsEmbedded   bool
}

// Global tracker for struct snippets across all components
var globalStructSnippets = make(map[string]SnippetData)

func generateSnippets(comp ComponentConfig, structInfo StructInfo, config *Config) error {
	snippets := []SnippetData{}

	// Generate component snippets
	for _, variant := range comp.Variants {
		snippet := SnippetData{
			Prefix:        variant.Prefix,
			ComponentName: variant.Name + comp.Name,
			StructName:    comp.Struct,
			Fields:        buildFieldListTopLevel(structInfo, config),
		}
		snippets = append(snippets, snippet)
	}

	// Collect struct snippets globally (don't add to component file)
	collectStructSnippets(structInfo, config)

	return generateLuaFile(snippets, comp.Name, config)
}
func buildFieldListTopLevel(structInfo StructInfo, config *Config) []SnippetField {
	fields := []SnippetField{}
	tabOrder := 1

	for _, field := range structInfo.Fields {
		// Handle embedded structs (fields with no name)
		if field.Name == "" && field.IsStruct {
			// This is an embedded struct - add it as a placeholder
			fields = append(fields, SnippetField{
				Name:         field.Type,
				StructType:   field.Type,
				DefaultValue: `""`,
				TabOrder:     tabOrder,
				IsEmbedded:   true,
			})
			tabOrder++
		} else if field.Name != "" {
			// This is a direct field on the struct
			defaultVal := getDefaultValue(field, config)
			fields = append(fields, SnippetField{
				Name:         field.Name,
				DefaultValue: defaultVal,
				TabOrder:     tabOrder,
				IsEmbedded:   false,
			})
			tabOrder++
		}
	}

	return fields
}
func generateStructSnippet(structName string, structInfo StructInfo, config *Config) {
	fields := []SnippetField{}

	for i, field := range structInfo.Fields {
		if field.Name == "" {
			continue // Skip embedded structs within these
		}

		defaultVal := getDefaultValue(field, config)
		fields = append(fields, SnippetField{
			Name:         field.Name,
			DefaultValue: defaultVal,
			TabOrder:     i + 1,
		})
	}

	prefix := "str" + strings.ToLower(structName)
	globalStructSnippets[prefix] = SnippetData{
		Prefix:        prefix,
		ComponentName: "",
		StructName:    structName,
		Fields:        fields,
	}
}
func collectStructSnippets(structInfo StructInfo, config *Config) {
	// Collect fields for each embedded struct
	commonFields := []SnippetField{}
	hxFields := []SnippetField{}
	formBehaviorFields := []SnippetField{}

	for _, field := range structInfo.Fields {
		defaultVal := getDefaultValue(field, config)
		snippetField := SnippetField{
			Name:         field.Name,
			DefaultValue: defaultVal,
		}

		switch getEmbeddedStructType(field.Name) {
		case "Common":
			snippetField.TabOrder = len(commonFields) + 1
			commonFields = append(commonFields, snippetField)
		case "Hx":
			snippetField.TabOrder = len(hxFields) + 1
			hxFields = append(hxFields, snippetField)
		case "FormBehaviors":
			snippetField.TabOrder = len(formBehaviorFields) + 1
			formBehaviorFields = append(formBehaviorFields, snippetField)
		}
	}

	// Add to global tracker (won't duplicate)
	if len(commonFields) > 0 {
		globalStructSnippets["strcommon"] = SnippetData{
			Prefix:        "strcommon",
			ComponentName: "",
			StructName:    "Common",
			Fields:        commonFields,
		}
	}

	if len(hxFields) > 0 {
		globalStructSnippets["strhx"] = SnippetData{
			Prefix:        "strhx",
			ComponentName: "",
			StructName:    "Hx",
			Fields:        hxFields,
		}
	}

	if len(formBehaviorFields) > 0 {
		globalStructSnippets["strformbehaviors"] = SnippetData{
			Prefix:        "strformbehaviors",
			ComponentName: "",
			StructName:    "FormBehaviors",
			Fields:        formBehaviorFields,
		}
	}
}

func generateStructSnippetsFile(config *Config) error {
	if len(globalStructSnippets) == 0 {
		return nil
	}

	// Convert map to slice
	snippets := []SnippetData{}
	for _, snippet := range globalStructSnippets {
		snippets = append(snippets, snippet)
	}

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
structs.{{.StructName}}{{openBrace}}
{{range .Fields}}	{{.Name}}: {},
{{end}}{{closeBrace}}
  ]], {
{{range .Fields}}    i({{.TabOrder}}, {{.DefaultValue}}),
{{end}}  })),
{{end}}
}
`

	t := template.Must(template.New("structs").Funcs(funcMap).Parse(tmplStr))

	outputPath := filepath.Join(config.Global.Output.SnippetsDir, "structs.lua")

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create structs.lua: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, snippets); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func getEmbeddedStructType(fieldName string) string {
	commonFields := map[string]bool{
		"ID": true, "Name": true, "Value": true, "Disabled": true,
	}

	hxFields := map[string]bool{
		"Method": true, "URL": true, "Target": true, "Include": true,
		"Trigger": true, "Swap": true, "Indicator": true, "Vals": true,
		"Confirm": true, "PushURL": true, "Boost": true,
	}

	formBehaviorFields := map[string]bool{
		"ConstraintForm": true, "DirtyWatch": true, "DirtyGroup": true,
		"EnableOnValid": true, "EnableTarget": true,
	}

	if commonFields[fieldName] {
		return "Common"
	}
	if hxFields[fieldName] {
		return "Hx"
	}
	if formBehaviorFields[fieldName] {
		return "FormBehaviors"
	}

	return ""
}

func getDefaultValue(field FieldInfo, config *Config) string {
	if val, ok := config.Global.FieldDefaults[field.Name]; ok {
		return val
	}

	if val, ok := config.Global.SnippetDefaults[field.Type]; ok {
		return val
	}

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
@components.{{.ComponentName}}(&structs.{{.StructName}}{{openBrace}}
{{range .Fields}}{{if .IsEmbedded}}	{{.Name}}: {},
{{else}}	{{.Name}}: {},
{{end}}{{end}}{{closeBrace}})
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

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer f.Close()

	if err := t.Execute(f, snippets); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
