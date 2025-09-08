package main

func CheckGuess(secret, guess string) ([]string, int) {
	if len(secret) != len(guess) {
		return nil, 0 // ou panic / erreur
	}

	result := make([]string, len(secret))
	secretRunes := []rune(secret)
	guessRunes := []rune(guess)
	used := make([]bool, len(secret)) // pour éviter les doublons
	correctCount := 0

	// Première passe : lettres correctes et bien placées
	for i := range secretRunes {
		if guessRunes[i] == secretRunes[i] {
			result[i] = "correct"
			used[i] = true
			correctCount++
		}
	}

	// Deuxième passe : lettres correctes mais mal placées
	for i := range guessRunes {
		if result[i] == "" {
			for j := range secretRunes {
				if !used[j] && guessRunes[i] == secretRunes[j] {
					result[i] = "present"
					used[j] = true
					break
				}
			}
			if result[i] == "" {
				result[i] = "absent"
			}
		}
	}

	return result, correctCount
}
