#！/bin/bash

# 获取脚本的执行目录
BasePath=$(cd `dirname $0`; pwd)
ProtoPath=$(cd `dirname $0`; pwd)

## 向外输出
echo "[  工作目录  ]:  ${BasePath}"
echo "[ ProtoPath ]:  $ProtoPath"

# 进入工作目录
cd $BasePath

protoc  --proto_path=.  --go_out=plugins=grpc:.  *.proto

protoc  --letmegrpc_out=. *.proto
mv ./standard.letmegrpc.go  ./standard.test.go

echo "[  完成  ]"
