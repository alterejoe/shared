package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	Source          SourceConfig      `yaml:"source"`
	Output          OutputConfig      `yaml:"output"`
	SnippetDefaults map[string]string `yaml:"snippet_defaults"`
	FieldDefaults   map[string]string `yaml:"field_defaults"`
}

type SourceConfig struct {
	StructsPath   string `yaml:"structs_path"`
	StructsImport string `yaml:"structs_import"` // ← Add this
}

type OutputConfig struct {
	TemplDir    string `yaml:"templ_dir"`
	SnippetsDir string `yaml:"snippets_dir"`
}

type ComponentConfig struct {
	Name     string              `yaml:"name"`
	Struct   string              `yaml:"struct"`
	Template string              `yaml:"template"`
	SvgPaths map[string][]string `yaml:"svg_paths"` // Add this
	Variants []VariantConfig     `yaml:"variants"`
}

type VariantConfig struct {
	Name   string            `yaml:"name"`
	Prefix string            `yaml:"prefix"`
	Props  map[string]string `yaml:"props"`
}

type Config struct {
	Global     GlobalConfig
	Components []ComponentConfig
}

func LoadConfig(configDir string) (*Config, error) {
	// Load global config
	globalPath := filepath.Join(configDir, "_global.yaml")
	globalData, err := os.ReadFile(globalPath)
	if err != nil {
		return nil, err
	}
	var global GlobalConfig
	if err := yaml.Unmarshal(globalData, &global); err != nil {
		return nil, err
	}

	// Load all component configs
	components := []ComponentConfig{}

	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		// Skip global config
		if entry.Name() == "_global.yaml" {
			continue
		}

		// Load component config
		componentPath := filepath.Join(configDir, entry.Name())
		fmt.Printf("Loading config: %s\n", entry.Name()) // ← Add this

		componentData, err := os.ReadFile(componentPath)
		if err != nil {
			fmt.Printf("  ⚠️  Failed to read %s: %v\n", entry.Name(), err) // ← Add this
			continue
		}

		var component ComponentConfig
		if err := yaml.Unmarshal(componentData, &component); err != nil {
			fmt.Printf("  ⚠️  Failed to parse %s: %v\n", entry.Name(), err) // ← Add this
			continue
		}

		fmt.Printf("  ✅ Loaded: %s (struct: %s)\n", component.Name, component.Struct) // ← Add this
		components = append(components, component)
	}

	return &Config{
		Global:     global,
		Components: components,
	}, nil
}
