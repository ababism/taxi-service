package postgres

import (
	"fmt"
	"strings"
)

func formatQuery(q string) string {
	return fmt.Sprintf("SQL Query: %s", strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " "))
}
