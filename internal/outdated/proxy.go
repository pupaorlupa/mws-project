package outdated

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/mod/module"
)

const (
	proxyBaseURL = "https://proxy.golang.org"
	userAgent    = "gomod-outdated/1.0"
)

func httpGet(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := strings.TrimSpace(string(body))
		if len(snippet) > 120 {
			snippet = snippet[:120] + "..."
		}
		return nil, fmt.Errorf("GET %s: %s: %s", url, resp.Status, snippet)
	}
	return body, nil
}

func latestVersion(client *http.Client, modPath string) (string, error) {
	escaped, err := module.EscapePath(modPath)
	if err != nil {
		return "", err
	}
	body, err := httpGet(client, fmt.Sprintf("%s/%s/@latest", proxyBaseURL, escaped))
	if err != nil {
		return "", err
	}
	var info struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", fmt.Errorf("разобрать ответ прокси: %w", err)
	}
	if info.Version == "" {
		return "", errors.New("пустое поле Version в ответе прокси")
	}
	return info.Version, nil
}
