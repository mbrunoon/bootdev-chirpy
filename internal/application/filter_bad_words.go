package application

import "strings"

func filterBadWords(text string) string {
	badwords := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}

	splitedText := strings.Split(text, " ")

	for index, word := range splitedText {
		if _, ok := badwords[strings.ToLower(word)]; ok {
			splitedText[index] = "****"
		}

	}

	return strings.Join(splitedText, " ")
}
