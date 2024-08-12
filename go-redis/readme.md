基于
1、https://github.com/HDT3213/godis
2、[154-576-深入Go底层原理，重写Redis中间件实战（完结）]课程


set key value
*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n

get key
*2\r\n$3\r\nget\r\n$3\r\nkey\r\n


select 1
*2\r\n$6\r\nselect\r\n$1\r\n1\r\n