package simpleviper_test

import (
	"fmt"
	"os"

	"github.com/andrewheberle/simpleviper"
	"github.com/spf13/pflag"
)

func ExampleViperlet_Init() {
	var example1, example2, example3, example4, example5 string

	// create flagset
	fs := pflag.NewFlagSet("example", pflag.ExitOnError)
	fs.StringVar(&example1, "example1", "", "Example flag 1")
	fs.StringVar(&example2, "example2", "", "Example flag 2")
	fs.StringVar(&example3, "example3", "", "Example flag 3")
	fs.StringVar(&example4, "example4", "default will be overridden", "Example flag 4")
	fs.StringVar(&example5, "example5", "from default value", "Example flag 5")

	// supply some args
	fs.Parse([]string{"--example1", "from command line"})

	// set some env vars
	os.Setenv("EXAMPLE1", "flag will take precedence")
	os.Setenv("EXAMPLE2", "from env var as flag is not set")
	os.Setenv("EXAMPLE3", "env var overrides config file")

	simpleviper.New(simpleviper.WithEnv(), simpleviper.WithConfig("example.yml")).Init(fs)

	fmt.Println(example1)
	fmt.Println(example2)
	fmt.Println(example3)
	fmt.Println(example4)
	fmt.Println(example5)
	// Output:
	// from command line
	// from env var as flag is not set
	// env var overrides config file
	// from config file
	// from default value
}
