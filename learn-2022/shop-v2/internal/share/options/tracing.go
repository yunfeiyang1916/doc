package options

import (
	"shop-v2/pkg/errors"

	"github.com/spf13/pflag"
)

type TelemetryOptions struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "mxshop",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
}

func (t *TelemetryOptions) Validate() []error {
	errs := []error{}
	if t.Batcher != "jaeger" && t.Batcher != "zipkin" {
		errs = append(errs, errors.New("opentelemetry batcher only supports jaeger or zipkin"))
	}
	return errs
}

// AddFlags adds flags related to open telemetry storage for a specific tracing to the specified FlagSet.
func (t *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&t.Name, "telemetry.name", t.Name, "opentelemetry name")
	fs.StringVar(&t.Endpoint, "telemetry.endpoint", t.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&t.Sampler, "telemetry.sampler", t.Sampler, "opentelemetry sampler")
	fs.StringVar(&t.Batcher, "telemetry.batcher", t.Batcher, "opentelemetry batcher,only support jaeger and zipkin")
}
