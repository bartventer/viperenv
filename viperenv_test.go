package viperenv_test

import (
	"os"
	"testing"

	"github.com/bartventer/viperenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	DatabaseURL   string `env:"DATABASE_URL" mapstructure:"database_url"`
	APIKey        string `env:"API_KEY" mapstructure:"api_key"`
	RequiredField string `env:"REQUIRED_FIELD,required" mapstructure:"required_field"`
	DefaultField  string `env:"DEFAULT_FIELD,default=default-value" mapstructure:"default_field"`
	IgnoreKey     string `env:"-" mapstructure:",omitempty"`
	IgnoreKey2    string `env:"" mapstructure:",omitempty"`
	IgnoreKey3    string
	Nested        struct {
		NestedField string `env:"NESTED_FIELD" mapstructure:"nested_field"`
	} `mapstructure:",squash"`
}

func TestBind(t *testing.T) {

	c := Config{}

	// Create a new viper instance.
	v := viper.New()
	v.AutomaticEnv()

	// expect error
	err := viperenv.Bind(&c, v, viperenv.BindOptions{})
	assert.Error(t, err)

	// Set up the test environment.
	os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb")
	os.Setenv("API_KEY", "my-secret-key")
	os.Setenv("REQUIRED_FIELD", "required-field")
	os.Setenv("IGNORE_KEY", "ignore-key")
	os.Setenv("IGNORE_KEY2", "ignore-key2")
	os.Setenv("IGNORE_KEY3", "ignore-key3")
	os.Setenv("NESTED_FIELD", "nested-field")

	// Bind the env vars to the config.
	if err := viperenv.Bind(&c, v, viperenv.BindOptions{
		AutoEnv:       true,
		EnvPrefix:     "",
		AllowEmptyEnv: true,
	}); err != nil {
		t.Fatal(err)
	}

	// Unmarshal the config.
	err = v.Unmarshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that the env vars are bound to the config.
	assert.Equal(t, "postgres://user:password@localhost:5432/mydb", c.DatabaseURL)
	assert.Equal(t, "my-secret-key", c.APIKey)
	assert.Equal(t, "default-value", c.DefaultField)
	assert.Equal(t, "required-field", c.RequiredField)
	assert.Equal(t, "nested-field", c.Nested.NestedField)

	// Assert that the ignored env vars are not bound to the config.
	assert.Equal(t, "", c.IgnoreKey)
	assert.Equal(t, "", c.IgnoreKey2)
	assert.Equal(t, "", c.IgnoreKey3)

}
