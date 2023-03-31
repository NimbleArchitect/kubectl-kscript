package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func InitCommand(rootCmd *cobra.Command) {

	log := logger{location: "InitSubCommands"}
	log.Debug("Start")

	KubernetesConfigFlags := genericclioptions.NewConfigFlags(false)
	rootCmd.SetHelpTemplate(helpTemplate)

	// install
	var cmdInstall = &cobra.Command{
		Use:     "install",
		Short:   scriptShort,
		Long:    fmt.Sprintf("%s\n\n%s", scriptShort, scriptDescription),
		Example: fmt.Sprintf(scriptExample, rootCmd.CommandPath()),
		// Aliases: []string{"script"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Script(ACTION_INSTALL, cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdInstall.Flags())
	// addCommonFlags(cmdScript)

	cmdInstall.Flags().BoolP("dry-run", "", false, `dont write changes to kubernetes`)
	rootCmd.AddCommand(cmdInstall)

	// install and upgrade
	var cmdUpdate = &cobra.Command{
		Use:     "update",
		Short:   scriptShort,
		Long:    fmt.Sprintf("%s\n\n%s", scriptShort, scriptDescription),
		Example: fmt.Sprintf(scriptExample, rootCmd.CommandPath()),
		// Aliases: []string{"script"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Script(ACTION_UPDATE, cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdUpdate.Flags())
	// addCommonFlags(cmdScript)

	cmdUpdate.Flags().BoolP("dry-run", "", false, `dont write changes to kubernetes`)
	rootCmd.AddCommand(cmdUpdate)

	// remove
	var cmdRemove = &cobra.Command{
		Use:     "remove",
		Short:   scriptShort,
		Long:    fmt.Sprintf("%s\n\n%s", scriptShort, scriptDescription),
		Example: fmt.Sprintf(scriptExample, rootCmd.CommandPath()),
		// Aliases: []string{"script"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Script(ACTION_REMOVE, cmd, KubernetesConfigFlags, args); err != nil {
				return err
			}

			return nil
		},
	}
	KubernetesConfigFlags.AddFlags(cmdRemove.Flags())
	// addCommonFlags(cmdScript)

	cmdRemove.Flags().BoolP("dry-run", "", false, `dont write changes to kubernetes`)
	rootCmd.AddCommand(cmdRemove)

}
