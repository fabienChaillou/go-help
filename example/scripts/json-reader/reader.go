package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Définir le flag pour le chemin du fichier JSON
	filePath := flag.String("file", "", "Chemin vers le fichier JSON à lire")
	flag.Parse()

	// Vérifier si un chemin de fichier a été spécifié
	if *filePath == "" {
		fmt.Println("Erreur: Veuillez spécifier un fichier JSON avec le flag -file")
		fmt.Println("Exemple: go run main.go -file data.json")
		os.Exit(1)
	}

	// Lire le contenu du fichier
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du fichier: %v\n", err)
		os.Exit(1)
	}

	// Créer une variable pour stocker les données JSON
	var jsonData interface{}

	// Décoder le JSON
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Printf("Erreur lors du décodage JSON: %v\n", err)
		os.Exit(1)
	}

	// Afficher les données JSON avec une indentation pour plus de lisibilité
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Printf("Erreur lors du formatage JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Contenu du fichier JSON:")
	fmt.Println(string(prettyJSON))
}
