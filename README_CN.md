# goredis

goredis 是一个Go语言实现的且基于NIO的Redis服务器。
本项目仅致力于测试基于go实现的redis服务器，在海量连接，高并发的情况下的极致性能

关键需要实现的及功能：
- 基于NIO(epoll,kqueue),go协程多路复用
- 支持string

#性能测试:
mac:

linux:

# 目录结构
参考包下的具体注释 mod.go