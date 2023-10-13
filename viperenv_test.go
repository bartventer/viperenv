package viperenv_test

import (
	"os"
	"testing"

	"github.com/bartventer/viperenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	APIKey      string `env:"API_KEY"`
	IgnoreKey   string `env:"-"`
	IgnoreKey2  string `env:""`
	IgnoreKey3  string
}

func TestBind(t *testing.T) {

	c := Config{}

	// Set up the test environment.
	os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb")
	os.Setenv("API_KEY", "my-secret-key")
	os.Setenv("IGNORE_KEY", "ignore-key")
	os.Setenv("IGNORE_KEY2", "ignore-key2")
	os.Setenv("IGNORE_KEY3", "ignore-key3")

	// Create a new viper instance.
	v := viper.New()
	v.AutomaticEnv()

	// Bind the env vars to the config.
	viperenv.Bind(&c, v, viperenv.BindOptions{
		EnvPrefix: "MY_APP",
	})

	// Unmarshal the config.
	err := v.Unmarshal(&c)
	if err != nil {
		panic(err)
	}

	// Assert that the env vars are bound to the config.
	assert.Equal(t, "postgres://user:password@localhost:5432/mydb", c.DatabaseURL)
	assert.Equal(t, "my-secret-key", c.APIKey)

	// Assert that the ignored env vars are not bound to the config.
	assert.Equal(t, "", c.IgnoreKey)
	assert.Equal(t, "", c.IgnoreKey2)
	assert.Equal(t, "", c.IgnoreKey3)
}
