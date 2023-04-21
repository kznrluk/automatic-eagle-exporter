package export

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/kznrluk/automatic-eagle-exporter/eagle"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func CreateZip(exportImages []*ExportImage) error {
	tempDir, err := os.MkdirTemp("", "zip")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	pack := eagle.Pack{}
	for _, exportImage := range exportImages {
		dirPath := filepath.Join(tempDir, fmt.Sprintf("%s.info", exportImage.EagleMeta.ID))
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}

		pngPath := filepath.Join(dirPath, exportImage.RealName)
		err = ioutil.WriteFile(pngPath, exportImage.PNGData, 0644)
		if err != nil {
			return err
		}

		metaPath := filepath.Join(dirPath, "metadata.json")
		metaData, err := json.Marshal(exportImage.EagleMeta)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(metaPath, metaData, 0644)
		if err != nil {
			return err
		}

		pack.Images = append(pack.Images, exportImage.EagleMeta)
	}

	packPath := filepath.Join(tempDir, "pack.json")
	packData, err := json.Marshal(pack)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(packPath, packData, 0644)
	if err != nil {
		return err
	}

	zipName := fmt.Sprintf("%s.eaglepack", time.Now().Format("20060102150405"))

	return compressFiles(tempDir, zipName)
}

func compressFiles(folderPath string, outputZip string) error {
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %v", err)
		}

		relPath, _ := filepath.Rel(folderPath, path)
		header.Name = filepath.ToSlash(relPath) // Use forward slash as path separator
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %v", err)
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to write file to zip: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to compress files: %v", err)
	}

	return nil
}
