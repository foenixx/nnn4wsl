// +build windows

package fclip

import (
	"encoding/binary"
	"os"
	"unsafe"
)

func getBitmap(hBMP hbitmap, file *os.File) error {
	var (
		pbih   pBitmapInfoHeader // bitmap info-header
		lpBits uintptr           // memory pointer
		hDC    hdc
		bi     = bitmapInfo{
			bmiHeader: bitmapInfoHeader{
				biSize: 40,
			},
		}
		pbi pBitmapInfo = &bi
	)

	r, _, err := getDC.Call(uintptr(0))
	if r == 0 {
		//m.Log.Error(context.Background(), "getDC", slog.F("result", r), slog.Error(err1))
		return err
	}
	hDC = hdc(r)

	r, _, err = getDIBits.Call(uintptr(hDC), uintptr(hBMP), 0, 0, 0, uintptr(unsafe.Pointer(pbi)), uintptr(dibRgbColors))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 1", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return err
	}

	// TODO: hack
	if pbi.bmiHeader.biCompression == biBitfields {
		pbi.bmiHeader.biCompression = biRgb
	}

	pbih = &pbi.bmiHeader

	r, _, err = globalAlloc.Call(gmemFixed, uintptr(pbih.biSizeImage))
	if r == 0 {
		//m.Log.Error(context.Background(), "globalAlloc 1", slog.F("result", r), slog.Error(err1))
		return err
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
	r, _, err = getDIBits.Call(uintptr(hDC), uintptr(hBMP), 0, uintptr(pbih.biHeight), lpBits, uintptr(unsafe.Pointer(pbi)), uintptr(dibRgbColors))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 2", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return err
	}

	var hdr = &bitmapFileHeader{}
	hdrSize := 14

	hdr.bfType = 0x4d42 // 0x42 = "B" 0x4d = "M"
	// Compute the size of the entire file.
	a := hdrSize
	b := unsafe.Sizeof(pbi.bmiColors)
	hdr.bfSize = (*(*uint32)(unsafe.Pointer(&a))) + uint32(pbih.biSize) + uint32(pbih.biClrUsed)*(*(*uint32)(unsafe.Pointer(&b))) + uint32(pbih.biSizeImage)
	hdr.bfReserved1 = 0
	hdr.bfReserved2 = 0

	// Compute the offset to the array of color indices.
	hdr.bfOffBits = uint32(uintptr(hdrSize) + uintptr(pbih.biSize) + uintptr(pbih.biClrUsed)*unsafe.Sizeof(pbi.bmiColors))

	err = binary.Write(file, binary.LittleEndian, hdr)
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, pbih)
	if err != nil {
		return err
	}

	bytes := make([]uint8, pbih.biSizeImage)

	r, _, err = copyMemory.Call(uintptr(unsafe.Pointer(&bytes[0])), lpBits, uintptr(pbih.biSizeImage))
	if r == 0 {
		//errnum, _, _ := getLastError.Call()
		//m.Log.Error(context.Background(), "getDIBits 2", slog.F("result", r), slog.F("errno", errnum), slog.Error(err1))
		return err
	}
	binary.Write(file, binary.LittleEndian, bytes)
	return err
}

