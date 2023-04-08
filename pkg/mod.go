package pkg

//base mod packages
//封装了一些基础模块，目前有

//FBI WARNING:
//
//|-- errors -> exceptions
//|-- hasher -> hash algorithms, for example `crc32`
//|-- log -> a simple log implementation by go log package and builder design mode
//|-- nio -> the non-blocking io framework based on epoll and zero copy(linux) or kqueue(mac) by event driven,on windows, debugging only
//|-- pool -> some pools implementation
//|   |-- bytepool -> a non-copy bytes pool
//|   |-- gopool -> go coprocessors pool implement by chan
//|-- utils -> some tools
