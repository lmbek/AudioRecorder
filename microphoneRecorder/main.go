package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	winmm               = syscall.NewLazyDLL("winmm.dll")
	mciSendString       = winmm.NewProc("mciSendStringW")
	mciGetErrorString   = winmm.NewProc("mciGetErrorStringW")
	mciSendStringFormat = "%s %s %s %s"
)

func mciSendStringW(command string, returnData *string, returnLength int, hwndCallback int) int {
	var buf [256]uint16
	ret, _, _ := mciSendString.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(command))),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(hwndCallback),
	)
	*returnData = syscall.UTF16ToString(buf[:])
	return int(ret)
}

func main() {
	filePath := "output.wav"

	openCommand := fmt.Sprintf("open new Type waveaudio Alias capture")
	var result string
	err := mciSendStringW(openCommand, &result, 0, 0)
	if err != 0 {
		handleError(err, result)
		return
	}

	recordCommand := "record capture"
	err = mciSendStringW(recordCommand, &result, 0, 0)
	if err != 0 {
		handleError(err, result)
		return
	}

	fmt.Println("Recording... Press Enter to stop.")
	fmt.Scanln()

	stopCommand := "stop capture"
	err = mciSendStringW(stopCommand, &result, 0, 0)
	if err != 0 {
		handleError(err, result)
		return
	}

	saveCommand := fmt.Sprintf("save capture %s", filePath)
	err = mciSendStringW(saveCommand, &result, 0, 0)
	if err != 0 {
		handleError(err, result)
		return
	}

	closeCommand := "close capture"
	err = mciSendStringW(closeCommand, &result, 0, 0)
	if err != 0 {
		handleError(err, result)
		return
	}

	fmt.Printf("Audio recorded and saved to %s successfully!\n", filePath)
}

func handleError(err int, result string) {
	var errBuf [256]uint16
	mciGetErrorString.Call(uintptr(err), uintptr(unsafe.Pointer(&errBuf[0])), uintptr(len(errBuf)))
	fmt.Printf("Error: %s\nResult: %s\n", syscall.UTF16ToString(errBuf[:]), result)
}
