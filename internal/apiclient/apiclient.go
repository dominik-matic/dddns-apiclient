package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/dominik-matic/dddns-apiclient/internal/models"
)

func SendRequest(token, mode, domain string) error {
	switch mode {
	case MODE_UPDATE:
		return SendUpdateRequest(token, domain)
	case MODE_DELETE:
		return SendDeleteRequest(token, domain)
	}
	return errors.New("invalid mode even after mode validation")
}

func SendUpdateRequest(token, domain string) error {
	public_ip, err := getMyPublicIp()
	if err != nil {
		return err
	}

	requestBody := models.UpdatePayload{
		Name:  domain,
		Value: public_ip,
	}

	return sendRequest(http.MethodPost, token, requestBody)
}

func SendDeleteRequest(token, domain string) error {
	requestBody := models.DeletePayload{
		Name: domain,
	}
	return sendRequest(http.MethodDelete, token, requestBody)
}

func sendRequest(method, token string, body any) error {
	url := "https://" + DDDNS_HOST + ":" + DDDNS_PORT

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("server returned status: %s, but failed to read body: %w", res.Status, err)
		}
		return fmt.Errorf("server returned status %s: %s", res.Status, string(bodyBytes))
	}

	return nil

}

// In the future maybe get rid of the curl dependency
// For now it's the fastest way to do this crap
func getMyPublicIp() (string, error) {
	var err error
	var out []byte
	for _, provider := range PUBLIC_IP_PROVIDERS {
		out, err = exec.Command("/usr/bin/curl", "--silent", "--ipv4", provider).Output()
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("curl failed to get public IP. %w", err)
	}
	public_ip := strings.TrimSpace(string(out))
	return public_ip, nil
}
