package outdated

import (
	"fmt"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func PrintHeader(mf *modfile.File) {
	fmt.Printf("\nМодуль:        %s\n", mf.Module.Mod.Path)
	if mf.Go != nil {
		fmt.Printf("Версия Go:     %s\n", mf.Go.Version)
	}
	fmt.Printf("Зависимостей:  %d\n\n", len(mf.Require))
}

func PrintTable(results []Result) []Result {
	maxPath, maxCur := len("МОДУЛЬ"), len("ТЕКУЩАЯ")
	for _, r := range results {
		if len(r.Path) > maxPath {
			maxPath = len(r.Path)
		}
		if len(r.Current) > maxCur {
			maxCur = len(r.Current)
		}
	}

	fmt.Printf("%-*s  %-*s  %s\n", maxPath, "МОДУЛЬ", maxCur, "ТЕКУЩАЯ", "ПОСЛЕДНЯЯ")
	fmt.Println(strings.Repeat("-", maxPath+maxCur+30))

	var outdated []Result
	for _, r := range results {
		switch {
		case r.Err != nil:
			fmt.Printf("%-*s  %-*s  ! %v\n", maxPath, r.Path, maxCur, r.Current, r.Err)
		case semver.Compare(r.Latest, r.Current) > 0:
			fmt.Printf("%-*s  %-*s  %s  (есть обновление)\n", maxPath, r.Path, maxCur, r.Current, r.Latest)
			outdated = append(outdated, r)
		default:
			fmt.Printf("%-*s  %-*s  %s\n", maxPath, r.Path, maxCur, r.Current, r.Latest)
		}
	}

	fmt.Printf("\nМожно обновить: %d модул(я/ей).\n", len(outdated))
	return outdated
}
