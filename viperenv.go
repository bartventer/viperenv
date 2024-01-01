package viperenv

import (
	"fmt"
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
func Bind(outPtr interface{}, v *viper.Viper, options BindOptions) error {
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
	return bind(outPtr, v)
}

// bind binds the env vars to the config.
func bind(cfg interface{}, v *viper.Viper) error {
	val := reflect.ValueOf(cfg).Elem() // Get the value of the config struct.
	t := val.Type()                    // Get the type of the config struct.
	switch t.Kind() {                  // Switch on the kind of the config struct.
	case reflect.Ptr, reflect.Interface: // If the config struct is a pointer or interface, call the function recursively.
		return bind(val.Interface(), v)
	case reflect.Struct: // If the config struct is a struct, iterate over its fields.
		for i := 0; i < t.NumField(); i++ { // Iterate over the fields of the config struct.
			field := t.Field(i)                      // Get the field.
			if field.Type.Kind() == reflect.Struct { // If the field is a struct, call the function recursively.
				if err := bind(val.Field(i).Addr().Interface(), v); err != nil {
					return err
				}
			} else {
				// If the field is not a struct, get the env var and bind it to the config.
				// Get the env tag e.g. `env:"DATABASE_URL,required,default=default-value"`
				envTag := field.Tag.Get("env")
				// Split the env tag by comma.
				envTagSplit := strings.Split(envTag, ",")
				envVar := envTagSplit[0] // Get the env var e.g. `DATABASE_URL`
				required, defaultValue := false, ""
				if len(envTagSplit) > 1 {
					for _, tag := range envTagSplit[1:] {
						switch tag {
						case "required": // If the tag is `required`, set required to true.
							required = true
						default: // If the tag is not `required`, check if it starts with `default=`. If so, set the default value.
							if strings.HasPrefix(tag, "default=") {
								defaultValue = strings.TrimPrefix(tag, "default=")
							}
						}
					}
				}
				switch envVar {
				case "-", "": // If the env var is "-" or "", skip it.
					continue
				default: // If the env var is not "-" or "", bind it to the config.
					v.BindEnv(envVar)
					if required && v.GetString(envVar) == "" {
						return fmt.Errorf("required environment variable %s is not set", envVar)

					}
					if v.GetString(envVar) == "" && defaultValue != "" {
						v.Set(envVar, defaultValue)
					}
				}
			}
		}
	default: // If the config struct is not a struct, pointer, or interface, return.
		return nil
	}
	return nil
}
