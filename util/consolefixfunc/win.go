//go:build windows

package consolefixfunc

import "golang.org/x/sys/windows"

const (
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
)

//windows 下直接运行exe可能由于没开启VIRTUAL_TERMINAL导致控制台文本无法上色（ANSI控制字符被显示）
func EnableANSIConsole() error {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}
	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return err
	}
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	err = windows.SetConsoleMode(handle, mode)
	if err != nil {
		return err
	}
	return nil
}
