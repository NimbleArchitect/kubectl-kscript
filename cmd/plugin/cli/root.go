package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/NimbleArchitect/kubectl-kscript/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// auto updated version via gorelaser
var version = "0.0.0"

var rootShort = "View pod information at the container level"

var rootDescription = `  Deploy applictaions from your pipeline using kscript.
	
 Suggestions and improvements can be made by raising an issue here: 
    https://github.com/NimbleArchitect/kubectl-kscript

`

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kubectl-kscript",
		Short:         rootShort,
		Long:          fmt.Sprintf("%s\n\n%s", rootShort, rootDescription),
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cobra.OnInitialize(initConfig)

	if strings.ToLower(os.Getenv("KSCRIPT_LOG")) == "debug" {
		plugin.LogDebug = true
	}

	plugin.InitCommand(cmd)

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
