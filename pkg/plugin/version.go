package plugin

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var versionsShort = "Display container versions and mount points"

var helpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
More information, documentation and examples at: https://nimblearchitect.github.io/kubectl-script/
 find this program useful? Please consider donating: https://nimblearchitect.github.io/kubectl-script/donations/

`

func Version(cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	// 1234567890123456789012345678901234567890123456789012345678901234567890123456789
	fmt.Printf(`kubectl-script kubernetes container viewer

version %s

the latest version can be found at: 
	https://nimblearchitect.github.io/kubectl-script/downloads/

to view the documentation:
	https://nimblearchitect.github.io/kubectl-script

or to raise issues: 
   https://github.com/NimbleArchitect/kubectl-script

if you find this program useful please consider saying thanks I can be reached
 on twitter @nimblearchitect or you can buy me a coffee:
	https://nimblearchitect.github.io/kubectl-script/donations/


if your just after the version string use: kubectl-script -v

`, cmd.Parent().Version)
	return nil
}
