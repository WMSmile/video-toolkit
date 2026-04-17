# 在 Mac 上开发时，正常拉取 Mac 包
go get github.com/k2-fsa/sherpa-onnx-go-macos

# 强制下载 Windows 包（即使报错也可以通过这种方式缓存到本地）
GOOS=windows go get github.com/k2-fsa/sherpa-onnx-go-windows


模型下载路径：
[asr-models](https://github.com/k2-fsa/sherpa-onnx/releases/tag/asr-models)

sherpa-onnx-paraformer-zh-2023-09-14.tar.bz2
这是阿里巴巴开源的工业级模型，也是目前 Sherpa-ONNX 社区最推荐的中文模型。


打包应用
```
wails3 build


brew install go-task

task package:universal


export GOPROXY=https://goproxy.cn,direct


```

