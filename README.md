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

#### 使用方法
因为外部模块有较多依赖，因此全部整理在该项目中，该项目代码是完整的，直接可以用的。
```bash
git clone https://github.com/yalay/go-thumbnail.git ./
go build
```
