package directory

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// Walk calls your callback function for every file name inside the directory.
// It doesn't distinguish between real files and directories.
func Walk(directory string, callBack func(string)) {
	file, err := os.Open(directory)

	if err != nil {
		panic(err)
	}

	defer file.Close()
	var bufferOnStack [4096]byte
	buffer := bufferOnStack[:]

	for {
		n, err := syscall.ReadDirent(int(file.Fd()), buffer)

		if err != nil {
			panic(err)
		}

		if n <= 0 {
			break
		}

		readBuffer := buffer[:n]

		for len(readBuffer) > 0 {
			dirent := (*syscall.Dirent)(unsafe.Pointer(&readBuffer[0]))
			readBuffer = readBuffer[dirent.Reclen:]

			// Skip deleted files
			if dirent.Ino == 0 {
				continue
			}

			// Skip hidden files
			if dirent.Name[0] == '.' {
				continue
			}

			for i, c := range dirent.Name {
				if c != 0 {
					continue
				}

				sliceHeader := reflect.SliceHeader{
					Len:  i,
					Cap:  len(dirent.Name),
					Data: uintptr(unsafe.Pointer(&dirent.Name[0])),
				}

				nameBytes := *(*[]byte)(unsafe.Pointer(&sliceHeader))
				name := string(nameBytes)
				callBack(name)
				break
			}
		}
	}
}
