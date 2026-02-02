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
}

func generateLuaFile(snippets []SnippetData, componentName string, config *Config) error {
	tmplStr := `local ls = require("luasnip")
local s = ls.snippet
local t = ls.text_node

return {
{{range .}}
  s("{{.Prefix}}", {
    t("@components.{{.ComponentName}}(&structs.{{.StructName}}{})"),
  }),
{{end}}
}
`

	t := template.Must(template.New("snippets").Parse(tmplStr))

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

func generateSnippets(comp ComponentConfig, config *Config) error {
	snippets := []SnippetData{}

	for _, variant := range comp.Variants {
		snippet := SnippetData{
			Prefix:        variant.Prefix,
			ComponentName: variant.Name + comp.Name,
			StructName:    comp.Struct,
		}
		snippets = append(snippets, snippet)
	}

	return generateLuaFile(snippets, comp.Name, config)
}

func generateStructSnippets(structs map[string]StructInfo, config *Config) error {
	structsToGenerate := []string{"Common", "Hx", "FormBehaviors"}

	type StructSnippet struct {
		Prefix     string
		StructName string
		Fields     []string
	}

	snippets := []StructSnippet{}

	for _, structName := range structsToGenerate {
		structInfo, ok := structs[structName]
		if !ok {
			continue
		}

		fieldLines := []string{}
		for _, field := range structInfo.Fields {
			if field.Name == "" {
				continue
			}
			defaultVal := getDefaultValue(field)
			fieldLines = append(fieldLines, fmt.Sprintf("\t%s: %s,", field.Name, defaultVal))
		}

		snippets = append(snippets, StructSnippet{
			Prefix:     "str" + strings.ToLower(structName),
			StructName: structName,
			Fields:     fieldLines,
		})
	}

	if len(snippets) == 0 {
		return nil
	}

	tmplStr := `local ls = require("luasnip")
local s = ls.snippet
local t = ls.text_node

return {
{{range .}}
  s("{{.Prefix}}", {
    t({
      [[structs.{{.StructName}}{]],
{{- range .Fields}}
      [[{{.}}]],
{{- end}}
      [[}]],
    }),
  }),
{{end}}
}
`

	tmpl := template.Must(template.New("structs").Parse(tmplStr))

	outputPath := filepath.Join(config.Global.Output.SnippetsDir, "structs.lua")

	if err := os.MkdirAll(config.Global.Output.SnippetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create snippets directory: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create structs.lua: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, snippets)
}

func getDefaultValue(field FieldInfo) string {
	switch field.Type {
	case "string":
		return `""`
	case "bool":
		return "false"
	case "int", "int64", "int32":
		return "0"
	case "float64", "float32":
		return "0.0"
	default:
		return `""`
	}
}
