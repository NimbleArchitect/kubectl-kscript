package plugin

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var scriptShort = "runs the named js script"

var scriptDescription = `not implemented`

var scriptExample = `
not implemented
`

func Script(action int, cmd *cobra.Command, kubeFlags *genericclioptions.ConfigFlags, args []string) error {
	var prog *goja.Program
	var err error
	var console jsConsole
	var dryrun bool
	// var globalObj globalObject

	// var k8s jsK8s

	log := logger{location: "Script"}
	log.Debug("Start")

	// loopinfo := ports{}
	fmt.Println("script args:", args)

	// start js code
	filename := args[0]

	cfile, err := os.ReadFile(filename)
	if err != nil {
		log.Tell("unable to read file:", err)
	}

	if cmd.Flag("dry-run") != nil {
		dryrun = cmd.Flag("dry-run").Value.String() == "true"
	}

	rt := goja.New()
	rt.SetFieldNameMapper(goja.UncapFieldNameMapper())

	builder := yamlBuilder{
		runtime: rt,
		dryrun:  dryrun,
		flags:   kubeFlags,
		action:  action,
	}

	prog, err = goja.Compile(filename, string(cfile), true)
	if err != nil {
		log.Tell("unable to compile script", err)
	}

	err = rt.Set("console", console)
	if err != nil {
		log.Yell(err)
	}

	err = rt.Set("global", builder.GlobalObj)
	if err != nil {
		log.Yell(err)
	}

	if err = rt.GlobalObject().Set("volumeClaim", builder.volumeClaim); err != nil {
		log.Yell("unable to attach function volumeClaim:", err)
	}
	if err = rt.GlobalObject().Set("volumeConfigMap", builder.volumeConfigMap); err != nil {
		log.Yell("unable to attach function volumeConfigMap:", err)
	}
	if err = rt.GlobalObject().Set("annotation", builder.annotation); err != nil {
		log.Yell("unable to attach function annotation: ", err)
	}
	if err = rt.GlobalObject().Set("create", builder.createMap); err != nil {
		log.Yell("unable to attach  function create: ", err)
	}
	if err = rt.GlobalObject().Set("deploy", builder.deploy); err != nil {
		log.Yell("unable to attach function deploy: ", err)
	}
	if err = rt.GlobalObject().Set("replica", builder.replica); err != nil {
		log.Yell("unable to attach function deploy: ", err)
	}

	out, err := rt.RunProgram(prog)

	if err != nil {
		log.Tell("script error:", err)
		log.Tell("script returns:", out)
	}

	return nil
}
