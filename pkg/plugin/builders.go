package plugin

import (
	"fmt"
)

type appObject struct {
	name          string
	suppiledState *state
	template      *appItems
	service       *appItems
	global        *appItems
	Pod           rawPod
}

func (o *appObject) init() {
	// TODO: check connected k8s version and adjust apiVersion tag accordingly

	m := make(map[string]interface{})
	m["spec"] = make(map[string]interface{})
	o.template = &appItems{data: m}

	m = make(map[string]interface{})
	m["spec"] = make(map[string]interface{})
	o.service = &appItems{data: m}
	// o = &appItems{data: m}

}

func buildTemplate(app *appObject) {
	state := app.suppiledState

	app.template.add("name", app.name, []string{"metadata"})
	app.template.add("namespace", "default", []string{"metadata"})

	path := []string{"spec", "template"}

	if state.Replicas == nil {
		app.template.add("replicas", 1, []string{"spec"})
	} else {
		app.template.add("replicas", state.Replicas, []string{"spec"})
	}

	for n, e := range state.Env {
		app.global.add(n, e, []string{"env"})
	}

	// set labels
	isSelectorSet := false
	if len(state.Labels) == 0 {
		state.Labels = make(map[string]interface{})
		state.Labels["app"] = app.name
	}
	for n, e := range state.Labels {
		app.template.add(n, e, append(path, "metadata", "labels"))
		if !isSelectorSet {
			if len(state.Global.Selector) == 0 {
				app.template.add(n, e, []string{"spec", "selector", "matchLabels"})
				isSelectorSet = true
			} else if state.Global.Selector == n {
				app.template.add(n, e, []string{"spec", "selector", "matchLabels"})
				isSelectorSet = true
			}
		}
	}

	buildDeploymentContainers(app)
}

func buildDeploymentContainers(app *appObject) {
	state := app.suppiledState

	path := []string{"spec", "template", "spec"}

	arrayInit := 0
	array := 0
	for _, v := range state.Image {
		var localPath []string
		if v.IsInit {
			localPath = append(path, fmt.Sprintf("initContainers[%d]", arrayInit))
			arrayInit++
		} else {
			localPath = append(path, fmt.Sprintf("containers[%d]", array))
			array++
		}

		app.template.add("name", app.name, localPath)
		app.template.add("image", v.Location, localPath)
		app.template.add("command", v.Command, localPath)

		envId := 0
		for n, e := range state.Env {
			// fmt.Println(">> global:", localPath, n, e)
			app.template.add("name", (n), append(localPath, fmt.Sprintf("env[%d]", envId)))
			app.template.add("value", (e), append(localPath, fmt.Sprintf("env[%d]", envId)))
			envId++
		}

		for n, e := range v.Env {
			app.template.add(n, (e), append(localPath, fmt.Sprintf("env[%d]", envId)))
			envId++
		}

		// resources:
		//     limits:
		//       memory: "2Gi"
		//       cpu: "1000m"
		//     requests:
		//       memory: "500Mi"
		//       cpu: "100m"
		// global resources
		if state.Memory["min"] != nil {
			app.template.add("memory", state.Memory["min"], append(localPath, "resources", "requests"))
		}
		if state.Memory["max"] != nil {
			app.template.add("memory", state.Memory["max"], append(localPath, "resources", "limits"))
		}
		if state.Cpu["min"] != nil {
			app.template.add("cpu", state.Cpu["min"], append(localPath, "resources", "requests"))
		}
		if state.Cpu["max"] != nil {
			app.template.add("cpu", state.Cpu["max"], append(localPath, "resources", "limits"))
		}
		// local image resources
		if v.Memory["min"] != nil {
			app.template.add("memory", v.Memory["min"], append(localPath, "resources", "requests"))
		}
		if v.Memory["max"] != nil {
			app.template.add("memory", v.Memory["max"], append(localPath, "resources", "limits"))
		}
		if v.Cpu["min"] != nil {
			app.template.add("cpu", v.Cpu["min"], append(localPath, "resources", "requests"))
		}
		if v.Cpu["max"] != nil {
			app.template.add("cpu", v.Cpu["max"], append(localPath, "resources", "limits"))
		}
	}
}

func buildService(app *appObject) {
	state := app.suppiledState
	if state.Service != nil {
		for _, v := range state.Service {
			for n, e := range v {
				app.service.add(n, e, []string{"metadata", "annotations"})
			}
		}
	}
}
