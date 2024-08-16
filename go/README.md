# Instalación

La instalación de `go` es sencilla, lo único que temenos que hacer es ir a la
[página de descarga](https://go.dev/doc/install) y seguir las instrucciones
listadas ahí (según el sistema operativo que tengas), se recomienda utilizar un
sistema linux para esto (si usas windows debes instalar WSL, las instrucciones
de instalación de esto se encuentran
[aquí](https://learn.microsoft.com/es-es/windows/wsl/install), para instalar go
necesitamos:

1. Descargar la versión de go correspondiente a nuestro [sistema operativo](https://go.dev/dl/)
2. Descomprimir el archivo descargado (ejemplo usando linux):
```bash
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.6.linux-amd64.tar.gz
```
3. Actualizar la variable `PATH` para que los ejecutables de `go` puedan ser
   reconocidos desde cualquier parte del sistema:
```bash
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile
```
4. Probar la instalación con `go version`

# Primer programa escrito en Go

Nuestro primer programa será un `hola mundo`, con el fin de aprender a cómo
generar un nuevo proyecto usando las herramientas que nos da el lenguaje.

## Preparación del proyecto

Antes de iniciar nuestro proyecto necesitamos un lugar para almacenar el código
que vámos a escribir, para ello crearemos una carpeta llamada `hello_go`:

```bash
mkdir -p ~/hello_go
cd ~/hello_go
```

Ahora necesitamos inicializar el proyecto con el comando:

```bash
 go mod init hello_go/hello
```

### Escribiendo el primer programa

Ahora necesitamos escribir el código, para ello necesitamos abrir nuestro editor
de texto (VSCode o Emacs) para crear un archivo llamado `hello.go` con el
siguiente contenido:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```
### Ejecutando el código

Para ver nuestro código en acción ejecutaremos el comando en la terminal:

```bash
cd ~/hello_go
go run .
```

Esto escribirá en la pantalla la cadena de texto `¡Hola Mundo!`

# Integración con VSCode

Para utilizar go con VSCode se recomienda [leer este artículo](https://code.visualstudio.com/docs/languages/go).
