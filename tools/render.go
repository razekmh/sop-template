package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	srcDir         = "src"
	contentDir     = "hugo/content"
	hugoConfigSrc  = "hugo/config.yml"
	configFile     = "config.yml"
	exampleConfig  = "config.example.yml"
	placeholderOpen  = "{{"
	placeholderClose = "}}"
)

func main() {
	mode := "render"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "render":
		runRender()
	case "lint":
		runLint()
	default:
		fmt.Printf("Unknown mode: %s. Use 'render' or 'lint'.\n", mode)
		os.Exit(1)
	}
}

// ── Render ────────────────────────────────────────────────────────────────────

func runRender() {
	config := loadConfig()
	validateConfig(config)

	if err := os.MkdirAll(contentDir, 0755); err != nil {
		fatalf("Could not create content dir: %v", err)
	}

	// Render src/*.md → hugo/content/*.md
	files, err := filepath.Glob(filepath.Join(srcDir, "*.md"))
	if err != nil || len(files) == 0 {
		fatalf("No .md files found in %s/", srcDir)
	}

	for _, f := range files {
		renderFile(f, filepath.Join(contentDir, filepath.Base(f)), config)
	}

	// Render hugo/config.yml in place
	renderFile(hugoConfigSrc, hugoConfigSrc, config)

	fmt.Println("\n✓ Render complete.")
	fmt.Printf("  %d document(s) written to %s/\n", len(files), contentDir)
	fmt.Println("  Hugo config updated.")
}

func renderFile(src, dst string, config map[string]string) {
	data, err := os.ReadFile(src)
	if err != nil {
		fatalf("Could not read %s: %v", src, err)
	}

	output := substitute(string(data), config)

	// Fail fast on any unreplaced placeholders
	if idx := strings.Index(output, placeholderOpen); idx != -1 {
		// Find the placeholder for a useful error message
		end := strings.Index(output[idx:], placeholderClose)
		placeholder := output[idx : idx+end+len(placeholderClose)]
		fatalf("Unreplaced placeholder %s in %s — add it to config.yml", placeholder, src)
	}

	if err := os.WriteFile(dst, []byte(output), 0644); err != nil {
		fatalf("Could not write %s: %v", dst, err)
	}

	fmt.Printf("  rendered: %s\n", filepath.Base(src))
}

func substitute(content string, config map[string]string) string {
	for k, v := range config {
		content = strings.ReplaceAll(content, placeholderOpen+k+placeholderClose, v)
	}
	return content
}

// ── Lint ──────────────────────────────────────────────────────────────────────

func runLint() {
	config := loadConfig()

	files, _ := filepath.Glob(filepath.Join(srcDir, "*.md"))
	if len(files) == 0 {
		fatalf("No .md files found in %s/", srcDir)
	}

	failed := false
	for _, f := range files {
		data, _ := os.ReadFile(f)
		content := string(data)

		for k, v := range config {
			if strings.TrimSpace(v) == "" {
				continue
			}
			if strings.Contains(content, v) {
				fmt.Printf("✗ Lint error in %s: real value for '%s' found (\"%s\")\n", f, k, v)
				fmt.Printf("  Replace with {{%s}}\n", k)
				failed = true
			}
		}
	}

	if failed {
		fmt.Println("\nLint failed — fix the above errors before merging.")
		os.Exit(1)
	}

	fmt.Printf("✓ Lint passed — no hardcoded values found in %d file(s).\n", len(files))
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func loadConfig() map[string]string {
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error: %s not found.\n", configFile)
		fmt.Printf("Copy %s to %s and fill in your values.\n", exampleConfig, configFile)
		os.Exit(1)
	}

	var config map[string]string
	if err := yaml.Unmarshal(data, &config); err != nil {
		fatalf("Error parsing %s: %v", configFile, err)
	}

	return config
}

func validateConfig(config map[string]string) {
	hasEmpty := false
	for k, v := range config {
		if strings.TrimSpace(v) == "" {
			fmt.Printf("Error: config key '%s' is empty — fill it in before rendering.\n", k)
			hasEmpty = true
		}
	}
	if hasEmpty {
		os.Exit(1)
	}
}

func fatalf(format string, args ...any) {
	fmt.Printf("Error: "+format+"\n", args...)
	os.Exit(1)
}
