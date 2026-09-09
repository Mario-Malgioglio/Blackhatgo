# Estrategias y Aplicaciones de Port Forwarding en Redes

El **port forwarding** (o reenvío de puertos) es una técnica en la que se emplea un sistema o equipo intermediario (*proxy* o *jump box*) para retransmitir o canalizar conexiones de red alrededor o a través de un cortafuegos (*firewall*).

## Principales Aplicaciones

* **Evasión de cortafuegos y controles de salida (*egress controls*)**: Permite establecer comunicación con servidores o puertos de destino cuyo acceso directo está bloqueado por las reglas de la red.
* **Acceso a redes segmentadas**: Permite redirigir el tráfico por medio de un servidor puente para acceder a equipos dentro de redes privadas segmentadas o a puertos configurados en interfaces restrictivas.

## Ejemplo de Funcionamiento

Si las reglas de un cortafuegos prohíben la conexión directa desde un equipo cliente hacia un destino no permitido (como `evil.com`), el cliente puede conectarse a un servidor intermediario autorizado por el firewall (como `stacktitan.com`). El servidor intermediario recibe la conexión y reenvía el tráfico hacia el destino final (`evil.com`), permitiendo así eludir las restricciones del cortafuegos.

# Implementación de un Proxy TCP y Port Forwarding en Go

A continuación se presenta el código completo y la explicación técnica para construir un **port forwarder** (proxy TCP) en Go [1].

## Código Fuente

```go
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
```

---

## Detalles Clave del Funcionamiento

* **Escucha e intermediación (`main`)**: El servidor escucha en el puerto local (por ejemplo, el puerto 80) mediante `net.Listen("tcp", ":80")` y acepta clientes entrantes. Cada conexión cliente (`src`) se pasa a la función `handle` dentro de una goroutine para que el proxy siga aceptando más peticiones en paralelo [2].
* **Conexión al servidor remoto (`net.Dial`)**: La función `handle` inicia un canal saliente hacia la máquina remota de destino mediante `net.Dial("tcp", "host_destino:puerto")` [2].
* **Transferencia bidireccional continua (`io.Copy`)**: Puesto que la función `io.Copy` es bloqueante hasta que la conexión se cierra, la retransmisión desde el cliente local hacia el servidor remoto (`io.Copy(dst, src)`) se envuelve en una **goroutine anónima**. Mientras tanto, la respuesta del servidor remoto hacia el cliente local (`io.Copy(src, dst)`) se ejecuta en el hilo de la función manejadora [2].
