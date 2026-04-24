package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	srcDir           = "src"
	contentDir       = "hugo/content"
	hugoConfigSrc    = "hugo/config.yml"
	configFile       = "config.yml"
	exampleConfig    = "config.example.yml"
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

	// Fail fast on any unreplaced config placeholders ({{key}}), but do not treat
	// Hugo shortcodes like {{< columns >}} or {{% hint %}} as config keys.
	if bad := findUnreplacedConfigPlaceholder(output, config); bad != "" {
		fatalf("Unreplaced placeholder %s in %s — add it to config.yml", bad, src)
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

// findUnreplacedConfigPlaceholder returns the first {{...}} that should come from
// config.yml but is still present, or malformed markup. Hugo shortcodes that use
// {{< ... >}} or {{% ... %}} are skipped so theme markup can live in src/.
func findUnreplacedConfigPlaceholder(s string, config map[string]string) string {
	validKeys := make(map[string]struct{}, len(config))
	for k := range config {
		validKeys[k] = struct{}{}
	}
	i := 0
	for i < len(s) {
		j := strings.Index(s[i:], placeholderOpen)
		if j < 0 {
			return ""
		}
		j += i
		afterOpen := j + len(placeholderOpen)
		if afterOpen >= len(s) {
			return "{{"
		}
		rest := s[afterOpen:]
		switch rest[0] {
		case '<':
			close := strings.Index(rest, ">}}")
			if close < 0 {
				return s[j:min(j+50, len(s))]
			}
			i = afterOpen + close + len(">}}")
			continue
		case '%':
			close := strings.Index(rest, "%}}")
			if close < 0 {
				return s[j:min(j+50, len(s))]
			}
			i = afterOpen + close + len("%}}")
			continue
		}
		closeRel := strings.Index(rest, placeholderClose)
		if closeRel < 0 {
			return s[j:]
		}
		key := rest[:closeRel]
		if strings.TrimSpace(key) != key || key == "" {
			return placeholderOpen + key + placeholderClose
		}
		if _, ok := validKeys[key]; !ok {
			return placeholderOpen + key + placeholderClose
		}
		// Known key still present after substitution — should not happen
		return placeholderOpen + key + placeholderClose
	}
	return ""
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
