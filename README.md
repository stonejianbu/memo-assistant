# memo-assistant

个人或公司小团队记忆助手，知识库数据保存在本地向量数据库weaviate
LLM模型基于豆包，需要配置apikey和开通接入模型语言模型（doubao-Seed-2.0-pro）和向量模型（doubao-embedding-vision）


## 安装weaviate环境
```shell
docker run -d \
  -p 8080:8080 \
  -p 50051:50051 \
  -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
  cr.weaviate.io/semitechnologies/weaviate:1.32.9
  
# 手动删除指定class的数据，假设class=TextData
curl -X DELETE http://localhost:8080/v1/schema/TextData

# 查看数据
curl http://localhost:8080/v1/schema
```

## 启动记忆助手
```shell
# 请先手动修改config/config.yaml配置doubao.apikey，然后再执行
make run
```

## 使用示例

### 记忆训练
```shell
curl -X POST http://127.0.0.1:9000/api/v1/train \
--header 'Content-Type: application/json' \
--data-raw '{
  "input": [
    "我是你的主人，以后叫我石头"
  ]
}'
```

### 搜索记忆
```shell
curl -X POST http://127.0.0.1:9000/api/v1/generate \
--header 'Content-Type: application/json' \
--data-raw '{
    "prompt": "谁是你的主人"
}'
```