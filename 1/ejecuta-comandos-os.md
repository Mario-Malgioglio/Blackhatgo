# Ejecución de comandos del sistema operativo desde Go

Guía completa para ejecutar comandos del sistema operativo desde Go usando el paquete estándar `os/exec`, con especial atención a **Windows** y notas para **Linux/macOS**.

---

## Índice

1. [Introducción](#introducción)
2. [Requisitos](#requisitos)
3. [El paquete `os/exec`](#el-paquete-osexec)
4. [Diferencias entre sistemas operativos](#diferencias-entre-sistemas-operativos)
5. [¿Por qué `cmd /C` en Windows?](#por-qué-cmd-c-en-windows)
6. [Funcionamiento del código](#funcionamiento-del-código)
7. [Casos de uso comunes](#casos-de-uso-comunes)
8. [Errores frecuentes](#errores-frecuentes)
9. [Código completo](#código-completo)
10. [Compilación y ejecución](#compilación-y-ejecución)
11. [Resumen](#resumen)

---

## Introducción

Go permite lanzar procesos externos mediante el paquete estándar `os/exec`. La API es **multiplataforma**, pero el comportamiento del shell subyacente cambia según el sistema operativo:

- **Windows** → `cmd.exe` (o `powershell.exe`).
- **Linux** → `/bin/sh` (o el shell indicado).
- **macOS** → `/bin/sh` (o `zsh`).

Conocer estas diferencias evita errores como `executable file not found` o caracteres mal codificados.

---

## Requisitos

- Go 1.18 o superior.
- Sistema operativo: Windows 10/11, Linux o macOS.
- Terminal: `cmd`, PowerShell, Windows Terminal, bash o zsh.

---

## El paquete `os/exec`

El tipo principal es `exec.Cmd`, que representa un comando externo.

| Método / Campo | Descripción |
|---|---|
| `exec.Command(name, args...)` | Crea el comando (no lo ejecuta todavía). |
| `cmd.Output()` | Ejecuta y devuelve `stdout` como `[]byte`. |
| `cmd.CombinedOutput()` | Devuelve `stdout` + `stderr` combinados. |
| `cmd.Run()` | Ejecuta y espera; no captura salida salvo que se asigne. |
| `cmd.Start()` / `cmd.Wait()` | Ejecución asíncrona manual. |
| `cmd.Dir` | Directorio de trabajo del proceso hijo. |
| `cmd.Env` | Variables de entorno (slice `"CLAVE=valor"`). |
| `cmd.Stdout` / `cmd.Stderr` / `cmd.Stdin` | Destinos de los flujos estándar. |

---

## Diferencias entre sistemas operativos

| Aspecto | Windows | Linux / macOS |
|---|---|---|
| Shell por defecto | `cmd.exe` | `/bin/sh` |
| Flag para ejecutar y salir | `/C` | `-c` |
| Ejecutable del shell | `cmd` | `sh` |
| Built-ins del shell | `dir`, `echo`, `copy`, `cd`, `type` | `ls`, `echo`, `cp`, `cd`, `cat` |
| Separador de rutas | `\` | `/` |
| Codificación consola | `cp850`, `cp1252` (por defecto) | UTF-8 |
| Cambiar página de códigos | `chcp 65001` | — |
| PowerShell | `powershell -Command` | — (usar `pwsh` si está instalado) |

### Ejemplos equivalentes

```go
// Windows
exec.Command("cmd", "/C", "dir")

// Linux / macOS
exec.Command("sh", "-c", "ls -la")
```

### Detección multiplataforma con `runtime.GOOS`

```go
var cmd *exec.Cmd
if runtime.GOOS == "windows" {
	cmd = exec.Command("cmd", "/C", "dir")
} else {
	cmd = exec.Command("sh", "-c", "ls -la")
}
```

---

## ¿Por qué `cmd /C` en Windows?

En Windows, comandos como `dir`, `echo`, `copy`, `cd` o `type` **no son ejecutables**: son *built-ins* del intérprete `cmd.exe`. Si intentas:

```go
exec.Command("dir")
```

Obtendrás:

```
exec: "dir": executable file not found in %PATH%
```

La solución es invocar el intérprete:

```go
exec.Command("cmd", "/C", "dir")
```

| Flag | Comportamiento |
|---|---|
| `/C` | Ejecuta el comando y termina. |
| `/K` | Ejecuta el comando y mantiene la consola abierta. |

Si el comando sí es un `.exe` real (`ping`, `ipconfig`, `git`, `go`, `notepad`...), no necesitas `cmd`:

```go
exec.Command("ping", "8.8.8.8")
```

---

## Funcionamiento del código

El archivo `main.go` incluye **ocho funciones** que cubren los escenarios más habituales.

### 1. `ejecutarSimple`
Ejecuta `dir` mediante `cmd /C` y captura la salida con `cmd.Output()`.

### 2. `ejecutarExeDirecto`
Ejecuta `ping` sin pasar por el shell, porque es un `.exe` real.

### 3. `ejecutarConStreaming`
Conecta `cmd.Stdout = os.Stdout` para ver la salida en tiempo real.

### 4. `ejecutarConBuffer`
Separa `stdout` y `stderr` en dos `bytes.Buffer` para procesarlos por separado.

### 5. `ejecutarConEnv`
Añade una variable de entorno con `append(os.Environ(), "MI_VAR=...")`.

### 6. `ejecutarEnDirectorio`
Cambia el directorio de trabajo del proceso hijo con `cmd.Dir`.

### 7. `ejecutarPowerShell`
Invoca PowerShell directamente con `-Command`.

### 8. `ejecutarUTF8`
Fuerza `chcp 65001` para evitar problemas de codificación.

---

## Casos de uso comunes

| Situación | Comando |
|---|---|
| Comando interno del shell (Windows) | `exec.Command("cmd", "/C", "dir")` |
| Comando interno del shell (Linux/macOS) | `exec.Command("sh", "-c", "ls -la")` |
| Ejecutable en `PATH` | `exec.Command("ping", "8.8.8.8")` |
| Script PowerShell | `exec.Command("powershell", "-Command", "...")` |
| Script batch | `exec.Command("cmd", "/C", "script.bat")` |
| Script shell | `exec.Command("sh", "script.sh")` |
| Ruta con espacios | `exec.Command("cmd", "/C", `"C:\Program Files\app.exe"`)` |
| Timeout | `exec.CommandContext(ctx, "cmd", "/C", "...")` |

---

## Errores frecuentes

### 1. `executable file not found`

Intentaste ejecutar un built-in sin `cmd /C` (Windows) o sin `sh -c` (Linux/macOS).

```go
// ❌ Windows
exec.Command("dir")

// ✅ Windows
exec.Command("cmd", "/C", "dir")

// ❌ Linux/macOS
exec.Command("ls -la")

// ✅ Linux/macOS
exec.Command("sh", "-c", "ls -la")
```

### 2. Caracteres extraños (cp850 / cp1252)

La consola de Windows no usa UTF-8 por defecto.

```go
exec.Command("cmd", "/C", "chcp 65001 >nul && tu_comando")
```

O decodificando manualmente con `golang.org/x/text/encoding/charmap`:

```go
import "golang.org/x/text/encoding/charmap"

decoder := charmap.CodePage850.NewDecoder()
utf8Out, _ := decoder.Bytes(out)
fmt.Println(string(utf8Out))
```

### 3. Rutas con espacios

Encierra entre comillas dobles:

```go
exec.Command("cmd", "/C", `"C:\Program Files\App\app.exe"`)
```

### 4. El proceso no termina

Usa `/C` (no `/K`) y, si hace falta, un `context.WithTimeout`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
exec.CommandContext(ctx, "cmd", "/C", "ping -t 8.8.8.8").Run()
```

### 5. Variables de entorno perdidas

Si asignas `cmd.Env` sin incluir `os.Environ()`, el proceso hijo pierde el `PATH`:

```go
// ❌ Sin PATH
cmd.Env = []string{"MI_VAR=hola"}

// ✅ Conserva el entorno actual
cmd.Env = append(os.Environ(), "MI_VAR=hola")
```

---

## Código completo

```go
// main.go
//
// Ejemplo de ejecución de comandos del sistema operativo desde Go.
// Cubre: comandos del shell, ejecutables directos, streaming,
// buffers separados, variables de entorno, directorio de trabajo,
// PowerShell y salida UTF-8.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ejecutarSimple ejecuta "dir" usando cmd /C y captura la salida.
func ejecutarSimple() {
	fmt.Println("=== 1. Ejecución simple (dir) ===")

	cmd := exec.Command("cmd", "/C", "dir")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(out))
}

// ejecutarExeDirecto ejecuta un .exe sin pasar por cmd.
func ejecutarExeDirecto() {
	fmt.Println("=== 2. Ejecutable directo (ping) ===")

	cmd := exec.Command("ping", "-n", "2", "8.8.8.8")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(out))
}

// ejecutarConStreaming muestra la salida en tiempo real.
func ejecutarConStreaming() {
	fmt.Println("=== 3. Streaming en tiempo real ===")

	cmd := exec.Command("cmd", "/C", "ping -n 3 8.8.8.8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

// ejecutarConBuffer separa stdout y stderr en buffers distintos.
func ejecutarConBuffer() {
	fmt.Println("=== 4. Captura con buffers separados ===")

	cmd := exec.Command("cmd", "/C", "dir")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	fmt.Println("--- stdout ---")
	fmt.Println(stdout.String())

	if stderr.Len() > 0 {
		fmt.Println("--- stderr ---")
		fmt.Println(stderr.String())
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// ejecutarConEnv añade una variable de entorno al proceso hijo.
func ejecutarConEnv() {
	fmt.Println("=== 5. Variables de entorno ===")

	cmd := exec.Command("cmd", "/C", "echo %MI_VAR%")
	cmd.Env = append(os.Environ(), "MI_VAR=hola desde Go")

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Salida:", string(out))
}

// ejecutarEnDirectorio cambia el directorio de trabajo del proceso hijo.
func ejecutarEnDirectorio() {
	fmt.Println("=== 6. Directorio de trabajo ===")

	cmd := exec.Command("cmd", "/C", "cd")
	cmd.Dir = `C:\Users\Public`

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Directorio:", string(out))
}

// ejecutarPowerShell invoca PowerShell directamente.
func ejecutarPowerShell() {
	fmt.Println("=== 7. PowerShell ===")

	cmd := exec.Command("powershell", "-Command",
		"Get-Process | Select-Object -First 5 | Format-Table -AutoSize")

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(out))
}

// ejecutarUTF8 fuerza la página de códigos UTF-8 antes de ejecutar.
func ejecutarUTF8() {
	fmt.Println("=== 8. Salida UTF-8 (chcp 65001) ===")

	cmd := exec.Command("cmd", "/C", "chcp 65001 >nul && echo ñ á é í ó ú")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(out))
}

func main() {
	ejecutarSimple()
	ejecutarExeDirecto()
	ejecutarConStreaming()
	ejecutarConBuffer()
	ejecutarConEnv()
	ejecutarEnDirectorio()
	ejecutarPowerShell()
	ejecutarUTF8()
}
```

---

## Compilación y ejecución

### Windows

```powershell
go mod init ejemplo-exec
go run main.go

# Generar binario
go build -o ejemplo.exe main.go
.\ejemplo.exe
```

### Linux / macOS

```bash
go mod init ejemplo-exec
go run main.go

# Generar binario
go build -o ejemplo main.go
./ejemplo
```

### Compilación cruzada (desde cualquier SO)

```bash
# Binario para Windows
GOOS=windows GOARCH=amd64 go build -o ejemplo.exe main.go

# Binario para Linux
GOOS=linux GOARCH=amd64 go build -o ejemplo main.go

# Binario para macOS
GOOS=darwin GOARCH=arm64 go build -o ejemplo main.go
```

---

## Resumen

- Usa `exec.Command("cmd", "/C", "comando")` en Windows para built-ins del shell.
- Usa `exec.Command("sh", "-c", "comando")` en Linux/macOS para built-ins del shell.
- Usa `exec.Command("programa", "args...")` para ejecutables reales.
- Captura salida con `Output()`, `CombinedOutput()` o asignando `cmd.Stdout` / `cmd.Stderr`.
- Configura `cmd.Dir` y `cmd.Env` según necesites.
- Cuida la codificación (`chcp 65001`) y las rutas con espacios.
- Usa `runtime.GOOS` si necesitas código multiplataforma.

---

> **Archivo:** `Ejecucion-comandos-os.md`
>
> Guarda este contenido en un archivo `.md` para consultarlo o subirlo a un repositorio.