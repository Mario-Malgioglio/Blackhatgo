# Manual de Programación Go para Hacking y Ciberseguridad Ofensiva

**Black Hat Go: Go Programming for Hackers and Pentesters** (publicado en 2020 por No Starch Press) es una obra escrita por **Tom Steele, Chris Patten y Dan Kottmann**. El libro enseña a profesionales de la seguridad ofensiva, *pentesters* y desarrolladores a construir sus propias herramientas de *hacking* automatizadas y eficientes aprovechando la velocidad, concurrencia, portabilidad y simplicidad de **Go**. La filosofía de los autores prioriza la **funcionalidad práctica sobre la elegancia** del diseño de software.

---

## Estructura y Contenidos Principales

### 1. Fundamentos y Entorno de Go (Capítulo 1)
* **Entorno de desarrollo**: Configuración del espacio de trabajo (`GOROOT`, `GOPATH`), herramientas principales de la línea de comandos (`go build`, `go run`, `go get`, `go fmt`) y selección de IDEs.
* **Sintaxis e Idiomas**: Repaso de tipos primitivos, *slices*, *maps*, *structs*, interfaces, manejo de errores y primitivas de concurrencia básica (*goroutines* y *channels*).

---

### 2. Desarrollo de Herramientas de Red y Protocolos (Capítulos 2 al 6)
* **TCP, Escáneres y Proxies (Cap. 2)**: Creación de escáneres de puertos TCP concurrentes optimizados con *Worker Pools* y `sync.WaitGroup`, desarrollo de servidores Echo, proxies TCP para redirección de puertos y réplica de la funcionalidad de ejecución remota de comandos de Netcat.
* **Clientes HTTP e Interacción Remota (Cap. 3)**: Peticiones HTTP, análisis de respuestas estructuradas en JSON/XML, integración con las APIs de **Shodan** y **Metasploit**, y extracción de metadatos de documentos analizando resultados web de Bing.
* **Servidores HTTP, Rutas y Middleware (Cap. 4)**: Construcción de servidores y enrutadores con `gorilla/mux` y `negroni`, recolectores de credenciales, *keyloggers* con **WebSockets** y multiplexación de conexiones C2 mediante proxies inversos.
* **Explotación de DNS (Cap. 5)**: Enumeración concurrente de subdominios, creación de un servidor y proxy DNS personalizado, y uso de **túneles DNS** para canales de Comando y Control (C2) en redes restrictivas.
* **Protocolos SMB y NTLM (Cap. 6)**: Implementación de comunicaciones binarias con SMB/NTLMSSP, *marshaling* y *unmarshaling* personalizado mediante reflexión y etiquetas de *struct*, ataques de adivinación de contraseñas, técnica de ***Pass-the-Hash*** y recuperación de hashes NTLMv2.

---

### 3. Extracción de Datos y Análisis de Paquetes (Capítulos 7 y 8)
* **Bases de Datos y Sistemas de Archivos (Cap. 7)**: Conexión y consulta a bases SQL (MySQL, PostgreSQL, MSSQL) y NoSQL (MongoDB), desarrollo de un minador de esquemas basado en expresiones regulares y rastreo recursivo del sistema de archivos.
* **Procesamiento de Paquetes Raw (Cap. 8)**: Captura y filtrado con `gopacket` y `libpcap` usando BPF, captura de credenciales en texto plano (FTP) y escaneo de puertos capaz de evadir protecciones **SYN-flood / SYN-cookies**.

---

### 4. Explotación Avanzada, Criptografía y Persistencia (Capítulos 9 al 14)
* **Fuzzing y Portabilidad de Exploits (Cap. 9)**: Creación de *fuzzers* para *buffer overflow* e inyección SQL, portabilidad de *exploits* desde Python y C a Go (ejemplo de elevación de privilegios con *Dirty COW*) y generación/transformación de *shellcode*.
* **Plugins y Herramientas Extensibles (Cap. 10)**: Arquitecturas modulares utilizando el sistema nativo de plugins de Go (`.so`/DLL) y la integración del motor de lenguaje **Lua** mediante `gopher-lua`.
* **Criptografía (Cap. 11)**: Cifrado simétrico (AES, RC2) y asimétrico (RSA), algoritmos de hashing (`bcrypt`, MD5, SHA-256), autenticación de mensajes (HMAC), autenticación mutua mediante certificados TLS y fuerza bruta a claves RC2.
* **Interacción con Windows API y Archivos PE (Cap. 12)**: Inyección de procesos con `syscall`, manipulación de memoria con `unsafe.Pointer` y `uintptr`, desarrollo de un analizador de la estructura de archivos ejecutables (PE) e interoperabilidad con C mediante **CGO**.
* **Estenografía (Cap. 13)**: Análisis de la estructura binaria e incrustación de *payloads* en imágenes PNG mediante la manipulación de *chunks* y cifrado XOR.
* **Construcción de un RAT de Comando y Control (Cap. 14)**: Desarrollo completo de una herramienta de acceso remoto (RAT) modular que incluye *implant* cliente, servidor e interfaz administrativa utilizando **gRPC** y **Protocol Buffers**.
