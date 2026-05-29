package outdated

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
)

func Interactive(mf *modfile.File, outdated []Result, outPath string, in io.Reader) error {
	reader := bufio.NewReader(in)

	updated := 0
	for _, o := range outdated {
		ok, err := askYesNo(reader, fmt.Sprintf("Обновить %s: %s -> %s ?", o.Path, o.Current, o.Latest))
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		if err := mf.AddRequire(o.Path, o.Latest); err != nil {
			fmt.Fprintf(os.Stderr, "  не удалось: %v\n", err)
			continue
		}
		updated++
	}

	if updated == 0 {
		fmt.Println("Изменений нет.")
		return nil
	}

	mf.Cleanup()
	out, err := mf.Format()
	if err != nil {
		return fmt.Errorf("отформатировать go.mod: %w", err)
	}
	if err := os.WriteFile(outPath, out, 0o644); err != nil {
		return fmt.Errorf("записать %s: %w", outPath, err)
	}
	fmt.Printf("Обновлённый go.mod (%d изменени(й)) записан в %s\n", updated, outPath)
	return nil
}

func askYesNo(r *bufio.Reader, question string) (bool, error) {
	fmt.Printf("%s [y/N]: ", question)
	line, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes" || ans == "д" || ans == "да", nil
}
