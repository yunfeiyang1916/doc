package config

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
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

type OssConfig struct {
	ApiKey      string `mapstructure:"key" json:"key"`
	ApiSecret   string `mapstructure:"secret" json:"secret"`
	Host        string `mapstructure:"host" json:"host"`
	CallBackUrl string `mapstructure:"callback_url" json:"callback_url"`
	UploadDir   string `mapstructure:"upload_dir" json:"upload_dir"`
}

type ServerConfig struct {
	Name        string       `mapstructure:"name" json:"name"`
	Host        string       `mapstructure:"host" json:"host"`
	Tags        []string     `mapstructure:"tags" json:"tags"`
	Port        int          `mapstructure:"port" json:"port"`
	JWTInfo     JWTConfig    `mapstructure:"jwt" json:"jwt"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
	NacosConfig NacosConfig  `mapstructure:"nacos" json:"nacos"`
	OssInfo     OssConfig    `mapstructure:"oss" json:"oss"`
}
