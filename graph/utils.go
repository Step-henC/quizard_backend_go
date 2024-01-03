package graph

import (
	"strings"
)

func StringifyQuery(query string) string {

	var builder strings.Builder
	builder.WriteString(query)

	return builder.String()
}
