package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func generateTemplComponent(comp ComponentConfig, structInfo StructInfo, config *Config) error {
	tmplFile := filepath.Join("templates", comp.Template)

	fmt.Printf("   Reading template: %s\n", tmplFile)

	tmplContent, err := os.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("template file not found: %s - %w", tmplFile, err)
	}

	fmt.Printf("   Template size: %d bytes\n", len(tmplContent))

	t, err := template.New("component").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	fmt.Printf("   Template parsed successfully\n")

	outputPath := filepath.Join(config.Global.Output.TemplDir, comp.Name+".templ")
	fmt.Printf("   Creating output file: %s\n", outputPath)

	if err := os.MkdirAll(config.Global.Output.TemplDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	fmt.Printf("   File created, executing template...\n")

	data := struct {
		Component     ComponentConfig
		Struct        StructInfo
		StructsImport string
	}{
		Component:     comp,
		Struct:        structInfo,
		StructsImport: config.Global.Source.StructsImport,
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("   ✅ Generated: %s\n", outputPath)

	return nil
}
