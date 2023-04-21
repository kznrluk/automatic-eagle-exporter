package pnginfo

import (
	"fmt"
	"strconv"
	"strings"
)

type Extras struct {
	UpscaleFactor int
	Upscaler      string
	Upscaler2     string
}

func ParseExtras(input string) (*Extras, error) {
	info := &Extras{}
	parts := strings.Split(input, ", ")

	for _, part := range parts {
		keyValue := strings.Split(part, ": ")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid input format")
		}

		key := keyValue[0]
		value := keyValue[1]

		switch key {
		case "Postprocess upscale by":
			factor, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid input: Upscale factor should be an integer, got %s", value)
			}
			info.UpscaleFactor = factor
		case "Postprocess upscaler":
			info.Upscaler = value
		case "Postprocess upscaler 2":
			info.Upscaler2 = value
		default:
			continue
		}
	}

	return info, nil
}
