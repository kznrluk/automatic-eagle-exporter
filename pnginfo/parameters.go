package pnginfo

import (
	"fmt"
	"strconv"
	"strings"
)

type Parameters struct {
	Prompt         string
	NegativePrompt string
	Steps          int
	Sampler        string
	CFGScale       float64
	Seed           int64
	Size           string
	ModelHash      string
	Model          string

	Raw string
}

func ParseParameters(input string) (*Parameters, error) {
	lines := strings.Split(input, "\n")
	if len(lines) < 1 || len(lines) > 3 {
		return nil, fmt.Errorf("invalid input: expected 1 to 3 lines, got %d", len(lines))
	}
	info := &Parameters{}
	if len(lines) >= 2 {
		if strings.HasPrefix(lines[1], "Negative prompt: ") {
			info.NegativePrompt = strings.TrimPrefix(lines[1], "Negative prompt: ")
		} else {
			info.Prompt = lines[0]
		}
	}
	if len(lines) == 3 {
		info.Prompt = lines[0]
		info.NegativePrompt = strings.TrimPrefix(lines[1], "Negative prompt: ")
	}
	stepLineParts := strings.Split(lines[len(lines)-1], ", ")
	for _, part := range stepLineParts {
		kv := strings.Split(part, ": ")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid input: expected key-value pair, got %s", part)
		}
		key, value := kv[0], kv[1]
		switch key {
		case "Steps":
			steps, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("invalid input: Steps should be an integer, got %s", value)
			}
			info.Steps = steps
		case "Sampler":
			info.Sampler = value
		case "CFG scale":
			scale, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid input: CFG scale should be a float, got %s", value)
			}
			info.CFGScale = scale
		case "Seed":
			seed, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid input: Seed should be an integer, got %s", value)
			}
			info.Seed = seed
		case "Size":
			info.Size = value
		case "Model hash":
			info.ModelHash = value
		case "Model":
			info.Model = value
		}
	}

	info.Raw = input
	return info, nil
}
