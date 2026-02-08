package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var skipDocker bool
var skipCI bool

func init() {
	initCmd.Flags().BoolVar(&skipDocker, "skip-docker", false, "Don't run docker init")
	initCmd.Flags().BoolVar(&skipCI, "skip-ci", false, "Don't generate GitHub Actions workflow")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new DevDock project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Detecting project type...")

		projectType := detectProjectType()
		if projectType == "unknown" {
			fmt.Println("⚠️ Could not detect project type automatically.")
		}

		// Run docker init unless skipped
		if !skipDocker {
			if projectType != "unknown" {
				fmt.Printf("✨ %s detected. Using DevDock optimized setup...\n", strings.Title(projectType))
				if err := generateConfig(projectType); err != nil {
					return err
				}
			} else {
				fmt.Println("🐳 Running standard docker init...")
				dockerCmd := exec.Command("docker", "init")
				dockerCmd.Stdin = os.Stdin
				dockerCmd.Stdout = os.Stdout
				dockerCmd.Stderr = os.Stderr
				if err := dockerCmd.Run(); err != nil {
					return fmt.Errorf("docker init failed: %w", err)
				}
			}
		}

		// Generate devdock.yaml
		cfg := Config{
			ProjectName: GetProjectName(),
			Type:        projectType,
		}
		if err := WriteConfig(cfg); err != nil {
			return err
		}
		fmt.Println("📄 Created devdock.yaml")

		// Generate basic GitHub Actions CI
		if !skipCI {
			fmt.Println("⚙️ Creating GitHub Actions workflow...")
			workflowDir := ".github/workflows"
			os.MkdirAll(workflowDir, 0755)
			ciContent := `name: CI
on: [push, pull_request]
jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build Docker image
        run: docker build -t ${{ secrets.DOCKER_REPO }} .
`
			ciPath := filepath.Join(workflowDir, "ci.yaml")
			if err := os.WriteFile(ciPath, []byte(ciContent), 0644); err != nil {
				return err
			}
			fmt.Println("✅ Created .github/workflows/ci.yaml")
		}

		fmt.Println("🚀 DevDock initialization complete! Try 'devdock dev'")
		return nil
	},
}

func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func detectProjectType() string {
	if checkFileExists("package.json") {
		// Differentiate between JS frameworks based on build tools
		if isNextApp() {
			return "next" // Next.js (SSR, Port 3000)
		}
		if isViteApp() {
			return "vite" // React/Vue/Svelte + Vite (CSR, Port 5173)
		}
		return "node" // Generic Node.js (Express/Nest, Port 3000)
	}
	if checkFileExists("requirements.txt") || checkFileExists("pyproject.toml") {
		return "python"
	}
	if checkFileExists("go.mod") {
		return "go"
	}
	if checkFileExists("pom.xml") || checkFileExists("build.gradle") {
		return "java"
	}
	if checkFileExists("composer.json") {
		return "php"
	}
	return "unknown"
}

func isNextApp() bool {
	return checkDependency("next")
}

func isViteApp() bool {
	return checkDependency("vite")
}

func checkDependency(depName string) bool {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return false
	}
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}
	for k := range pkg.Dependencies {
		if strings.Contains(k, depName) {
			return true
		}
	}
	for k := range pkg.DevDependencies {
		if strings.Contains(k, depName) {
			return true
		}
	}
	return false
}

func generateConfig(projectType string) error {
	var dockerfile, compose string
	dockerignore := nodeDockerignore // Default ignore list

	switch projectType {
	case "vite":
		dockerfile = viteDockerfile
		compose = viteCompose
	case "next":
		dockerfile = nextDockerfile
		compose = nextCompose
	case "node":
		dockerfile = nodeDockerfile
		compose = nodeCompose
	case "python":
		dockerfile = pythonDockerfile
		compose = pythonCompose
	case "go":
		dockerfile = goDockerfile
		compose = goCompose
	case "java":
		dockerfile = javaDockerfile
		compose = javaCompose
	case "php":
		dockerfile = phpDockerfile
		compose = phpCompose
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}

	if err := os.WriteFile("Dockerfile", []byte(dockerfile), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	if err := os.WriteFile("compose.yaml", []byte(compose), 0644); err != nil {
		return fmt.Errorf("failed to write compose.yaml: %w", err)
	}
	if err := os.WriteFile(".dockerignore", []byte(dockerignore), 0644); err != nil {
		return fmt.Errorf("failed to write .dockerignore: %w", err)
	}
	fmt.Printf("📄 Created Dockerfile, compose.yaml, .dockerignore for %s\n", projectType)
	return nil
}
