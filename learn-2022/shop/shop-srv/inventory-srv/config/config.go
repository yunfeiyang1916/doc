package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
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
	Name        string       `mapstructure:"name" json:"name"`
	Host        string       `mapstructure:"host" json:"host"`
	Port        int          `mapstructure:"port" json:"port"`
	Tags        []string     `mapstructure:"tags" json:"tags"`
	MysqlInfo   MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	RedisInfo   RedisConfig  `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
	NacosConfig NacosConfig  `mapstructure:"nacos" json:"nacos"`
}
