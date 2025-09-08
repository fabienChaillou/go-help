package main

import (
	"strings"
	"testing"
)

func TestGetFileContent(t *testing.T) {
	content, err := GetFileContent("a.txt")
	if err != nil {
		t.Fatalf("Erreur de lecture de a.txt: %v", err)
	}
	if !strings.Contains(content, "fichier A") {
		t.Errorf("Le contenu ne contient pas 'fichier A' : %s", content)
	}
}

func TestListFiles(t *testing.T) {
	files, err := ListFiles()
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du répertoire assets: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("Il manque des fichiers dans assets/: %d trouvés", len(files))
	}
}
