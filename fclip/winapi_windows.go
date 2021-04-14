//+build windows

package fclip

import "syscall"

type xint int32
type xchr int8
type _byte uint8
type dword uint32
type long int32
type word uint16
type hbitmap uintptr
type lpvoid uintptr
type handle uintptr
type lpbyte *_byte
type hdc uintptr
type ulargeInteger uint64
type ulong uint64

// POINT struct
type POINT struct {
	X xint
	Y xint
}

// DROPFILES struct
//
// typedef struct _DROPFILES {
// DWORD pFiles;
// POINT pt;
// BOOL  fNC;
// BOOL  fWide; } DROPFILES, *LPDROPFILES;
//
type DROPFILES struct {
	A xint
	B POINT
	C xint
	D xint
}

const (
	cfUnicodetext                 = 13
	cfHdrop                       = 15
	cfBitmap                      = 2
	gmemMoveable                  = 0x0002
	gmemZeroinit                  = 0x0040
	dvaspectContent         int32 = 1
	tymedHGlobal            int32 = 1
	tymedIStream            int32 = 4
	tymedIStorage           int32 = 8
	lptr                          = 0x0040
	biRgb                         = 0x0000
	biBitfields                   = 0x0003
	dibRgbColors                  = 0x00
	null                          = 0
	gmemFixed                     = 0x0000
	fileNameSize                  = 260
	fileGroupDescriptorName       = "FileGroupDescriptorW"
	fileContentsName              = "FileContents"
	notFoundErr                   = 2147745898
	successErrorStr               = "The operation completed successfully."
)

var (
	user32                   = syscall.MustLoadDLL("user32")
	openClipboard            = user32.MustFindProc("OpenClipboard")
	closeClipboard           = user32.MustFindProc("CloseClipboard")
	emptyClipboard           = user32.MustFindProc("EmptyClipboard")
	getClipboardData         = user32.MustFindProc("GetClipboardData")
	setClipboardData         = user32.MustFindProc("SetClipboardData")
	getDC                    = user32.MustFindProc("GetDC")
	registerClipboardFormatA = user32.MustFindProc("RegisterClipboardFormatA")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
	copyMemory   = kernel32.NewProc("RtlCopyMemory")
	getLastError = kernel32.NewProc("GetLastError")
	globalSize   = kernel32.NewProc("GlobalSize")

	shell32       = syscall.MustLoadDLL("shell32")
	dragQueryFile = shell32.MustFindProc("DragQueryFileW")

	gdi32     = syscall.MustLoadDLL("gdi32")
	getObject = gdi32.MustFindProc("GetObjectW")
	getDIBits = gdi32.MustFindProc("GetDIBits")

	ole32                        = syscall.MustLoadDLL("ole32")
	oleGetClipboard              = ole32.MustFindProc("OleGetClipboard")
	oleInitialize                = ole32.MustFindProc("OleInitialize")
	oleUninitialize              = ole32.MustFindProc("OleUninitialize")
	createILockBytesOnHGlobal    = ole32.MustFindProc("CreateILockBytesOnHGlobal")
	stgCreateDocfileOnILockBytes = ole32.MustFindProc("StgCreateDocfileOnILockBytes")
)

type bitmapInfo struct {
	bmiHeader bitmapInfoHeader
	bmiColors rgbQuad // [1]
}
type pBitmapInfo *bitmapInfo

type bitmapInfoHeader struct {
	biSize          dword
	biWidth         long
	biHeight        long
	biPlanes        word
	biBitCount      word
	biCompression   dword
	biSizeImage     dword
	biXPelsPerMeter long
	biYPelsPerMeter long
	biClrUsed       dword
	biClrImportant  dword
}

type pBitmapInfoHeader *bitmapInfoHeader

type bitmapFileHeader struct {
	bfType      uint16
	bfSize      uint32
	bfReserved1 uint16
	bfReserved2 uint16
	bfOffBits   uint32
}

type rgbQuad struct {
	rgbBlue     _byte
	rgbGreen    _byte
	rgbRed      _byte
	rgbReserved _byte
}

type bitmap struct {
	bmType       long
	bmWidth      long
	bmHeight     long
	bmWidthBytes long
	bmPlanes     word
	bmBitsPixel  word
	bmBits       lpvoid
}