# e-backend

`e-backend` CLI.

![cat](./assets/android-chrome-192x192.png)

```
▗▄▄▄▖▗▄▄▖  ▗▄▖  ▗▄▄▖▗▖ ▗▖▗▄▄▄▖▗▖  ▗▖▗▄▄▄ 
▐▌   ▐▌ ▐▌▐▌ ▐▌▐▌   ▐▌▗▞▘▐▌   ▐▛▚▖▐▌▐▌  █
▐▛▀▀▘▐▛▀▚▖▐▛▀▜▌▐▌   ▐▛▚▖ ▐▛▀▀▘▐▌ ▝▜▌▐▌  █
▐▙▄▄▖▐▙▄▞▘▐▌ ▐▌▝▚▄▄▖▐▌ ▐▌▐▙▄▄▖▐▌  ▐▌▐▙▄▄▀
```

## Installation

```bash
go install github.com/jhekasoft/e-backend-cli@latest
```
## Application generation

```bash
e-backend-cli app create [name]
```
Where `name` is name of application is `lower-kebab-case`.

## Module generation

```bash
e-backend-cli module create [name] -t crud
```

Where `name` is name of module is `lowerCamelCase`, `-t` is template name
(simple, crud).

## Application template generation (for development purposes)

```bash
e-backend-cli app createTemplate -b ../e-backend-boilerplate-min -t ./generator/app/templates/simple -p e-backend-boilerplate
```
