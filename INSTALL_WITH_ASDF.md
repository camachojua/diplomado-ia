# Instalación usando ASDF

ASDF es un manejador de versiones multiruntime. Funciona a través de plugins y hace fácil la gestión de versiones de entornos de ejecución a manera de `rvm`, `nvm`, `sdkman` y otros.

La página oficial es: [https://asdf-vm.com/](https://asdf-vm.com/).

## Dependencias

Para usar `asdf` es necesario contar con `git` y `curl`. En linux pueden usarse los gestores de paquetes propios de cada distribución. En MacOS puede usarse homebrew.

### MacOS

```bash
brew install coreutils curl git
```

### Derivados de Debian

```bash
apt-get install curl git
```

### Derivados de Arch

```bash
pacman -S curl git
```

## Instalación de ASDF

Para instalar `asdf` con el método recomendado por sus creadores, se clona el repositorio en un directorio oculto dentro del directorio personal:

```bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
```

Para dejarlo activo dentro de la linea de comandos debe agregarse una linea al archivo `~/.bash_profile` o `~/.bashrc`. En el caso de MacOS con versiones de Catalina o superiores, el shell por default es ZSH, por lo que debe agregarse al archivo `~/.zprofile` o `~/.zshrc`:

```bash
. "$HOME/.asdf/asdf.sh"
```

> **Nota**: hay que cerrar la terminal y volver a abrirla para que los cambios tengan efecto. Otra opción es ejecutar directamente `. "$HOME/.asdf/asdf.sh"`

## Instalación de plugins

`asdf` al ser un manejador de múltiples entornos de ejecución, se gestiona a través de plugins. La lista completa puede verse ejecutando:

```bash
asdf plugin list all
```

Para nuestro caso usaremos los correspondientes para Julia y Go:

```bash
asdf plugin add julia
asdf plugin add golang
```

## Instalación de los runtimes

Para visualizar las versiones disponibles de cada runtime, se ejecuta:

```bash
asdf list all golang
```
o

```bash
asdf list all julia
```

Instalaremos las últimas disponibles al momento de escribir este README:

```bash
asdf install golang 1.22.6
```

y

```bash
asdf install julia 1.10.4
```

Para luego listar las instalaciones disponibles:

```bash
asdf list
```

Entonces podemos establecer una versión por default para ser usada en todo el sistema:

```bash
asdf global julia 1.10.4
```

y

```bash
asdf global golang 1.22.6
```

Si se ejecuta nuevamente `asdf list`, se verán las versiones por default marcadas por un asterisco `*`.

---