package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/cobra"
)

type Config struct {
	KeyId   string
	AppKey  string
	Bucket  string
	Path    string
	Recurse bool
}

var k = koanf.New(".")

func LoadConfig(cmd *cobra.Command) Config {
	cfgFile, _ := cmd.Flags().GetString("config-file")
	if cfgFile != "" {
		if err := k.Load(file.Provider(cfgFile), toml.Parser()); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config file %s: %v\n", cfgFile, err)
			os.Exit(1)
		}
		localOverride := strings.Replace(cfgFile, ".toml", ".local.toml", 1)
		if _, err := os.Stat(localOverride); err == nil {
			if err := k.Load(file.Provider(localOverride), toml.Parser()); err != nil {
				log.Printf("Warning: could not load local override config: %v", err)
			}
		}
	}

	k.Load(env.Provider("B2_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "B2_"))
	}), nil)

	k.Load(posflag.Provider(cmd.Root().PersistentFlags(), ".", k), nil)
	k.Load(posflag.Provider(cmd.Flags(), ".", k), nil)

	return Config{
		KeyId:   k.String("key_id"),
		AppKey:  k.String("app_key"),
		Bucket:  k.String("bucket"),
		Path:    k.String("path"),
		Recurse: k.Bool("recurse"),
	}
}
