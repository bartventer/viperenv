# viperenv

[![Go Reference](https://pkg.go.dev/badge/github.com/bartventer/viperenv.svg)](https://pkg.go.dev/github.com/bartventer/viperenv)
[![Go Report Card](https://goreportcard.com/badge/github.com/bartventer/viperenv)](https://goreportcard.com/report/github.com/bartventer/viperenv)

`viperenv` is a Go package that provides a simple way to bind environment variables to a struct using [Viper](https://github.com/spf13/viper). It allows you to define a struct with tags that specify the environment variable names. You can then use the `viperenv.Bind()` function to bind the environment variables to the struct. The upside of this is that you don't have to manually set each environment variable with `viper.BindEnv()`, you can just bind them all at once.

## Installation

To install `viperenv`, use `go get`:

```sh
go get github.com/bartventer/viperenv
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/spf13/viper"
    "github.com/bartventer/viperenv"
)

type Config struct {
    Host string `env:"HOST"`
    Port int    `env:"PORT"`
}

func main() {
    v := viper.New()
    var config Config
    # ... set up viper, and read config file, etc.
    err := viperenv.Bind(&config, v, viperenv.BindOptions{
        AutoEnv: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(config.Host)
    fmt.Println(config.Port)
}

```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
