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
	if len(lines) < 1 {
		return nil, fmt.Errorf("invalid input: expected at least 1 line, got %d", len(lines))
	}
	info := &Parameters{}
	promptLines := []string{}
	negativePromptLines := []string{}
	// Separate prompt and negative prompt lines
	inNegativePrompt := false
	var stepLine string
	for _, line := range lines {
		if strings.HasPrefix(line, "Negative prompt: ") {
			inNegativePrompt = true
			negativePromptLines = append(negativePromptLines, strings.TrimPrefix(line, "Negative prompt: "))
		} else if inNegativePrompt {
			if strings.Contains(line, "Steps:") {
				stepLine = line
				break
			}
			negativePromptLines = append(negativePromptLines, line)
		} else {
			if strings.Contains(line, "Steps:") {
				stepLine = line
				break
			}
			promptLines = append(promptLines, line)
		}
	}
	info.Prompt = strings.Join(promptLines, "\n")
	info.NegativePrompt = strings.Join(negativePromptLines, "\n")
	stepLineParts := strings.Split(stepLine, ", ")
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
