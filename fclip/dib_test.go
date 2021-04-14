package fclip

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestDIBImage_At(t *testing.T) {
	err := waitOpenClipboard()
	assert.NoError(t, err)
	defer closeClipboard.Call()

	h, _, err := getClipboardData.Call(cfBitmap)
	if h == 0 {
		assert.Fail(t, "cannot get image from clipboard", err)
	}
	img, err := NewDibImage(h)
	assert.NotNil(t, img)
	assert.NoError(t, err)
	err = writeImageToFile(img, "test.jpg", "jpg")
	assert.NoError(t, err)

}


func writeImageToFile(img image.Image, name string, format string) error {
	outputFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	switch format {
	case "png":
		return png.Encode(outputFile, img)
	case "jpg":
		return jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 50})
	}
	return errors.New("unknown file format")

}