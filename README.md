# go-thumbnail
#### 功能描述
1. 非失真缩略图压缩，支持jpg、png、gif、webp格式
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
8. [image](https://github.com/golang/image.git)

#### 编译方法 
**特别注意：因为解码webp依赖go1.6版本库，因此该项目必须是go1.6以上版本才可以编译。**

因为外部模块有部分依赖地址无效，因此修改为有效的依赖地址，并且全部整理在该项目中，该项目代码是完整的，不需要重新下载其他代码，直接可以用的。
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
