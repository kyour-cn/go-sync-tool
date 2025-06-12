//go:build windows

package mutex

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

var mutexHandle windows.Handle

func Create() {

	// 创建一个具有唯一名称的互斥体
	mutexName := "Global\\JxSyncUtilV3"
	mutex, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		fmt.Println("Error creating mutex:", err)
		os.Exit(1)
	}
	if errCode := syscall.GetLastError(); errors.Is(errCode, windows.ERROR_ALREADY_EXISTS) {
		fmt.Println("Another instance of the program is already running.")
		os.Exit(1)
	} else if errCode != nil {
		fmt.Println("Error checking for existing instance:", errCode)
		os.Exit(1)
	}

	mutexHandle = mutex

	// 在程序退出时释放互斥体资源
	defer func(handle windows.Handle) {
		err := windows.CloseHandle(handle)
		if err != nil {
			os.Exit(1)
		}
	}(mutex)
}

func Close() {
	err := windows.CloseHandle(mutexHandle)
	if err != nil {
		os.Exit(1)
	}
}
