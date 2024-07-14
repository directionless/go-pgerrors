package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	// sourceURL = "https://github.com/postgres/postgres/blob/master/src/backend/utils/errcodes.txt"
	sourceURL = "https://raw.githubusercontent.com/postgres/postgres/master/src/backend/utils/errcodes.txt"
)

func main() {
	ctx := context.Background()

	pgErrors := map[string]pgErrType{}
	if err := fetchAndParse(ctx, pgErrors); err != nil {
		fmt.Printf("Failed: %s\n", err)
		os.Exit(1)
	}

	if err := outputLookupSqlCode(ctx, "lookup-sql-code.go", pgErrors); err != nil {
		fmt.Printf("Failed: %s\n", err)
		os.Exit(1)

	}

}

func fetchAndParse(ctx context.Context, pgErrors map[string]pgErrType) error {
	req, err := http.NewRequestWithContext(ctx, "GET", sourceURL, nil)
	if err != nil {
		return err
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting remote URL: got status %d", resp.StatusCode)
	}

	rdr := bufio.NewReader(resp.Body)
	for {
		line, err := rdr.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		pgErr, ok := newPgError(line)
		if !ok {
			continue
		}

		pgErrors[pgErr.sqlstate] = pgErr
	}

	return nil
}

func outputLookupSqlCode(ctx context.Context, filename string, pgErrors map[string]pgErrType) error {
	out, err := createGoFile("pgerrors", filename)
	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		os.Exit(1)
	}

	keys := make([]string, 0, len(pgErrors))

	for k := range pgErrors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(out, "func LookupSqlCode(sqlcode string) string {\n")
	fmt.Fprintf(out, "switch sqlcode {\n")
	for _, sqlcode := range keys {
		fmt.Fprintf(out, "case `%s`: return `%s`\n", sqlcode, pgErrors[sqlcode].Name())
	}
	fmt.Fprintf(out, "default:\n")
	fmt.Fprintf(out, "return ``\n")
	fmt.Fprintf(out, "}\n")

	fmt.Fprintf(out, "}\n")

	return out.Close()
}
