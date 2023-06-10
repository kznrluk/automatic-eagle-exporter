package export

import (
	"fmt"
	"github.com/kznrluk/sdweb-eaglepack/eagle"
	"github.com/kznrluk/sdweb-eaglepack/pnginfo"
	"github.com/kznrluk/sdweb-eaglepack/utils"
	"os"
	"path/filepath"
)

type ExportImage struct {
	RealName   string
	PNGData    []byte
	Parameters *pnginfo.Parameters
	Upscale    bool
	Extras     *pnginfo.Extras
	EagleMeta  eagle.Image
}

func CreateExportImage(path string) (*ExportImage, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", path, err)
	}

	pngData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading PNG file %s: %v", path, err)
	}

	rawParams, err := pnginfo.FindValue(pngData, "parameters")
	if err != nil {
		return nil, fmt.Errorf("error finding png chunk with keyword parameters in file %s: %v", path, err)
	}

	parameters, err := pnginfo.ParseParameters(rawParams)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameters: %v", err)
	}

	realName := filepath.Base(path)
	ext := filepath.Ext(path)
	name := realName[:len(realName)-len(ext)]

	width, height, err := pnginfo.GetPNGDimensions(pngData)
	if err != nil {
		return nil, fmt.Errorf("error parsing dimensions: %v", err)
	}

	time := fileInfo.ModTime().Unix() * 1000

	image := &ExportImage{
		RealName:   realName,
		Parameters: parameters,
		PNGData:    pngData,

		EagleMeta: eagle.Image{
			ID:               utils.CreateULID(),
			Name:             name,
			Size:             fileInfo.Size(),
			BTime:            time,
			MTime:            time,
			Ext:              "png",
			Tags:             []string{},
			Annotation:       rawParams,
			ModificationTime: time,
			NoThumbnail:      true,
			Width:            width,
			Height:           height,
			Palettes:         []string{},
			LastModified:     time,
		},
	}

	image.EagleMeta.Tags = createTags(image)

	rawExtras, err := pnginfo.FindValue(pngData, "extras")
	if err != nil {
		image.Upscale = false
	} else {
		extras, err := pnginfo.ParseExtras(rawExtras)
		if err != nil {
			return nil, fmt.Errorf("error parsing parameters: %v", err)
		}

		image.Extras = extras
		image.EagleMeta.Tags = append(image.EagleMeta.Tags, "Upscale")
	}

	return image, nil
}

func createTags(image *ExportImage) []string {
	var tags []string

	tags = append(tags, fmt.Sprintf("Steps:%d", image.Parameters.Steps))
	tags = append(tags, fmt.Sprintf("Sampler:%s", image.Parameters.Sampler))
	tags = append(tags, fmt.Sprintf("CFGScale:%f", image.Parameters.CFGScale))
	// tags = append(tags, fmt.Sprintf("Seed:%d", image.Parameters.Seed))
	// tags = append(tags, fmt.Sprintf("Size:%s", image.Parameters.Size))
	tags = append(tags, fmt.Sprintf("ModelHash:%s", image.Parameters.ModelHash))
	tags = append(tags, fmt.Sprintf("Model:%s", image.Parameters.Model))
	tags = append(tags, fmt.Sprintf("Model:%s", image.Parameters.Model))

	for i, t := range image.Parameters.LoRAs {
		tags = append(tags, fmt.Sprintf("LoRA%d:%s", i, t))
	}

	return tags
}
