package winapi

import (
	"syscall"
	"unsafe"

	"projectThree/wintypes"

	"golang.org/x/sys/windows"
)

var (
	// Caricamento delle librerie Windows
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// Funzioni da user32.dll
	setWindowsHookEx         = user32.NewProc("SetWindowsHookExA")
	callNextHookEx           = user32.NewProc("CallNextHookEx")
	unhookWindowsHookEx      = user32.NewProc("UnhookWindowsHookEx")
	setWinEventHook          = user32.NewProc("SetWinEventHook")
	unhookWinEvent           = user32.NewProc("UnhookWinEvent")
	getMessage               = user32.NewProc("GetMessageW")
	translateMessage         = user32.NewProc("TranslateMessage")
	dispatchMessage          = user32.NewProc("DispatchMessage")
	getWindowTextLength      = user32.NewProc("GetWindowTextLengthW")
	getWindowText            = user32.NewProc("GetWindowTextW")
	getKeyboardState         = user32.NewProc("GetKeyboardState")
	attachThreadInput        = user32.NewProc("AttachThreadInput")
	toUnicodeEx              = user32.NewProc("ToUnicodeEx")
	getKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	loadKeyboardLayout       = user32.NewProc("LoadKeyboardLayoutW")
	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")

	// Funzioni da kernel32.dll
	getCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	getThreadId        = kernel32.NewProc("GetThreadId")
)

/*
LoadKeyboardLayout - Carica un layout di tastiera specifico.
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-loadkeyboardlayoutw
*/
func LoadKeyboardLayout(layout string, flags uint32) wintypes.HKL {
	ret, _, _ := loadKeyboardLayout.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(layout))),
		uintptr(flags),
	)
	return wintypes.HKL(ret)
}

/*
GetKeyboardLayout - Ottiene il layout corrente della tastiera
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardlayout
*/
func GetKeyboardLayout(threadID wintypes.DWORD) wintypes.HKL {
	ret, _, _ := getKeyboardLayout.Call(uintptr(threadID))
	return wintypes.HKL(ret)
}

/*
ToUnicodeEx - Converte il tasto premuto in carattere Unicode, supportando diversi layout.
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-tounicodeex
*/
func ToUnicodeEx(uVirtKey wintypes.DWORD, uScanCode wintypes.DWORD, lpKeyState *[256]byte, lpChar *uint16, cchBuff int, wFlags wintypes.DWORD, dwhkl wintypes.HKL) int {
	ret, _, _ := toUnicodeEx.Call(
		uintptr(uVirtKey),
		uintptr(uScanCode),
		uintptr(unsafe.Pointer(lpKeyState)),
		uintptr(unsafe.Pointer(lpChar)),
		uintptr(cchBuff),
		uintptr(wFlags),
		uintptr(dwhkl),
	)
	return int(ret)
}

/*
SetWindowsHookEx - Imposta un hook globale per intercettare eventi
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexa
*/
func SetWindowsHookEx(idHook int, lpfn wintypes.HOOKPROC, hMod wintypes.HINSTANCE, dwThreadId wintypes.DWORD) wintypes.HHOOK {
	ret, _, _ := setWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return wintypes.HHOOK(ret)
}

/*
CallNextHookEx - Chiama il prossimo hook nella catena
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-callnexthookex
*/
func CallNextHookEx(hhk wintypes.HHOOK, nCode int, wParam wintypes.WPARAM, lParam wintypes.LPARAM) wintypes.LRESULT {
	ret, _, _ := callNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return wintypes.LRESULT(ret)
}

/*
UnhookWindowsHookEx - Rimuove l'hook specificato
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
*/
func UnhookWindowsHookEx(hhk wintypes.HHOOK) bool {
	ret, _, _ := unhookWindowsHookEx.Call(uintptr(hhk))
	return ret != 0
}

