# memo-assistant

个人或小团队记忆助手，知识库数据保存在本地向量数据库 Weaviate。

LLM 模型基于豆包，需要在[火山引擎控制台](https://console.volcengine.com/ark)开通以下模型并获取 API Key：
- 语言模型：`doubao-Seed-2.0-pro`
- 向量模型：`doubao-embedding-vision`

---

## 部署方式

### 方式一：Docker Compose（推荐）

一键启动 memo-assistant 和 Weaviate：

```shell
# 设置豆包 API Key 环境变量
export DOUBAO_API_KEY=your_api_key_here

# 编译memo-assistant镜像
make package

# 启动所有服务
docker compose up -d
```

服务说明：
- **weaviate**：向量数据库，端口 8080（HTTP）和 50051（gRPC），数据持久化到 Docker volume
- **memo-assistant**：记忆助手服务，端口 9000，自动等待 Weaviate 健康后启动

停止服务：

```shell
docker compose down

# 同时删除数据卷（清空所有记忆数据）
docker compose down -v
```

---

### 方式二：手动部署

**1. 启动 Weaviate**

```shell
docker run -d \
  -p 8080:8080 \
  -p 50051:50051 \
  -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
  cr.weaviate.io/semitechnologies/weaviate:1.32.9
```

**2. 修改配置**

编辑 `config/config.yaml`，填入豆包 API Key：

```yaml
doubao:
  apiKey: your_api_key_here
```

**3. 启动服务**

```shell
make run
```

---

## 使用方式

### Web界面（推荐）

打开浏览器访问 `static/index.html`，通过可视化界面操作：

```shell
# 方式一：直接打开文件
open static/index.html

# 方式二：使用HTTP服务器（推荐，避免CORS问题）
python3 -m http.server 8000
# 然后访问 http://localhost:8000/static/index.html
```

**功能特性：**
- 📝 **训练记忆**：输入文本内容训练到知识库（支持 Ctrl+Enter 提交）
- 🔍 **搜索记忆**：输入问题搜索相关记忆（支持 Enter 搜索）
- 🎨 **Markdown显示**：搜索结果支持Markdown格式渲染
- 📦 **可折叠结果**：点击结果标题栏展开/收缩内容
- ⚙️ **API配置**：可自定义API地址

---

### API接口

#### 记忆训练

```shell
curl -X POST http://127.0.0.1:9000/api/v1/train \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "input": [
      "我是你的主人，以后叫我石头"
    ]
  }'
```

#### 搜索记忆

```shell
curl -X POST http://127.0.0.1:9000/api/v1/generate \
  --header 'Content-Type: application/json' \
  --data-raw '{
    "prompt": "谁是你的主人"
  }'
```

响应格式：

```json
{
  "code": 200,
  "data": "{\"answer\":\"根据记忆，你的主人是石头\"}"
}
```

---

## 数据管理

```shell
# 查看已有数据结构
curl http://localhost:8080/v1/schema

# 删除指定 class 的全部数据
curl -X DELETE http://localhost:8080/v1/schema/TextData
```
