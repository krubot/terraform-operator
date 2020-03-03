package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

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

func checkURL(url string, header http.Header, expectedStatus int) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("unexpected response: got %d, want %d", resp.StatusCode, expectedStatus)
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
