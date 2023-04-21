package export

import (
	"fmt"
	"github.com/kznrluk/automatic-eagle-exporter/eagle"
	"github.com/kznrluk/automatic-eagle-exporter/pnginfo"
	"github.com/kznrluk/automatic-eagle-exporter/utils"
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

	image := &ExportImage{
		RealName:   realName,
		Parameters: parameters,
		PNGData:    pngData,

		EagleMeta: eagle.Image{
			ID:               utils.CreateULID(),
			Name:             name,
			Size:             fileInfo.Size(),
			BTime:            fileInfo.ModTime().Unix(),
			MTime:            fileInfo.ModTime().Unix(),
			Ext:              "png",
			Tags:             []string{},
			Annotation:       rawParams,
			ModificationTime: fileInfo.ModTime().Unix(),
			NoThumbnail:      true,
			Width:            width,
			Height:           height,
			Palettes:         []string{},
			LastModified:     fileInfo.ModTime().Unix(),
		},
	}

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
