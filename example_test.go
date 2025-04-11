package simpleviper_test

import (
	"fmt"
	"os"

	"github.com/andrewheberle/simpleviper"
	"github.com/spf13/pflag"
)

// This example demonstrates values coming from the command line, defaults, environment variables and a configuration file.
func ExampleViperlet_Init() {
	var example1, example2, example3, example4, example5, example6 string

	// create flagset, which in a real program (not an example) would use pflag.ExitOnError
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	fs.StringVar(&example1, "example1", "", "Example flag 1")
	fs.StringVar(&example2, "example2", "", "Example flag 2")
	fs.StringVar(&example3, "example3", "", "Example flag 3")
	fs.StringVar(&example4, "example4", "default will be overridden", "Example flag 4")
	fs.StringVar(&example5, "example5", "from default value", "Example flag 5")
	fs.StringVar(&example6, "example6", "", "Example flag 6")

	// as this is an example the command line options are provided
	fs.Parse([]string{
		"--example1", "from command line",
		"--example6", "will override config file",
	})

	// set some env vars
	os.Setenv("EXAMPLE1", "flag will take precedence")
	os.Setenv("EXAMPLE2", "from env var as flag is not set")
	os.Setenv("EXAMPLE3", "env var overrides config file")

	_ = simpleviper.New(simpleviper.WithEnv(), simpleviper.WithConfig("example.yml")).Init(fs)

	fmt.Println(example1)
	fmt.Println(example2)
	fmt.Println(example3)
	fmt.Println(example4)
	fmt.Println(example5)
	fmt.Println(example6)
	// Output:
	// from command line
	// from env var as flag is not set
	// env var overrides config file
	// from config file
	// from default value
	// will override config file
}

// This example demonstrates attempting to load a config file that is missing and is treated as a fatal error.
func ExampleViperlet_Init_missingConfigFile() {
	var exampleVar string

	// create flagset, which in a real program (not an example) would use pflag.ExitOnError
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	fs.StringVar(&exampleVar, "example", "", "Example flag")
	fs.Parse([]string{"--example", "from command line"})

	if err := simpleviper.New(simpleviper.WithConfig("missing.yml")).Init(fs); err != nil {
		fmt.Println("error: config file missing")

		// a real program (not an example) would exit with a non-zero exit code at this point, but in this case we return
		return
	}

	// this is not executed
	fmt.Println(exampleVar)
	// Output: error: config file missing
}

// This example demonstrates specifying a config file that is optional, so that a missing file is not an error.
func ExampleViperlet_Init_optionalConfigFile() {
	var exampleVar string

	// create flagset, which in a real program (not an example) would use pflag.ExitOnError
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	fs.StringVar(&exampleVar, "example", "", "Example flag")
	fs.Parse([]string{"--example", "from command line"})

	if err := simpleviper.New(simpleviper.WithOptionalConfig("missing.yml")).Init(fs); err != nil {
		fmt.Printf("error: %s\n", err)

		// a real program (not an example) would exit with a non-zero exit code at this point, but in this case we return
		return
	}

	fmt.Println(exampleVar)
	// Output: from command line
}
