package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Directory dove vengono salvati i file caricati
const uploadDir = "./uploads"

func main() {
	// Crea la cartella per salvare i file, se non esiste
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			fmt.Println("Errore nella creazione della directory:", err)
			return
		}
	}

	// Gestione degli endpoint
	http.HandleFunc("/upload", uploadHandler)          // Riceve i file
	http.HandleFunc("/files", listFilesHandler)        // Elenca i file salvati
	http.HandleFunc("/download/", downloadFileHandler) // Scarica un file specifico

	fmt.Println("Server2!!!! in ascolto su http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore nell'avvio del server:", err)
	}
}

/*
uploadHandler - Riceve un file tramite una richiesta HTTP POST multipart.
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	// Riceve il file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Errore nel caricamento del file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Salva il file con un prefisso "uploaded_"
	filename := filepath.Join(uploadDir, "uploaded_"+header.Filename)
	out, err := os.Create(filename)
	fmt.Println("Siamo entrati nella funzione UploadHandler")
	if err != nil {
		http.Error(w, "Errore nella creazione del file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Scrive il contenuto del file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Errore nella scrittura del file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File ricevuto con successo: %s\n", filename)
	fmt.Println("File salvato come:", filename)
}

/*
listFilesHandler - Elenca tutti i file presenti nella directory uploads.
*/
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Errore nella lettura della directory", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File disponibili per il download:")
	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(w, "http://%s/download/%s\n", r.Host, file.Name())
		}
	}
}

/*
downloadFileHandler - Permette di scaricare un file specifico.
*/
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Ottiene il nome del file dall'URL
	filename := filepath.Base(r.URL.Path)
	filePath := filepath.Join(uploadDir, filename)

	// Controlla se il file esiste
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File non trovato", http.StatusNotFound)
		return
	}

	// Serve il file per il download
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, filePath)
}
