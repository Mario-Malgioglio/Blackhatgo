# Escáner de Puertos Concurrente en Go

Para construir un escáner de puertos concurrente en Go, se requiere una evolución por etapas que va desde la conexión básica individual hasta un diseño eficiente y controlado mediante un **Worker Pool** y **comunicación multicanal**.

---

## 1. Verificación básica de puertos (`net.Dial`)

La piedra angular para comprobar si un puerto TCP está disponible es la función `net.Dial("tcp", address)` del paquete `net`.

* Si `err == nil`, la conexión fue exitosa, lo que significa que el puerto está **abierto** (se debe cerrar la conexión con `conn.Close()`).
* Si devuelve un error, el puerto está **cerrado o filtrado**.

Un escaneo secuencial utiliza un ciclo `for` para probar puerto por puerto (por ejemplo, del 1 al 1024), pero resulta extremadamente lento al procesar una conexión a la vez.

---

## 2. Primer intento de concurrencia y sincronización

Para acelerar el proceso se recurre a la concurrencia nativa de Go:

1. **Goroutines sin sincronización (demasiado rápido):** Si se envuelve cada llamada a `net.Dial` dentro de una goroutine (`go func()`), el bucle `for` termina casi de inmediato y la función `main` finaliza antes de que las conexiones completen su intercambio de paquetes en la red.
2. **Uso de `sync.WaitGroup`:** Se puede corregir la finalización prematura utilizando un `sync.WaitGroup`:
   * Se incrementa el contador con `wg.Add(1)` por cada puerto.
   * Se decrementa con `defer wg.Done()` dentro de la goroutine.
   * Se bloquea la ejecución principal con `wg.Wait()` hasta que todas las goroutines terminen.

> **Problema:** Lanzar miles de goroutines de forma simultánea puede saturar la red o sobrepasar los límites del sistema operativo, produciendo resultados inconsistentes o erróneos.

---

## 3. Solución eficiente: Worker Pool y comunicación multicanal

La arquitectura recomendada combina un **Worker Pool** (grupo limitado de trabajadores) con **canales con búfer** para gestionar las tareas y evitar sobrecargas.

### Componentes clave

* **Canal de tareas (`ports`):** Un canal con búfer (por ejemplo, de capacidad 100) que contiene los números de puerto a escanear.
* **Canal de resultados (`results`):** Un canal donde los trabajadores envían el puerto si está abierto, o `0` si está cerrado.
* **Trabajadores (`workers`):** Un número fijo de goroutines (ej. 100) que leen constantemente del canal `ports` y transmiten su hallazgo a `results`.
* **Goroutine emisora:** Carga de forma independiente los puertos en el canal `ports` para no bloquear el flujo principal.
* **Ordenamiento:** Como la concurrencia entrega resultados desordenados, el hilo principal recolecta los puertos abiertos en un *slice* y los ordena con `sort.Ints()` antes de mostrarlos.

---

## Código completo del escáner concurrente

```go
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
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// 3. Recolectar los resultados de los 1024 puertos
	for i := 0; i < 1024; i++ {
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
```
