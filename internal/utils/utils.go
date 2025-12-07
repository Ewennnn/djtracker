package utils

import "strings"

func SafeTrim(data *string) *string {
	if data == nil {
		return nil
	}
	v := strings.TrimSpace(*data)
	if v == "" {
		return nil
	}
	return &v
}
