# Variables
BINARY_NAME=main
PKG=./...
DOT_FILE=tree.dot
PNG_FILE=tree.png

# Compilation
build:
	go build -o $(BINARY_NAME) .

# Exécution
run: build
	./$(BINARY_NAME)

# Tests
test:
	go test $(PKG)

# Tests avec couverture
cover:
	go test -cover $(PKG)

# Benchmarks
bench:
	go test -bench=. -benchmem $(PKG)

# Générer et convertir le fichier .dot (nécessite Graphviz)
treepng:
	go run .  # Assure-toi que ExportDOT() est dans main
	dot -Tpng $(DOT_FILE) -o $(PNG_FILE)
	@echo "Image générée : $(PNG_FILE)"

# Nettoyage
clean:
	go clean
	rm -f $(BINARY_NAME) $(DOT_FILE) $(PNG_FILE)

# Génère un graphe de benchmark avec Python
fibograph:
	python3 plot_benchmark.py
