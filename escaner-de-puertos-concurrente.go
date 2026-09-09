package main

import (
	"fmt"
	"net"
	"sort"
)

// worker procesa los puertos recibidos y envía el resultado
func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	// 1. Crear el pool de trabajadores (100 goroutines)
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// 2. Enviar los puertos a escanear en una goroutine separada
	go func() {
		for i := 1; i <= 80; i++ {
			ports <- i
		}
	}()

	// 3. Recolectar los resultados de los primeros 80 puertos
	for i := 0; i < 80; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// 4. Cerrar canales y presentar los resultados ordenados
	close(ports)
	close(results)
	sort.Ints(openports)

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
