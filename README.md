# keylogger-go
Progetto composto da un server HTTP e un client per il monitoraggio e la gestione di eventi di sistema. Il sistema include funzionalità per rilevare dati sensibili, salvare eventi come cambi di finestra e pressioni di tasti in un file di log locale, e inviare tali log al server.

## Struttura del progetto

- **server2.go:** Server HTTP per gestire i file.

- **uploads/:** Cartella per i file che vengono inviati al server (creata automaticamente)

- **projectThree/:** Directory del client
  - **main.go:** file principale del client
    
  - **go.mod:** Gestione delle dipendenze del modulo
    
  - **go.sum:** Verifica delle dipendenze
    
  - **callhome/:**
     - **callhome.go:** Gestione del log e interazione con il server
  - **detection/:**
     - **detection.go:** Rilevamento di dati sensibili
  - **logger/:**
     - **logger.go:** Implementazione del keylogger e monitoraggio della clipboard
  - **winapi/:**
     - **winapi.go:** Wrapper per le API di Windows
  - **wintypes/**
     - **wintypes.go:** Definizioni di tipi e costanti di Windows

## Funzinalità principali
**SERVER (server2.go)**
- Gestione di upload e download di file tramite HTTP.
- Archiviazione dei file caricati in una directory dedicata (./uploads).
- Endpoint per elencare e scaricare file salvati.

**CLIENT (projectThree)**
- Monitoraggio di eventi di sistema (tasti premuti, cambi finestra, contenuto della clipboard).
- Salvataggio di log in una directory nascosta.
- Analisi dei log per rilevare dati sensibili (email, IBAN, password, ecc.).
- Invio automatico dei log contenenti dati sensibili al server.

## Prerequisiti
- **Sistema Operativo:** Windows (necessario per il client, che utilizza API specifiche di Windows).
- **Go:** Versione 1.18 o successiva. Puoi scaricarlo da golang.org.
- **Permessi di Amministratore:** Potrebbero essere richiesti per accedere ad alcune API di sistema.

## Esecuzione
1. **Avvio del Server:**
    - Spostati nella directory principale del progetto.
    - Avvia il server:
      ```bash
         go run server2.go
      ```
2. **Avvio del Client:**
   - Spostati nella directory del client:
      ```bash
         cd projectThree
      ```
   - Assicurati che le dipendenze siano installate correttamente:
      ```bash
         go mod tidy
      ```
   - Avvia il client:
      ```bash
         go run main.go
      ```
3. **Scaricamento file di log dal server:**
   
