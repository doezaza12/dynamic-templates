package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/doezaza12/dynamic-templates/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "dynamic-templates",
	Short: "Dynamic Templates is a template generator",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("require at least 2 arguments. $REMOTE_REPO/$LOCAL_REPO, $NAME")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		noCache, err := cmd.Flags().GetBool("no-cache")
		if err != nil {
			panic(err)
		}
		fmt.Println(args)
		err = core.GitClone(args[0], noCache)
		if err != nil {
			panic(err)
		}
		valueFiles, err := cmd.Flags().GetStringArray("value")
		if err != nil {
			panic(err)
		}

		value := make(map[string]interface{})
		for _, valueFile := range valueFiles {
			rawData, err := os.ReadFile(valueFile)
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(rawData, &value)
			if err != nil {
				panic(err)
			}
		}

		// patch name into value, this is reserve variable
		value["name"] = args[1]

		core.RenderTemplate(strings.ReplaceAll(filepath.Base(args[0]), ".git", ""), "", args[1], value)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
	rootCmd.Flags().StringArray("value", []string{}, "yaml value file for use to render template")
	rootCmd.Flags().Bool("no-cache", false, "to clone remote template even it has already existed in local")
}
