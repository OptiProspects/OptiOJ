# OptiOJ

## JWT 密钥管理（狗屁通写的）

- JWT 密钥存储在项目根目录的 `jwtKey` 文件中
- 该文件不应提交到版本控制系统
- 部署时需要手动复制 key 文件到目标服务器
- 如需更换密钥，删除 key 文件后重启服务即可

```bash
protoc --go_out=. --go_opt=paths=source_relative `
>>     --go-grpc_out=. --go-grpc_opt=paths=source_relative `
>>     src/proto/judge_grpc_service/judge_grpc_service.proto
```
