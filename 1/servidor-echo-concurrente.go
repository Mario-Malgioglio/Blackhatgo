package main

import (
	"io"
	"log"
	"net"
)

// echo es una función manejadora que simplemente retransmite los datos recibidos
func echo(conn net.Conn) {
	defer conn.Close()

	// Crear un búfer para almacenar los datos recibidos
	b := make([]byte, 512)
	for {
		// Recibir datos a través de conn.Read
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Cliente desconectado")
			break
		}
		if err != nil {
			log.Println("Error inesperado")
			break
		}
		log.Printf("Recibidos %d bytes: %s\n", size, string(b))

		// Enviar datos de vuelta a través de conn.Write
		log.Println("Escribiendo datos")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatal("No se pudo escribir los datos")
		}
	}
}

func main() {
	// Enlazar al puerto TCP 20080 en todas las interfaces
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatal("No se pudo enlazar al puerto")
	}
	log.Println("Escuchando en 0.0.0.0:20080")

	for {
		// Esperar conexión. Crea una instancia de net.Conn al establecerse
		conn, err := listener.Accept()
		log.Println("Conexión recibida")
		if err != nil {
			log.Fatal("No se pudo aceptar la conexión")
		}
		// Manejar la conexión de forma concurrente con una goroutine
		go echo(conn)
	}
}
