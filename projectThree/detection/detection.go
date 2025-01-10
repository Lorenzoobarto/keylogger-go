package detection

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Regex per rilevare dati sensibili
var (
	emailRegex    = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	ibanRegex     = regexp.MustCompile(`\b[A-Z]{2}[0-9]{2}[A-Z0-9]{11,30}\b`)
	phoneRegex    = regexp.MustCompile(`\b\+?[0-9]{1,4}[ -]?[0-9]{2,5}[ -]?[0-9]{4,}\b`)
	passwordRegex = regexp.MustCompile(`(?i)(password|pass|pwd)[:=]\s*\S+`)
	usernameRegex = regexp.MustCompile(`(?i)(username|user|usr)[:=]\s*\S+`)
	codiceFiscale = regexp.MustCompile(`\b[A-Z]{6}[0-9]{2}[A-Z][0-9]{2}[A-Z][0-9]{3}[A-Z]\b`)
)

// DetectSensitiveData analizza una stringa e rileva dati sensibili
func DetectSensitiveData(text string) map[string]string {
	results := make(map[string]string)

	if emailRegex.MatchString(text) {
		results["Email"] = emailRegex.FindString(text)
	}
	if ibanRegex.MatchString(text) {
		results["IBAN"] = ibanRegex.FindString(text)
	}
	if phoneRegex.MatchString(text) {
		results["Phone"] = phoneRegex.FindString(text)
	}
	if passwordRegex.MatchString(text) {
		results["Password"] = passwordRegex.FindString(text)
	}
	if usernameRegex.MatchString(text) {
		results["Username"] = usernameRegex.FindString(text)
	}
	if codiceFiscale.MatchString(text) {
		results["Codice Fiscale"] = codiceFiscale.FindString(text)
	}

	return results
}

// DetectSensitiveDataFromFile analizza ogni riga del file log.txt per rilevare dati sensibili
func DetectSensitiveDataFromFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Errore nell'apertura del file log.txt: %v\n", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	foundSensitive := false

	for scanner.Scan() {
		line := scanner.Text()
		results := DetectSensitiveData(line)

		// Stampa i risultati e aggiorna lo stato
		for key, value := range results {
			fmt.Printf("[ALERT] %s rilevato: %s\n", key, value)
			foundSensitive = true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Errore durante la lettura del file: %v\n", err)
	}

	return foundSensitive
}
