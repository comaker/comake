package cmd

import (
	"fmt"

	"github.com/comaker/comake/util"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available steps",
	Run: func(cmd *cobra.Command, args []string) {
		buildfile, err := cmd.Flags().GetString("buildfile")
		if err != nil {
			panic(err)
		}

		config, err := util.ReadBuildConfig(buildfile)
		if err != nil {
			panic(err)
		}

		fmt.Println("Build steps:")

		for _, step := range config.Steps {
			fmt.Println("-", step.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
