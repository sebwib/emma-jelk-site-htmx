package dump

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func JSONAsText(v any) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

func JSON(v any) {
	data, _ := json.MarshalIndent(v, "", "  ")
	log.Println(string(data))
}

func MapsAsCSV[V any](maps []map[string]V) string {
	if len(maps) == 0 {
		return ""
	}

	var keys []string
	for k := range maps[0] {
		keys = append(keys, fmt.Sprintf("%v", k))
	}
	var rows [][]string
	rows = append(rows, keys)

	for _, v := range maps {
		var row []string
		for _, k := range keys {
			row = append(row, fmt.Sprint(v[k]))
		}
		rows = append(rows, row)
	}
	var builder strings.Builder
	for _, row := range rows {
		builder.WriteString(strings.Join(row, ","))
		builder.WriteString("\n")
	}
	return builder.String()
}
