package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

var compteur int
var mutex sync.Mutex

func incrementer(wg *sync.WaitGroup) {
	mutex.Lock() // verrouille le mutex
	compteur++
	mutex.Unlock() // libère le mutex
	wg.Done()
}

// peux lister les goroutines et obtenir leur stack trace
func dumpGoroutines() {
	buf := make([]byte, 1<<20) // 1 MB buffer
	stacklen := runtime.Stack(buf, true)
	os.Stdout.Write(buf[:stacklen])
}

func main() {
	var wg sync.WaitGroup
	dumpGoroutines()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go incrementer(&wg)
	}

	wg.Wait()
	fmt.Println("Valeur finale du compteur :", compteur)
}
