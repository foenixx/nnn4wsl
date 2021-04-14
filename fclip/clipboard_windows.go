// +build windows

package fclip

import (
	"errors"
	"io/ioutil"
	"syscall"
	"time"
	"unsafe"
)

// PathsToClipboard func
func PathsToClipboard(paths ...string) error {
	const dropFilesLen = 20
	var dropFiles = DROPFILES{20, POINT{0, 0}, 0, 1}


	err := waitOpenClipboard()
	if err != nil {
		return err
	}
	defer closeClipboard.Call()

	r, _, err := emptyClipboard.Call(0)
	if r == 0 {
		return err
	}

	dropFilesPtr := unsafe.Pointer(&dropFiles)
/*	dropFilesPtrArr := *((*[20]byte)(dropFilesPtr))
	dataLength := uintptr(len(dropFilesPtrArr))
	data := make([]byte, dataLength)
	copy(data[0:], dropFilesPtrArr[:])*/

	pathsAsArr := make([][]uint16, len(paths))
	arrLen := 0
	for index, fileName := range paths {
		var dbytes []uint16
		dbytes, err = syscall.UTF16FromString(fileName)
		if err != nil {
			return err
		}
		pathsAsArr[index] = dbytes
		arrLen += len(dbytes)
	}
	arrLen++
	fileNamesData := make([]uint16, arrLen)
	offset := 0
	for _, fileName := range pathsAsArr {
		copy(fileNamesData[offset:], fileName)
		offset += len(fileName)
	}
	fileNamesData[offset] = 0

	//dataSize := uintptr(len(data) * int(unsafe.Sizeof(data[0])))
	fileNamesSize := uintptr(arrLen * int(unsafe.Sizeof(fileNamesData[0])))

	h, _, err := globalAlloc.Call(gmemMoveable|gmemZeroinit, dropFilesLen + fileNamesSize)
	if h == 0 {
		return err
	}
	defer func() {
		if h != 0 {
			globalFree.Call(h)
		}
	}()

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return err
	}

	/*r, _, err = copyMemory.Call(l, uintptr(unsafe.Pointer(&data[0])), dropFilesLen)*/
	r, _, err = copyMemory.Call(l, uintptr(dropFilesPtr), dropFilesLen)
	if r == 0 {
		return err
	}
	r, _, err = copyMemory.Call(l+dropFilesLen, uintptr(unsafe.Pointer(&fileNamesData[0])), fileNamesSize)
	if r == 0 {
		return err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		if err.(syscall.Errno) != 0 {
			return err
		}
	}

	r, _, err = setClipboardData.Call(cfHdrop, h)
	if r == 0 {
		return err
	}
	h = 0 // suppress deferred cleanup
	return nil
}

func GetBitmapFromClipboard() (string, error) {
	err := waitOpenClipboard()
	if err != nil {
		return "", err
	}
	defer closeClipboard.Call()

	h, _, err := getClipboardData.Call(cfBitmap)
	if h == 0 {
		return "", err
	}
	return readBitmap(h)
}

func GetPathsFromClipboard() ([]string, error) {
	err := waitOpenClipboard()
	if err != nil {
		return nil, err
	}

	defer closeClipboard.Call()

	h, _, err := getClipboardData.Call(cfHdrop)
	if h == 0 {
		if errors.Is(err, syscall.Errno(0)) {
			// there are no content with specified format in clipboard
			return nil, nil
		}
		// an actual error occurred
		return nil, err
	}
	var paths []string
	paths, err = readFileList(h)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func  readBitmap(h uintptr) (tmpFile string, err error) {

	file, err := ioutil.TempFile("", "clipboard_image")
	if err != nil {
		return
	}
	defer file.Close()

	err = getBitmap(hbitmap(h), file)
	if err != nil {
		return
	}

	return file.Name(), nil
}

func readFileList(h uintptr) (paths []string, err error) {
	//paths = make([]string, 0)
	var err1 error
	defer func() {
		r, _, err1 := globalUnlock.Call(h)

		if r == 0 && !errors.Is(err1, syscall.Errno(0)) {
			//clear it and return error
			paths = nil
			err = err1
		}
	}()

	l, _, err1 := globalLock.Call(h)
	if l == 0 {
		err = err1
		return
	}

	count, _, err1 := dragQueryFile.Call(h, 0xFFFFFFFF, 0, 0)
	if count == 0 {
		err = err1
		return
	}

	paths = make([]string, count)

	for i := 0; i < int(count); i++ {
		buffSize, _, err1 := dragQueryFile.Call(h, uintptr(i), 0, 0)
		if buffSize == 0 {
			err = err1
			paths = nil
			return
		}
		var buffer []uint16 = make([]uint16, buffSize)
		r, _, err1 := dragQueryFile.Call(h, uintptr(i), uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&buffSize)))
		if r == 0 {
			err = err1
			paths = nil
			return
		}
		paths[i] = syscall.UTF16ToString(buffer)
	}
	return
}

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	var err error
	for time.Now().Before(limit) {
		r, _, err = openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return err
}

