grpc 连接池
grpc-pool-channel使用管道模式实现grpc连接池，单个连接同时只能处理一个请求，不能复用，并且一旦设置空闲连接不足，就容易长时间阻塞，所以不建议使用。
参考：https://github.com/processout/grpc-go-pool

参考文章：https://www.jb51.net/jiaoben/2958901tv.htm，文章中使用的库：https://github.com/shimingyah/pool
这个库性能毕竟高，可以使用
grpc单连接压测研究：https://xiaorui.cc/archives/6001，文章中使用的库：https://github.com/rfyiamcool/grpc-client-pool

