package trace

const TraceName = "mxshop"

// 这个跟业务里的option一样，是因为 如果以后想改这个trace会直接影响到业务 所以业务得隔离开配置
type Options struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"`
	Batcher  string  `json:"batcher"`
}
