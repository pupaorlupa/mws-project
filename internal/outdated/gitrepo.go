package outdated

import (
	"fmt"
	"net/url"
	"strings"
)

func GoModRawURL(repo, ref string) (string, error) {
	repo = strings.TrimSpace(repo)
	repo = strings.TrimSuffix(repo, ".git")

	if !strings.Contains(repo, "://") {
		repo = "https://" + repo
	}

	u, err := url.Parse(repo)
	if err != nil {
		return "", err
	}
	if u.Host != "github.com" {
		return "", fmt.Errorf("поддерживается только github.com, получено %q", u.Host)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("URL должен содержать owner и repo: %q", repo)
	}
	if ref == "" {
		ref = "HEAD"
	}
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/go.mod", parts[0], parts[1], ref), nil
}
