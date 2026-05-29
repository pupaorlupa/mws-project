package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pupaorlupa/mws-project/internal/outdated"
)

const httpTimeout = 20 * time.Second

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	branch := fs.String("branch", "HEAD", "Ветка, тег или коммит, из которого берётся go.mod")
	doUpdate := fs.Bool("update", false, "Интерактивно обновить устаревшие require и записать go.mod локально")
	outPath := fs.String("o", "go.mod.updated", "Путь для сохранения обновлённого go.mod (используется с --update)")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Использование: ./mws-project [флаги] <git-repo-url>\n\nФлаги:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("требуется ровно один URL репозитория")
	}

	rawURL, err := outdated.GoModRawURL(fs.Arg(0), *branch)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: httpTimeout}
	checker := outdated.NewChecker(client)

	fmt.Printf("Загружаю %s\n", rawURL)
	mf, err := checker.FetchModFile(rawURL)
	if err != nil {
		return err
	}

	outdated.PrintHeader(mf)
	results := checker.CheckAll(mf.Require)
	outdatedList := outdated.PrintTable(results)

	if *doUpdate && len(outdatedList) > 0 {
		if err := outdated.Interactive(mf, outdatedList, *outPath, os.Stdin); err != nil {
			return err
		}
	}
	return nil
}
