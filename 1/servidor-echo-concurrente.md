# Desarrollo de Servidores Echo Concurrentes en Go

Un **servidor Echo** en Go es un servicio de red básico que se utiliza habitualmente para aprender a leer y escribir datos en sockets, y cuya función principal es **retransmitir (*echo*) al cliente exactamente los mismos datos que este le envía** a través de la conexión.

---

## Componentes clave en Go

1. **Uso de `net.Conn` e interfaces I/O**: En Go, el tipo `net.Conn` representa una conexión de red TCP orientada a flujo. Puesto que `net.Conn` implementa los métodos `Read([]byte)` y `Write([]byte)`, actúa simultáneamente como un `io.Reader` y un `io.Writer` bidireccional.
2. **Escucha y aceptación de conexiones**: El servidor utiliza `net.Listen("tcp", ":20080")` para abrir un puerto de escucha. Luego, dentro de un bucle infinito, invoca `listener.Accept()` para aguardar la llegada de clientes.
3. **Manejo concurrente con goroutines**: Al aceptar una conexión, el servidor delega su procesamiento a una función manejadora mediante una goroutine (`go echo(conn)`). Esto permite que el hilo principal vuelva inmediatamente a `listener.Accept()` y siga atendiendo a nuevos clientes sin ser bloqueado.

---

## Enfoques para procesar los datos

Las fuentes presentan tres formas de implementar la función manejadora `echo(net.Conn)`:

* **Lectura y escritura explícita con buffer**: Se crea un slice de bytes (por ejemplo, `make([]byte, 512)`), se leen los datos recibidos mediante `conn.Read()` y se escriben de vuelta con `conn.Write()` hasta que el cliente se desconecta (`io.EOF`).
* **Lógica con I/O en búfer (`bufio`)**: Se envuelven los flujos con `bufio.NewReader(conn)` y `bufio.NewWriter(conn)`, lo que permite leer cadenas delimitadas (como saltos de línea con `ReadString('\n')`) y requiere llamar a `writer.Flush()` para transmitir el búfer.
* **Implementación idiomática con `io.Copy`**: Dado que `net.Conn` es tanto *Reader* como *Writer*, se puede pasar la misma variable `conn` como origen y destino a `io.Copy(conn, conn)`. Esta función copia continuamente todos los datos leídos de vuelta a la conexión de forma limpia.

---

## Código del servidor Echo en Go

A continuación se muestra la implementación del servidor Echo utilizando el manejo explícito con buffer:

```go
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
```

---

### Versión simplificada de la función `echo` con `io.Copy`

Si prefieres simplificar el manejador utilizando `io.Copy`, la función se reduce a:

```go
func echo(conn net.Conn) {
	defer conn.Close()
	// Copia datos de io.Reader a io.Writer a través de io.Copy
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatal("No se pudo leer/escribir datos")
	}
}
```
