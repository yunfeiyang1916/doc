package options

import (
	"shop-v2/pkg/errors"

	"github.com/spf13/pflag"
)

type RegistryOptions struct {
	Address string `json:"address" mapstructure:"address,omitempty"`
	// 指明基于什么进行注册 如：nacos consul etcd等
	Scheme string `json:"scheme" mapstructure:"scheme,omitempty"`
}

func NewRegistryOptions() *RegistryOptions {
	return &RegistryOptions{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}
}
func (o *RegistryOptions) Validate() []error {
	errs := []error{}
	if o.Address == "" || o.Scheme == "" {
		errs = append(errs, errors.New("address or scheme is empty"))
	}
	return errs
}
func (o *RegistryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Address, "consul.address", o.Address, ""+
		"consul address, if left , default is 127.0.0.1:8500")
	fs.StringVar(&o.Scheme, "consul.scheme", o.Scheme, ""+
		"registry scheme, if left , default is http")

}