/*
SetWinEventHook - Imposta un hook per eventi di sistema (es. cambio finestra attiva)
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwineventhook
*/
func SetWinEventHook(eventMin wintypes.DWORD, eventMax wintypes.DWORD, hmodWinEventProc wintypes.HMODULE, pfnWinEventProc wintypes.WINEVENTPROC, idProcess wintypes.DWORD, idThread wintypes.DWORD, dwFlags wintypes.DWORD) wintypes.HWINEVENTHOOK {
	ret, _, _ := setWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		uintptr(syscall.NewCallback(pfnWinEventProc)),
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags),
	)
	return wintypes.HWINEVENTHOOK(ret)
}

/*
UnhookWinEvent - Rimuove un hook per eventi di sistema
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwinevent
*/
func UnhookWinEvent(hWinEventHook wintypes.HWINEVENTHOOK) wintypes.BOOL {
	ret, _, _ := unhookWinEvent.Call(uintptr(hWinEventHook))
	return wintypes.BOOL(ret)
}

/*
GetKeyboardState - Ottiene lo stato attuale della tastiera
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeyboardstate
*/
func GetKeyboardState(lpKeyState *[256]byte) wintypes.BOOL {
	ret, _, _ := getKeyboardState.Call(uintptr(unsafe.Pointer(&(*lpKeyState)[0])))
	return wintypes.BOOL(ret)
}

/*
GetCurrentThreadId - Ottiene l'ID del thread corrente
https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getcurrentthreadid
*/
func GetCurrentThreadId() wintypes.DWORD {
	ret, _, _ := getCurrentThreadId.Call()
	return wintypes.DWORD(ret)
}

/*
AttachThreadInput - Collega un input thread per condividere lo stato della tastiera
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-attachthreadinput
*/
func AttachThreadInput(idAttach wintypes.DWORD, idAttachTo wintypes.DWORD, fAttach wintypes.BOOL) wintypes.BOOL {
	ret, _, _ := attachThreadInput.Call(
		uintptr(idAttach),
		uintptr(idAttachTo),
		uintptr(fAttach),
	)
	return wintypes.BOOL(ret)
}

/*
GetForegroundWindow - Restituisce l'handle della finestra attualmente in primo piano
https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
*/
func GetForegroundWindow() wintypes.HWND {
	ret, _, _ := getForegroundWindow.Call()
	return wintypes.HWND(ret)
}

/*
GetWindowText - Ottiene il titolo della finestra specificata
*/
func GetWindowText(hwnd wintypes.HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1
	buf := make([]uint16, textLen)
	getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen),
	)
	return syscall.UTF16ToString(buf)
}

/*
GetWindowTextLength - Restituisce la lunghezza del titolo della finestra
*/
func GetWindowTextLength(hwnd wintypes.HWND) int {
	ret, _, _ := getWindowTextLength.Call(uintptr(hwnd))
	return int(ret)
}

/*
TranslateMessage - Traduzione dei messaggi per DispatchMessage
*/
func TranslateMessage(msg *wintypes.MSG) wintypes.BOOL {
	ret, _, _ := translateMessage.Call(uintptr(unsafe.Pointer(msg)))
	return wintypes.BOOL(ret)
}

/*
DispatchMessage - Invia il messaggio alla finestra di destinazione
*/
func DispatchMessage(msg *wintypes.MSG) wintypes.LRESULT {
	ret, _, _ := dispatchMessage.Call(uintptr(unsafe.Pointer(msg)))
	return wintypes.LRESULT(ret)
}

/*
GetThreadId - Restituisce l'ID del thread associato a un handle di processo
*/
func GetThreadId(threadHandle wintypes.HANDLE) wintypes.DWORD {
	ret, _, _ := getThreadId.Call(uintptr(threadHandle))
	return wintypes.DWORD(ret)
}

/*
GetMessage - Ottiene un messaggio dalla coda dei messaggi
*/
func GetMessage(msg *wintypes.MSG, hwnd wintypes.HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := getMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax),
	)
	return int(ret)
}

func GetWindowThreadProcessId(hwnd wintypes.HWND) wintypes.DWORD {
	var threadID wintypes.DWORD
	getWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&threadID)))
	return threadID
}
