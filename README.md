# go-thumbnail
#### 功能描述
1. 非失真缩略图压缩
2. 缓存机制
3. 在线http服务

#### 引用公共模块
1. [gin](https://github.com/gin-gonic/gin.git)
2. [protobuf](https://github.com/golang/protobuf.git)
3. [sse](https://github.com/manucorporat/sse.git)
4. [resize](https://github.com/nfnt/resize.git)
5. [cutter](https://github.com/oliamb/cutter.git)
6. [context](https://github.com/golang/net.git)
7. [validator](https://github.com/go-playground/validator.git)

#### 编译方法
因为外部模块有部分依赖地址更换，因此修改了引用模块的依赖，全部整理在该项目中，该项目代码是完整的，直接可以用的。
```bash
export GOPATH=$PWD
git clone https://github.com/yalay/go-thumbnail.git src
cd src
go build
```

#### 使用方法
```bash
http://127.0.0.1:6789/pure/22.jpg?s=100x100
```

#### demo演示
 - 原始图片  
![原始图片](https://github.com/yalay/go-thumbnail/blob/master/public/pure/22.jpg)
 - 缩略图  
![缩略图](https://github.com/yalay/go-thumbnail/blob/master/cache/%252FPure%252F22.jpg100x100)
