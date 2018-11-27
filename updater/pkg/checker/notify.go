package checker

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type NotifyMsg struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

func NotifyUpdate(url, service , version string) error {
	notifyMsg := NotifyMsg{Service:service, Version:version}
	marshaledObject, err := json.MarshalIndent(notifyMsg, "", "  ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(marshaledObject))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	c := &http.Client{}
	_, err = c.Do(req)
	if err != nil {
		return err
	}

	return nil
}
