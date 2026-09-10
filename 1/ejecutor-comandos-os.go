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
