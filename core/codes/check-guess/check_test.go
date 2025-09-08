package main

import (
	"reflect"
	"testing"
)

func TestCheckGuess(t *testing.T) {
	secret := "pomme"
	guess := "moppe"

	expectedResult := []string{"present", "present", "correct", "correct", "absent"}
	expectedCorrectCount := 2

	result, correctCount := CheckGuess(secret, guess)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Résultat incorrect. Attendu %v, obtenu %v", expectedResult, result)
	}

	if correctCount != expectedCorrectCount {
		t.Errorf("Nombre de lettres correctes incorrect. Attendu %d, obtenu %d", expectedCorrectCount, correctCount)
	}
}
