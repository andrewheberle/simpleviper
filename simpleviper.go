// Package simpleviper a convenience wrapper around [viper] to avoid repeated boilerplate code in my own projects when integrating [viper] with [cobra] or [github.com/bep/simplecobra] and [pflag].
package simpleviper

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Errors returned by Init
var (
	ErrInvalidFlagset = errors.New("invalid flagset")
)

// A Viperlet is used to bind flags with env vars based on the options provided to New.
//
// Although it is safe to use an unitialised Viperlet, it is equivalent to calling New without any options, so it's usefulness is limited.
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
// Creating a new Viperlet with no options is valid but it does not provide any specific features without manually using the underlying [*viper.Viper] instance via the Viper method.
func New(opts ...Option) *Viperlet {
	v := new(Viperlet)

	// set options
	for _, o := range opts {
		o(v)
	}

	return v
}

// Viper provides access to the underlying [*viper.Viper] instance
func (v *Viperlet) Viper() *viper.Viper {
	if v.viper == nil {
		v.viper = viper.New()
	}

	return v.viper
}

// Init binds the provided [*pflag.FlagSet] and env vars to the underlying [*viper.Viper] instance
func (v *Viperlet) Init(flagset *pflag.FlagSet) error {
	if flagset != nil {
		// bind *pflag.FlagSet to *viper.Viper instance
		if err := v.Viper().BindPFlags(flagset); err != nil {
			return err
		}
	}

	// bind to env
	if v.bindEnv {
		if v.envPrefix != "" {
			v.Viper().SetEnvPrefix(v.envPrefix)
		}

		if v.envKeyReplacer != nil {
			v.Viper().SetEnvKeyReplacer(v.envKeyReplacer)
		}

		v.Viper().AutomaticEnv()
	}

	// read in config if specified
	if v.configFile != "" {
		v.Viper().SetConfigFile(v.configFile)
		if err := v.Viper().ReadInConfig(); err != nil {
			// return all errors if allowMissingConfig is not true
			if !v.allowMissingConfig {
				return err
			}

			// otherwise only return error if it is NOT a viper.ConfigFileNotFoundError error
			if !errors.Is(err, viper.ConfigFileNotFoundError{}) && !errors.Is(err, os.ErrNotExist) {
				// error was something else so return it
				return err
			}
		}
	}

	// set any values from viper as flags once other steps are done
	flagset.VisitAll(func(f *pflag.Flag) {
		if v.Viper().IsSet(f.Name) && v.Viper().GetString(f.Name) != "" {
			flagset.Set(f.Name, v.Viper().GetString(f.Name))
		}
	})

	return nil
}

type Option func(*Viperlet)

// WithViper allows passing your own [*viper.Viper] instance
func WithViper(viper *viper.Viper) Option {
	return func(v *Viperlet) {
		v.viper = viper
	}
}

// WithEnv enables environment variable binding. See [viper.AutomaticEnv] for details.
func WithEnv() Option {
	return func(v *Viperlet) {
		v.bindEnv = true
	}
}

// WithEnvPrefix enables environment variable binding using the provided prefix. See [viper.SetEnvPrefix] for details.
func WithEnvPrefix(prefix string) Option {
	return func(v *Viperlet) {
		v.bindEnv = true
		v.envPrefix = prefix
	}
}

// WithEnvKeyReplacer uses the provided [*strings.Replacer] for environment variable names. See [viper.SetEnvKeyReplacer] for details.
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
