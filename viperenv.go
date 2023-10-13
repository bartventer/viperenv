package viperenv

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// BindOptions represents various options for binding env vars to the config.
type BindOptions struct {
	AutoEnv        bool              // Whether to call viper.AutomaticEnv() or not.
	EnvPrefix      string            // The env prefix; calls viper.SetEnvPrefix() if not empty.
	AllowEmptyEnv  bool              // Whether to call viper.AllowEmptyEnv() or not.
	EnvKeyReplacer *strings.Replacer // The env key replacer; calls viper.SetEnvKeyReplacer() if not nil.
}

// Bind binds the env vars to the config.
//
// Parameters:
//   - outPtr: The pointer to the config struct.
//   - v: The viper instance.
//   - options: The options.
func Bind(outPtr interface{}, v *viper.Viper, options BindOptions) {
	// Viper defaults.
	if options.AutoEnv {
		v.AutomaticEnv()
	}
	if options.EnvPrefix != "" {
		v.SetEnvPrefix(options.EnvPrefix)
	}
	if options.AllowEmptyEnv {
		v.AllowEmptyEnv(true)
	}
	if options.EnvKeyReplacer != nil {
		v.SetEnvKeyReplacer(options.EnvKeyReplacer)
	}
	// Bind the env vars to the config.
	bind(outPtr, v)
}

// bind binds the env vars to the config.
func bind(cfg interface{}, v *viper.Viper) {
	val := reflect.ValueOf(cfg).Elem() // Get the value of the config struct.
	t := val.Type()                    // Get the type of the config struct.
	switch t.Kind() {                  // Switch on the kind of the config struct.
	case reflect.Ptr, reflect.Interface: // If the config struct is a pointer or interface, call the function recursively.
		bind(val.Interface(), v)
	case reflect.Struct: // If the config struct is a struct, iterate over its fields.
		for i := 0; i < t.NumField(); i++ { // Iterate over the fields of the config struct.
			field := t.Field(i)                      // Get the field.
			if field.Type.Kind() == reflect.Struct { // If the field is a struct, call the function recursively.
				bind(val.Field(i).Addr().Interface(), v)
			} else { // If the field is not a struct, get the env var and bind it to the config.
				switch envVar := field.Tag.Get("env"); envVar {
				case "-", "": // If the env var is "-" or "", skip it.
					continue
				default: // If the env var is not "-" or "", bind it to the config.
					v.BindEnv(field.Name, envVar)
				}
			}
		}
	default: // If the config struct is not a struct, pointer, or interface, return.
		return
	}
}
