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
   - Scarica il file di log inviato al server andando sul browser e cercando:
     ```bash
        http://localhost:8080/download/uploaded_log.txt
     ```

## Funzionamento generale
**1. REGISTRAZIONE DEGLI EVENTI:**
  - Il client monitora diversi eventi di sistema:
     - Cambiamenti della finestra attiva (titolo della finestra).
     - Tasti premuti sulla tastiera (keylogger).
     - Contenuto della clipboard.
  - Gli eventi catturati vengono salvati nel file log.txt, situato in una cartella nascosta che di default è C:\Users\Public\MyHiddenLogs

**2. CARTELLA NASCOSTA E FILE DI LOG:**
  - Se la cartella nascosta non esiste, il client la crea al momento dell'esecuzione (funzione InitLogDirectory).
  - Se il file log.txt non esiste, viene creato automaticamente (funzione InitLogFile).

**3. ANALISI DEL FILE log.txt:**
  - Ogni 20 secondi, il client legge il contenuto del file log.txt (funzione MonitorFileAndSend).
  - Utilizza espressioni regolari (nel file detection.go) per cercare dati sensibili come:
     - Email.
     - IBAN.
     - Numeri di telefono.
     - Password.
     - Codici fiscali.

**4. INVIO AL SERVER:**
  - Se vengono rilevati dati sensibili nel file log.txt, il client:
     - Segnala che sono stati trovati dati sensibili (si può togliere per fare meno "rumore").
     - Invia il file log.txt al server tramite una richiesta HTTP POST.

**5. SERVER:**
  - Il server riceve il file e lo salva nella directory uploads.
  - Il file viene rinominato con il prefisso uploaded_ (ad esempio: uploaded_log.txt).
