//go:build !windows

package mutex

func Create() {
	//walk只支持windows，请看notify_icon_windows.go
}

func Close() {
}
