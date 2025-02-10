package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/config"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	Key1 string  `env:"KEY_1,notEmpty"`
	Key2 int     `env:"KEY_2"`
	Key3 *string `env:"KEY_3"`
}

func TestParseConfigWithEnvFile(t *testing.T) {
	t.Parallel()

	envFile := must.Return(os.CreateTemp(t.TempDir(), ".env"))
	t.Cleanup(func() {
		os.Remove(envFile.Name())
	})

	_ = must.Return(envFile.WriteString("KEY_1=key1\nKEY_2=2"))
	must.Do(envFile.Close())

	os.Setenv("ENV_FILE", envFile.Name())

	cfg := testConfig{}
	config.ParseConfig(&cfg)

	assert.Equal(t, "key1", cfg.Key1)
	assert.Equal(t, 2, cfg.Key2)
	assert.Nil(t, cfg.Key3)

	os.Unsetenv("KEY_1")
	os.Unsetenv("KEY_2")
}

func TestParseConfigWithEnvVars(t *testing.T) {
	t.Parallel()

	os.Setenv("KEY_1", "key1")
	os.Setenv("KEY_2", "2")

	cfg := testConfig{}
	config.ParseConfig(&cfg)

	assert.Equal(t, "key1", cfg.Key1)
	assert.Equal(t, 2, cfg.Key2)
	assert.Nil(t, cfg.Key3)

	os.Unsetenv("KEY_1")
	os.Unsetenv("KEY_2")
}

func TestParseConfigPanicsWhenRequiredFieldsMissing(t *testing.T) {
	t.Parallel()

	cfg := testConfig{}
	assert.Panics(t, func() {
		config.ParseConfig(&cfg)
	})
	fmt.Println(cfg.Key1)
}
