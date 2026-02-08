package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	detach bool
	watch  bool
)

func init() {
	devCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run in detached mode")
	devCmd.Flags().BoolVar(&watch, "watch", false, "Enable hot-reload (Docker Compose watch)")
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting development environment...")

		// Load config to show project name
		cfg, err := LoadConfig()
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if cfg.ProjectName != "" {
			fmt.Printf("Project: %s\n", cfg.ProjectName)
		}

		dockerArgs := []string{"compose", "up", "--build"}
		if detach {
			dockerArgs = append(dockerArgs, "--detach")
		}
		if watch {
			dockerArgs = append(dockerArgs, "--watch")
		}

		dockerCmd := exec.Command("docker", dockerArgs...)
		dockerCmd.Stdin = os.Stdin
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		return dockerCmd.Run()
	},
}