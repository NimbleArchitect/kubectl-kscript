package plugin

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
)

func (b *yamlBuilder) pushChange(k8info schema.GroupVersionResource, listobj map[string]interface{}) error {
	client, err := newClient(b.flags)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	res := client.Resource(k8info).Namespace("default")

	list := unstructured.Unstructured{}

	list.SetUnstructuredContent(listobj)

	// TODO: need a better comparasion, using name on its own is poor practice
	name := listobj["metadata"].(map[string]interface{})["name"].(string)
	out, _ := res.Get(ctx, name, metav1.GetOptions{})

	// TODO: how do we uninstall or force upgrade/force reinstall
	switch b.action {
	case ACTION_INSTALL:
		if out == nil {
			_, err = res.Create(ctx, &list, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("install error: %w", err)
			}
		} else {
			// TODO: what happens when the resource already exists??
			return errors.New("resourse already exists")
		}
	case ACTION_UPDATE:
		if out != nil {
			_, err = res.Update(ctx, &list, metav1.UpdateOptions{})
		} else {
			_, err = res.Create(ctx, &list, metav1.CreateOptions{})
		}
		if err != nil {
			return fmt.Errorf("install error: %w", err)
		}
	case ACTION_REMOVE:
		if out != nil {
			err = res.Delete(ctx, name, metav1.DeleteOptions{})
		} else {
			// TODO: what happens if the resource dosent exist
			return errors.New("resource doesn't exist")
		}
		if err != nil {
			return fmt.Errorf("install error: %w", err)
		}
	}

	// fmt.Println(out.GetUID())
	return nil
}

func newClient(kubeFlags *genericclioptions.ConfigFlags) (dynamic.Interface, error) {
	// config, err := rest.InClusterConfig()

	config, err := kubeFlags.ToRESTConfig()

	if err != nil {
		return nil, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}
