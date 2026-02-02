package main

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/alterejoe/shared/create"
)

func main() {
	logger := create.CreateLogger()

	// Load all configs from config directory
	config, err := LoadConfig("config")
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		return
	}

	logger.Info("Loaded components", "count", len(config.Components))
	logger.Info("Structs path", "path", config.Global.Source.StructsPath)
	logger.Info("Structs import", "import", config.Global.Source.StructsImport)
	logger.Info("Output templ dir", "dir", config.Global.Output.TemplDir)
	logger.Info("Output snippets dir", "dir", config.Global.Output.SnippetsDir)

	// List the components
	for i, comp := range config.Components {
		logger.Info("Component", "number", i+1, "name", comp.Name, "struct", comp.Struct, "template", comp.Template)
	}

	// Parse structs from configured path
	logger.Info("Parsing structs", "path", config.Global.Source.StructsPath)
	structs, err := ParseStructs(config.Global.Source.StructsPath)
	if err != nil {
		logger.Error("Failed to parse structs", "error", err)
		return
	}
	logger.Info("Found structs", "count", len(structs))

	// Expand embedded structs for templ generation
	expandedStructs := expandEmbeddedStructs(structs)

	// Generate components
	for _, comp := range config.Components {
		logger.Info("Processing component", "name", comp.Name)

		expandedStructInfo, ok := expandedStructs[comp.Struct]
		if !ok {
			logger.Warn("Skipping component - struct not found", "component", comp.Name, "struct", comp.Struct)
			continue
		}

		// Generate templ component
		if err := generateTemplComponent(comp, expandedStructInfo, config); err != nil {
			logger.Error("Error generating templ", "component", comp.Name, "error", err)
			continue
		}

		// Generate Lua snippets
		if err := generateSnippets(comp, config); err != nil {
			logger.Error("Error generating snippets", "component", comp.Name, "error", err)
			continue
		}
	}

	// Generate struct snippets
	logger.Info("Generating struct snippets")
	if err := generateStructSnippets(structs, config); err != nil {
		logger.Error("Failed to generate struct snippets", "error", err)
	}

	logger.Info("✅ Generation complete!")

	// Run templ generate
	logger.Info("Running templ generate...")
	if err := runTemplGenerate(config.Global.Output.TemplDir, logger); err != nil {
		logger.Error("Failed to run templ generate", "error", err)
		return
	}
	logger.Info("✅ Templ generate complete!")
}

func runTemplGenerate(templDir string, logger *slog.Logger) error {
	cmd := exec.Command("templ", "generate")
	cmd.Dir = templDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Templ generate output", "output", string(output))
		return fmt.Errorf("templ generate failed: %w", err)
	}
	if len(output) > 0 {
		logger.Info("Templ output", "output", string(output))
	}
	return nil
}

func expandEmbeddedStructs(structs map[string]StructInfo) map[string]StructInfo {
	expanded := make(map[string]StructInfo)
	for name, info := range structs {
		expandedInfo := StructInfo{
			Name:   info.Name,
			Fields: []FieldInfo{},
		}
		for _, field := range info.Fields {
			if field.Name == "" && field.IsStruct {
				if embeddedStruct, exists := structs[field.Type]; exists {
					expandedInfo.Fields = append(expandedInfo.Fields, embeddedStruct.Fields...)
				}
			} else {
				expandedInfo.Fields = append(expandedInfo.Fields, field)
			}
		}
		expanded[name] = expandedInfo
	}
	return expanded
}
