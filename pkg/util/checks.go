package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func CheckURL(url string, header http.Header, expectedStatus int) error {
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
