package simpleviper

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Viperlet struct {
	viper *viper.Viper

	// options
	bindEnv            bool
	envPrefix          string
	envKeyReplacer     *strings.Replacer
	configFile         string
	allowMissingConfig bool
}

// New returns an initialised Viperlet instance. The behaviour of the returned *Viperlet can be alted by passing various Option's.
//
// Creating a new Viperlet with no options is valid but it does not provide any specific features without manually using the underlying *viper.Viper instance via the Viper method.
func New(opts ...Option) *Viperlet {
	v := new(Viperlet)

	// set options
	for _, o := range opts {
		o(v)
	}

	// set up new viper if not passed
	if v.viper == nil {
		v.viper = viper.New()
	}

	return v
}

// Viper provides access to the underlying *viper.Viper instance
func (v *Viperlet) Viper() *viper.Viper {
	return v.viper
}

// Init binds flags and env vars to the underlying *viper.Viper instance
func (v *Viperlet) Init(flagset *pflag.FlagSet) error {
	// bind flagset to viper instance
	if err := v.viper.BindPFlags(flagset); err != nil {
		return err
	}

	// bind to env
	if v.bindEnv {
		if v.envPrefix != "" {
			v.viper.SetEnvPrefix(v.envPrefix)
		}

		if v.envKeyReplacer != nil {
			v.viper.SetEnvKeyReplacer(v.envKeyReplacer)
		}

		v.viper.AutomaticEnv()
	}

	// read in config if specified
	if v.configFile != "" {
		v.viper.SetConfigFile(v.configFile)
		if err := v.viper.ReadInConfig(); err != nil {
			if v.allowMissingConfig {
				// check if error was not a viper.ConfigFileNotFoundError
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					// error was something else so return it
					return err
				}
			} else {
				// missing config not allowed
				return err
			}
		}
	}

	// set any values from viper as flags
	flagset.VisitAll(func(f *pflag.Flag) {
		if v.viper.IsSet(f.Name) && v.viper.GetString(f.Name) != "" {
			flagset.Set(f.Name, v.viper.GetString(f.Name))
		}
	})

	return nil
}

type Option func(*Viperlet)

// WithViper allows passing your own viper.Viper instance
func WithViper(viper *viper.Viper) Option {
	return func(v *Viperlet) {
		v.viper = viper
	}
}

// WithEnv enables env var binding. See viper.AutomaticEnv for details.
func WithEnv() Option {
	return func(v *Viperlet) {
		v.bindEnv = true
	}
}

// WithEnvPrefix enables env var binding using the provided prefix. See viper.SetEnvPrefix for details.
func WithEnvPrefix(prefix string) Option {
	return func(v *Viperlet) {
		v.bindEnv = true
		v.envPrefix = prefix
	}
}

// WithEnvKeyReplacer uses the provided replacer for env var names. See viper.SetEnvKeyReplacer for details.
func WithEnvKeyReplacer(replacer *strings.Replacer) Option {
	return func(v *Viperlet) {
		v.bindEnv = true
		v.envKeyReplacer = replacer
	}
}

// WithConfig enables the reading of the provided config file. All errors, including if the config file is missing are treated as a failure.
func WithConfig(config string) Option {
	return func(v *Viperlet) {
		v.configFile = config
		v.allowMissingConfig = false
	}
}

// WithOptionalConfig enables the reading of the provided config file however this differs from WithConfig as a missing config file is not fatal.
func WithOptionalConfig(config string) Option {
	return func(v *Viperlet) {
		v.configFile = config
		v.allowMissingConfig = true
	}
}
