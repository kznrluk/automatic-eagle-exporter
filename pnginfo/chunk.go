package pnginfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func FindValue(pngData []byte, keyword string) (string, error) {
	r := bytes.NewReader(pngData)
	// Check if the PNG header is correct
	header := make([]byte, 8)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return "", err
	}
	pngHeader := []byte{137, 80, 78, 71, 13, 10, 26, 10}
	if !bytes.Equal(header, pngHeader) {
		return "", fmt.Errorf("invalid PNG header")
	}
	for {
		chunkLengthBytes := make([]byte, 4)
		_, err := io.ReadFull(r, chunkLengthBytes)

		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
		chunkLength := binary.BigEndian.Uint32(chunkLengthBytes)
		chunkType := make([]byte, 4)
		_, err = io.ReadFull(r, chunkType)
		if err != nil {
			return "", err
		}
		if string(chunkType) == "tEXt" {
			chunkData := make([]byte, chunkLength)
			_, err = io.ReadFull(r, chunkData)
			if err != nil {
				return "", err
			}
			chunkReader := bytes.NewReader(chunkData)
			chunkKeyword := ""
			for {
				b, err := chunkReader.ReadByte()
				if err != nil {
					return "", err
				}
				if b == 0 {
					break
				}
				chunkKeyword += string(b)
			}
			if chunkKeyword == keyword {
				value := ""
				for {
					b, err := chunkReader.ReadByte()
					if err == io.EOF {
						break
					} else if err != nil {
						return "", err
					}
					value += string(b)
				}
				return value, nil
			}
		}
		_, err = r.Seek(int64(chunkLength)+4, io.SeekCurrent)
		if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("png chunk with keyword %s not found", keyword)
}

func GetPNGDimensions(pngData []byte) (width int, height int, err error) {
	if len(pngData) < 33 || string(pngData[:8]) != "\x89PNG\x0d\x0a\x1a\x0a" {
		return 0, 0, fmt.Errorf("not a PNG file")
	}

	if string(pngData[12:16]) != "IHDR" {
		return 0, 0, fmt.Errorf("IHDR chunk not found")
	}

	width = int(binary.BigEndian.Uint32(pngData[16:20]))
	height = int(binary.BigEndian.Uint32(pngData[20:24]))

	return width, height, nil
}
