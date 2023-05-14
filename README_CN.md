# goredis

goredis 是一个Go语言实现的且基于NIO的Redis服务器。
本项目仅致力于测试基于go实现的redis服务器，在海量连接，高并发的情况下的极致性能

关键需要实现的及功能：
- 基于NIO(epoll,kqueue),go协程多路复用
- 支持string


#性能测试:(以下为本机模拟测试 redis 版本7.x)
- mac环境 :m1pro 8核 16g go1.20.4 macos 12.1
- goredis              redis 7.x
- ping 182481 请求/秒   200000 请求/秒
-  
- win:
- linux:

# 目录结构
参考包下的具体注释 mod.go