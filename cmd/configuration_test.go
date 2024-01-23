package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGet_virtual_links(t *testing.T) {
	os.Setenv(env_links_path, "test")
	viper.BindEnv(env_links_path)
	if viper.Get(env_links_path) != "test" {
		t.Log("Failed to set VIRTUAL_LINKS_PATH")
		t.Fail()
	}
}

func Test_boolean_flags(t *testing.T) {
	os.Setenv(env_links, "true")
	viper.BindEnv(env_links)
	if !viper.GetBool(env_links) {
		t.Log("Failed to set VIRTUAL_LINKS")
		t.Fail()
	}
}
