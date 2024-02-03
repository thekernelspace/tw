package main

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// global config object
var globalConfig *Config = nil

type Config struct {
  ShowIcons bool `yaml:"show_icons"`
}

// default options
func NewConfig() *Config {
  return &Config{
    ShowIcons: true,
  }
}

func (c *Config) load() {
  // Load the config from the $HOME/.config/.tw.yaml or $HOME/.tw.yaml file
  home := os.Getenv("HOME")
  paths := []string{filepath.Join(home, "/.config/.tw.yaml"), filepath.Join(home, "/.tw.yaml")}

  // try all of the paths until one of them works
  foundpath := ""
  for _, path := range paths {
    if fi, err := os.Stat(path); err == nil {
      if fi.Mode().IsRegular() {
        foundpath = path
        break
      }
    }
  }

  // parse it and use the values
  if foundpath != "" {
    data, err := os.ReadFile(foundpath)
    if err != nil {
      log.Fatalf("Error reading config %s: %s\n", foundpath, err)
    }
    err = yaml.Unmarshal([]byte(data), &c)
    if err != nil {
      log.Fatalf("Error parsing config %s: %s\n", foundpath, err)
    }
  }
}

// load it to global
func LoadGlobalCfg() {
  globalConfig = NewConfig()
  globalConfig.load()
}

func GetGlobalCfg() *Config {
  return globalConfig
}
