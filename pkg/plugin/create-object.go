package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	ACTION_INSTALL int = 1
	ACTION_UPDATE  int = 2
	ACTION_REMOVE  int = 3
)

type yamlBuilder struct {
	runtime   *goja.Runtime
	GlobalObj globalObject
	dryrun    bool
	flags     *genericclioptions.ConfigFlags
	action    int
}

type state struct {
	Cpu      map[string]interface{}
	Env      map[string]interface{}
	Global   *globalObject
	Image    []imageState
	Labels   map[string]interface{}
	Memory   map[string]interface{}
	Mount    []map[string]interface{}
	Replicas *int
	Service  []map[string]interface{}
}

type imageState struct {
	Cpu      map[string]interface{}
	Command  interface{}
	Env      map[string]interface{}
	IsInit   bool
	Location string
	Name     string
	Memory   map[string]interface{}
	Mount    []map[string]interface{}
	Port     []map[string]interface{}
}

func (b *yamlBuilder) createMap(name string, value goja.Value) *appObject {
	var appOut appObject
	var newState state

	switch val := value.Export().(type) {
	case string:
		newState.Image = make([]imageState, 1)
		newState.Image[0].Location = value.String()

	case map[string]interface{}:
		object := value.ToObject(b.runtime)
		out, err := object.MarshalJSON()
		if err != nil {
			fmt.Println("create object decode error:", err)
			// panic(err)
		}
		err = json.Unmarshal(out, &newState)
		if err != nil {
			fmt.Println("create decode error:", err)
		}

	default:
		panic(fmt.Sprintf("recieved %T, expected string or object", val))
	}

	appOut.init()
	newState.Global = &b.GlobalObj
	appOut.suppiledState = &newState
	appOut.name = name

	buildTemplate(&appOut)
	buildService(&appOut)

	// create reference to link up appOut.Pod.items with deployment template
	podspec := appOut.template.data["spec"].(map[string]interface{})
	appOut.Pod.items = &appItems{data: podspec["template"].(map[string]interface{})}
	appOut.Pod.hasChanged = &appOut.template.hasChanged

	return &appOut
}

func (b *yamlBuilder) deploy(object appObject) error {
	var err error

	object.template.add("apiVersion", "apps/v1", []string{})
	object.template.add("kind", "Deployment", []string{})

	printOut(*object.template)
	if !b.dryrun {
		k8info := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
		err = b.pushChange(k8info, object.template.data)
	}

	object.service.add("apiVersion", "v1", []string{})
	object.service.add("kind", "Service", []string{})
	printOut(*object.service)

	return err
}

func (b *yamlBuilder) replica(object appObject) error {
	var err error

	object.template.add("apiVersion", "apps/v1", []string{})
	object.template.add("kind", "ReplicaSet", []string{})

	printOut(*object.template)
	if !b.dryrun {
		k8info := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "replicasets"}
		err = b.pushChange(k8info, object.template.data)
	}

	object.service.add("apiVersion", "v1", []string{})
	object.service.add("kind", "Service", []string{})
	printOut(*object.service)

	return err
}

func (b *yamlBuilder) annotation(name goja.Value, value goja.Value) map[string]interface{} {
	// func(name goja.Value, value goja.Value) map[string]interface{} {
	// fmt.Println("builtin annotation:", name.String())
	strName := name.String()
	if len(strName) == 0 {
		panic("invalid path used for annotation")
	}

	out := make(map[string]interface{})
	out[strName] = value.Export()
	return out
	// }
}

func (b *yamlBuilder) volumeClaim(value goja.Value) *goja.Object {
	fmt.Println("builtin claim")
	var v map[string]string
	err := json.Unmarshal([]byte(""), &v)
	if err != nil {
		fmt.Println("error:", err)
	}

	jsVal := b.runtime.ToValue(v)

	return jsVal.ToObject(b.runtime)
}

func (b *yamlBuilder) volumeConfigMap(path goja.Value, mapName goja.Value, extras ...goja.Value) *goja.Object {
	v := make(map[string]goja.Value)

	v["volumeMounts.mountPath"] = path.ToString()
	v["volumeMounts.name"] = mapName.ToString()

	v["volumes.name"] = mapName.ToString()
	v["volumes.configMap.name"] = mapName.ToString()

	// err := json.Unmarshal([]byte(""), &v)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }

	// fmt.Println(v)
	jsVal := b.runtime.ToValue(v)
	jsObj := jsVal.ToObject(b.runtime)

	// val, _ := jsObj.MarshalJSON()

	// fmt.Println("configmap:", string(val))

	return (jsObj)
}
