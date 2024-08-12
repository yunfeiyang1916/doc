package config

// grpc服务配置
type SrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

// 阿里短信apikey配置
type AliSmsConfig struct {
	ApiKey    string `mapstructure:"key" json:"key"`
	ApiSecret string `mapstructure:"secret" json:"secret"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	// 单位秒
	Expire int `mapstructure:"expire" json:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

// 阿里开源的分布式配置中心
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type ServerConfig struct {
	Name          string       `mapstructure:"name" json:"name"`
	Host          string       `mapstructure:"host" json:"host"`
	Tags          []string     `mapstructure:"tags" json:"tags"`
	Port          int          `mapstructure:"port" json:"port"`
	UserOpSrvInfo SrvConfig    `mapstructure:"userop_srv" json:"userop_srv"`
	GoodsSrvInfo  SrvConfig    `mapstructure:"goods_srv" json:"goods_srv"`
	JWTInfo       JWTConfig    `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo    AliSmsConfig `mapstructure:"sms" json:"sms"`
	RedisInfo     RedisConfig  `mapstructure:"redis" json:"redis"`
	ConsulInfo    ConsulConfig `mapstructure:"consul" json:"consul"`
	NacosConfig   NacosConfig  `mapstructure:"nacos" json:"nacos"`
}
