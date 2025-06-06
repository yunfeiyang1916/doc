package config

import (
	"shop-v2/internal/share/options"
	cliflag "shop-v2/pkg/common/cli/flag"
	"shop-v2/pkg/log"
)

type Config struct {
	Log       *log.Options              `json:"log" mapstructure:"log"`
	Server    *options.ServerOptions    `json:"server" mapstructure:"server"`
	Registry  *options.RegistryOptions  `json:"registry" mapstructure:"registry"`
	Telemetry *options.TelemetryOptions `json:"telemetry" mapstructure:"telemetry"`
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Telemetry.AddFlags(fss.FlagSet("telemetry"))
	return fss
}

func (c *Config) Validate() []error {
	var errs []error
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	errs = append(errs, c.Telemetry.Validate()...)
	return errs
}

func New() *Config {
	return &Config{
		Log:       log.NewOptions(),
		Server:    options.NewServerOptions(),
		Registry:  options.NewRegistryOptions(),
		Telemetry: options.NewTelemetryOptions(),
	}
}
