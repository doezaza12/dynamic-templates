package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/doezaza12/dynamic-templates/core"
	"github.com/doezaza12/dynamic-templates/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "dynamic-templates",
	Short: "Dynamic Templates is a template generator",
	Long:  ``,
	Example: `dynamic-templates /path/to/your/local/template rendered-template
dynamic-templates git@github.com/doezaza12/dummy-template.git rendered-template --values core-values.yaml --values template-values.yaml
dynamic-templates https://github.com/doezaza12/dummy-template.git rendered-template --values template-values.yaml
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("require at least 2 arguments. $REMOTE_REPO/$LOCAL_REPO, $NAME")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// get flag values
		noCache, err := cmd.Flags().GetBool("no-cache")
		if err != nil {
			panic(err)
		}
		valueFiles, err := cmd.Flags().GetStringArray("values")
		if err != nil {
			panic(err)
		}
		outputDir, err := cmd.Flags().GetString("out")
		if err != nil {
			panic(err)
		}

		var templateFullPath string

		if util.IsRemoteTemplate(args[0]) {
			// clone remote template into work dir
			err = core.GitClone(args[0], noCache)
			if err != nil {
				panic(err)
			}
			// by default, remote template will store at $HOME/dynamic-templates/$REMOTE_TEMPLATE_NAME
			templateFullPath = filepath.Base(args[0])
		} else {
			//  local template path
			templateFullPath = args[0]
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

		core.RenderTemplate(templateFullPath, outputDir, args[1], value)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringArray("values", []string{}, "yaml value files for use to render template, example = --values project-val.yaml --values kustomize-val.yaml")
	rootCmd.Flags().Bool("no-cache", false, "to clone remote template even it has already existed in local (only affect remote template), default = false")
	rootCmd.Flags().String("out", "", "specify output directory for rendered template (user must have RW permission)")
}
