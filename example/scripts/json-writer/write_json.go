package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Personne est une structure représentant les données que nous voulons stocker
type Personne struct {
	Nom       string `json:"nom"`
	Prenom    string `json:"prenom"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}

func main() {
	// Créer des données à écrire
	personnes := []Personne{
		{
			Nom:       "Dupont",
			Prenom:    "Jean",
			Age:       32,
			Email:     "jean.dupont@exemple.fr",
			Telephone: "01 23 45 67 89",
		},
		{
			Nom:       "Martin",
			Prenom:    "Sophie",
			Age:       28,
			Email:     "sophie.martin@exemple.fr",
			Telephone: "01 98 76 54 32",
		},
	}

	// Convertir les données en JSON
	donnees, err := json.MarshalIndent(personnes, "", "  ")
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON:", err)
		return
	}

	// Écrire les données dans un fichier
	err = os.WriteFile("personnes.json", donnees, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier:", err)
		return
	}

	fmt.Println("Le fichier personnes.json a été créé avec succès!")

	// Afficher le contenu du fichier JSON pour vérification
	fmt.Println("Contenu du fichier JSON:")
	fmt.Println(string(donnees))
}
