package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func listNamespaces(c client.Client) (corev1.NamespaceList, error) {
	// Fetch the Namespace list instance
	backendNamespaceList := corev1.NamespaceList{}
	backendOpts := client.ListOptions{}

	// This is a hack, sometimes we can return nothing so we need to cycle till we get something
	// Fill free to tell me what I'm doing wrong here!
	for len(backendNamespaceList.Items) == 0 {
		if err := c.List(context.Background(), &backendNamespaceList, &backendOpts); err != nil {
			return backendNamespaceList, err
		}
	}

	return backendNamespaceList, nil
}
