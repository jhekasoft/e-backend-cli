# e-backend-cli

```
░█▀▀░█▀▄░█▀█░█▀▀░█░█░█▀▀░█▀█░█▀▄
░█▀▀░█▀▄░█▀█░█░░░█▀▄░█▀▀░█░█░█░█
░▀▀▀░▀▀░░▀░▀░▀▀▀░▀░▀░▀▀▀░▀░▀░▀▀CLI
```
CLI for generating golang backend applications with [e-backend](https://github.com/jhekasoft/e-backend).

[![Go Report Card](https://goreportcard.com/badge/github.com/jhekasoft/e-backend-cli)](https://goreportcard.com/report/github.com/jhekasoft/e-backend-cli)

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
