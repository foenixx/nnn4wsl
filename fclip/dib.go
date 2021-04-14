package fclip

import (
	"github.com/phuslu/log"
	"image"
	"image/color"
	"unsafe"
)

//DIBImage represents a windows DIB image
type DIBImage struct {
	data     []byte
	width 	int
	height int
	scanLine int
	palette  color.Palette
	pbi pBitmapInfo
	h uintptr
}

var _ image.Image = (*DIBImage)(nil)

//NewDibImage constructs new object
func NewDibImage(hDIB uintptr) (*DIBImage, error) {
	img := DIBImage{
		h: hDIB,
	}
	log.Info().Msg("creating DIBImage")

	var (
		pbih   pBitmapInfoHeader // bitmap info-header
		lpBits uintptr           // memory pointer
		hDC    hdc
		bi     = bitmapInfo{
			bmiHeader: bitmapInfoHeader{
				biSize: 40,
			},
		}
	)
	img.pbi = &bi

	r, _, err := getDC.Call(uintptr(0))
	if r == 0 {
		//m.Log.Error(context.Background(), "getDC", slog.F("result", r), slog.Error(err1))
		return nil, err
	}
	hDC = hdc(r)

	r, _, err = getDIBits.Call(uintptr(hDC), uintptr(hDIB), 0, 0, 0, uintptr(unsafe.Pointer(img.pbi)), uintptr(dibRgbColors))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 1", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return nil, err
	}
	img.width = int(img.pbi.bmiHeader.biWidth)
	img.height = int(img.pbi.bmiHeader.biHeight)

	// request uncompressed data without palette
	if img.pbi.bmiHeader.biCompression == biBitfields {
		img.pbi.bmiHeader.biCompression = biRgb
	}
	if img.pbi.bmiHeader.biBitCount != 32 {
		log.Error().Uint16("biBitCount", uint16(img.pbi.bmiHeader.biBitCount)).Msg("unsupported color depth")
		img.pbi.bmiHeader.biBitCount = 32
	}
	// 32bits = 4 bytes per pixel
	img.scanLine = 4 * img.width
	pbih = &img.pbi.bmiHeader

	r, _, err = globalAlloc.Call(gmemFixed, uintptr(pbih.biSizeImage))
	if r == 0 {
		//m.Log.Error(context.Background(), "globalAlloc 1", slog.F("result", r), slog.Error(err1))
		return nil, err
	}
	defer func() {
		if lpBits != 0 {
			// Free memory.
			globalFree.Call(lpBits)
		}
	}()
	lpBits = r
	// Retrieve the color table (RGBQUAD array) and the bits
	// (array of palette indices) from the DIB.
	r, _, err = getDIBits.Call(uintptr(hDC), uintptr(hDIB), 0, uintptr(pbih.biHeight), lpBits, uintptr(unsafe.Pointer(img.pbi)), uintptr(dibRgbColors))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 2", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return nil, err
	}

	img.data = make([]uint8, pbih.biSizeImage)

	r, _, err = copyMemory.Call(uintptr(unsafe.Pointer(&img.data[0])), lpBits, uintptr(pbih.biSizeImage))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 2", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return nil, err
	}
	//binary.Write(file, binary.LittleEndian, bytes)

	return &img, nil
}

// ColorModel returns the Image's color model.
func (i *DIBImage) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i *DIBImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.width,i.height)
}

// At returns the color of the pixel at (x, y).
func (i *DIBImage) At(x, y int) color.Color {
	offset := (i.height-y-1)*i.scanLine + x*4

	return &color.RGBA{B: i.data[offset], G: i.data[offset+1], R: i.data[offset+2], A: 255}
}

