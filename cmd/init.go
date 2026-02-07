package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

		var projectType string
		if _, err := os.Stat("package.json"); err == nil {
			projectType = "node"
		} else if _, err := os.Stat("requirements.txt"); err == nil {
			projectType = "python"
		} else if _, err := os.Stat("pyproject.toml"); err == nil {
			projectType = "python"
		} else if _, err := os.Stat("go.mod"); err == nil {
			projectType = "go"
		} else {
			return fmt.Errorf("could not detect project type (no known files found)")
		}
		fmt.Printf("✅ Detected %s project\n", projectType)

		// Run docker init unless skipped
		if !skipDocker {
			fmt.Println("🐳 Running docker init...")
			dockerCmd := exec.Command("docker", "init")
			dockerCmd.Stdin = os.Stdin
			dockerCmd.Stdout = os.Stdout
			dockerCmd.Stderr = os.Stderr
			if err := dockerCmd.Run(); err != nil {
				return fmt.Errorf("docker init failed: %w", err)
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
