# simpleviper

[![Go Report Card](https://goreportcard.com/badge/github.com/andrewheberle/simpleviper)](https://goreportcard.com/report/github.com/andrewheberle/simpleviper)
[![GoDoc](https://godoc.org/github.com/andrewheberle/simpleviper?status.svg)](https://godoc.org/github.com/andrewheberle/simpleviper)
[![codecov](https://codecov.io/gh/andrewheberle/simpleviper/graph/badge.svg?token=JEFWB2U0GY)](https://codecov.io/gh/andrewheberle/simpleviper)

This module is a convenience wrapper around [viper](https://github.com/spf13/viper) to avoid repeated boilerplate code in my own projects when integrating Viper with [cobra](https://github.com/spf13/cobra) or [simplecobra](https://github.com/bep/simplecobra) and [pflag](https://github.com/spf13/pflag).

This provides a `Viperlet` type that is created using the `New` function, which by using the `Init` method will bind a `*pflag.Flagset` to the underlying `*viper.Viper` instance, while also optionally loading a configuration file and retrieving values from the environment.

The following example shows the integration with [simplecobra](https://github.com/bep/simplecobra) and allows the value of the `--stringflag` command line option to be set using the `STRINGFLAG` environment variable.

```go
import (
    "fmt"
    "os"

    "github.com/andrewheberle/simpleviper"
    "github.com/bep/simplecobra"
    "github.com/spf13/pflag"
)

type rootCommand struct {
    name string

    // flags
    stringFlag string

    commands []simplecobra.Commander
}

func (c *rootCommand) Name() string {
	return c.name
}

func (c *rootCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *rootCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "simpleviper example command"

    cmd.Flags().StringVar(&c.stringFlag, "stringflag", "", "Example string flag")

    return nil
}

func (c *rootCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	cmd := this.CobraCommand

    return viperlet.New(WithEnv()).Init(cmd.Flags())
}

func (c *rootCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
    fmt.Printf("string flag = %s\n", c.stringFlag)

    return nil
}

func main() {
    rootCmd := &rootCommand{
        name: "simpleviper-example",
        commands: []simplecobra.Commander{}
    }

    x, err := simplecobra.New(rootCmd)
	if err != nil {
		panic(err)
	}

    if _, err := x.Execute(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error during execution: %s\n", err)
		os.Exit(1)
	}
}
```

## Default Values, Flags, Environment and Configration

The precendece for a value is unchanged from [viper](https://github.com/spf13/viper), which is as follows where each item takes precedence over the item below it:

* flag
* env
* config
* default

## Using Viper Directly

The underlying `*viper.Viper` is exposed using the `Viper` method, so you are not restricted to just the features this module provides.
