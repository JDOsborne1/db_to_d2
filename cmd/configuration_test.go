package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGet_virtual_links(t *testing.T) {
	os.Setenv("VIRTUAL_LINKS_PATH", "test")
	viper.BindEnv("VIRTUAL_LINKS_PATH", "VIRTUAL_LINKS_PATH")
	if viper.Get("VIRTUAL_LINKS_PATH") != "test" {
		t.Log("Failed to set VIRTUAL_LINKS_PATH")
		t.Fail()
	}
}

func Test_boolean_flags(t *testing.T) {
	os.Setenv("VIRTUAL_LINKS", "true")
	viper.BindEnv("VIRTUAL_LINKS", "VIRTUAL_LINKS")
	if !viper.GetBool("VIRTUAL_LINKS") {
		t.Log("Failed to set VIRTUAL_LINKS")
		t.Fail()
	}
}
