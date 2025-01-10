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
