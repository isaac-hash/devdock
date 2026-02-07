package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build production Docker image with standardized tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, _ := LoadConfig()
		projectName := cfg.ProjectName
		if projectName == "" {
			projectName = GetProjectName()
		}

		tags := []string{fmt.Sprintf("%s:latest", projectName)}

		// Add git short hash if in a repo
		if gitHash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
			hash := string(gitHash)
			hash = hash[:len(hash)-1] // trim newline
			tags = append(tags, fmt.Sprintf("%s:%s", projectName, hash))
		}

		tagArgs := []string{}
		for _, t := range tags {
			tagArgs = append(tagArgs, "-t", t)
		}

		buildArgs := append([]string{"build"}, tagArgs...)
		buildArgs = append(buildArgs, ".")

		fmt.Printf("🔨 Building image with tags: %v\n", tags)
		buildCmd := exec.Command("docker", buildArgs...)
		buildCmd.Stdin = os.Stdin
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		return buildCmd.Run()
	},
}