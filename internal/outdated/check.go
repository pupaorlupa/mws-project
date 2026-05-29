package outdated

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/mod/modfile"
)

const maxParallel = 8

type Result struct {
	Path    string
	Current string
	Latest  string
	Err     error
}

type Checker struct {
	http *http.Client
}

func NewChecker(client *http.Client) *Checker {
	return &Checker{http: client}
}

func (c *Checker) FetchModFile(url string) (*modfile.File, error) {
	data, err := httpGet(c.http, url)
	if err != nil {
		return nil, err
	}
	mf, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("разобрать go.mod: %w", err)
	}
	return mf, nil
}

func (c *Checker) CheckAll(reqs []*modfile.Require) []Result {
	results := make([]Result, len(reqs))
	sem := make(chan struct{}, maxParallel)
	var wg sync.WaitGroup

	for i, r := range reqs {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int, r *modfile.Require) {
			defer wg.Done()
			defer func() { <-sem }()

			latest, err := latestVersion(c.http, r.Mod.Path)
			results[i] = Result{
				Path:    r.Mod.Path,
				Current: r.Mod.Version,
				Latest:  latest,
				Err:     err,
			}
		}(i, r)
	}
	wg.Wait()
	return results
}
