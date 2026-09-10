package main

import (
	"io"
	"log"
	"net"
)

func handle(src net.Conn) {
	// 1. Establecer conexión con el puerto/host de destino
	dst, err := net.Dial("tcp", "joescatcam.website:80")
	if err != nil {
		log.Fatal("Unable to connect to our unreachable host")
	}
	defer dst.Close()

	// 2. Ejecutar en una goroutine para evitar que io.Copy bloquee el flujo
	go func() {
		// Copiar los datos del cliente local hacia el host de destino
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatal(err)
		}
	}()

	// 3. Copiar la respuesta del host de destino de vuelta al cliente local
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Escuchar conexiones en el puerto local 80
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal("Unable to bind to port")
	}

	for {
		// Aceptar la conexión entrante del cliente
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Unable to accept connection")
		}

		// Manejar la conexión de forma concurrente
		go handle(conn)
	}
}
