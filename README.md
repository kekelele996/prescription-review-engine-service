# 处方智能审核引擎（rxcheck）

为医院药房提供处方自动审核服务的纯后端 API 服务：基于药品说明书、临床指南与用药规则库，自动检测处方中的用药禁忌、药物相互作用、剂量异常、适应症不符与重复用药，标记风险等级并生成结构化审核报告。

## 快速启动（Docker Compose，首选）

```bash
cd /Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/医疗健康主题项目提示词/lp-345
cp .env.example .env   # 首次使用按需修改
docker compose up -d --build
```

启动后：

- 后端 API：http://localhost:19945
- 健康检查：http://localhost:19945/healthz
- Swagger 文档：http://localhost:19945/swagger/index.html
- PostgreSQL：localhost:44018，Redis：localhost:46318，MinIO API：localhost:47020

默认种子数据（幂等，后端启动时自动写入）：

| 项目 | 值 |
| --- | --- |
| 管理员 | 用户名 `admin`，密码 `admin123` |
| 药品档案 | 氨氯地平、赖诺普利、依那普利、华法林、阿司匹林、二甲双胍、头孢呋辛、地西泮 |
| 相互作用规则 | 5 条（华法林×阿司匹林高风险等） |

## 项目主要功能

1. **处方接收与解析**：RESTful API 接收电子处方，支持 JSON 与 XML 两种格式（`Content-Type: application/json` / `application/xml`），自动解析患者信息、诊断（ICD-10）与药品列表。
2. **适应症审核**：按药品适应症（ICD-10）与处方诊断匹配，不符时标记「适应症不符」。
3. **用法用量审核**：按成人/儿童/老人剂量上限审核日剂量与频次（qd/bid/tid…），超量时给出推荐剂量。
4. **药物相互作用审核**：两两药品基于相互作用规则库分级（高/中/低风险），高风险标记「严重相互作用」。
5. **禁忌与特殊人群审核**：过敏史、妊娠、肝肾功能、老年人 Beers 标准、儿童年龄限制。
6. **重复用药与审核报告**：相同作用机制重复开具检测（如两种 ACEI），生成通过/警告/拒绝的审核报告，支持医生对警告类处方强制通过并记录原因。
7. **异步审核队列**：处方接收后通过 Redis Stream 异步审核（Redis 不可用时自动降级同步审核）。
8. **横切能力**：API Key + JWT 双认证 + RBAC、操作审计日志、请求追踪（request_id）、Redis 限流、MinIO 报告快照导出。

## 技术栈

| 层 | 技术栈 |
| --- | --- |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存/队列 | Redis 7（限流 + Redis Stream 异步审核队列） |
| 对象存储 | MinIO（审核报告快照导出） |
| 认证 | JWT + RBAC + API Key |
| 日志 | `log/slog` 结构化日志 |
| 参数校验 | `github.com/go-playground/validator/v10` |
| 接口文档 | OpenAPI 3.0（/swagger/index.html） |

## 项目目录结构

```
lp-345/
├── docker-compose.yml
├── .env / .env.example
├── README.md
├── database/
│   └── init.sql                 # 数据库初始化脚本
├── deploy/
│   └── README.md                # 部署说明
└── backend/
    ├── Dockerfile               # Go 多阶段构建
    ├── go.mod / go.sum
    ├── api/swagger.json         # OpenAPI 3.0 定义（router 内嵌副本）
    ├── migrations/001_init.sql  # 数据库级迁移脚本
    ├── cmd/server/main.go       # 装配依赖、启动 HTTP 服务与审核 worker
    ├── internal/
    │   ├── config/              # 环境变量配置
    │   ├── database/            # 连接、AutoMigrate、种子数据
    │   ├── model/               # 每实体一个文件
    │   ├── dto/                 # 每实体一个 DTO 文件
    │   ├── repository/          # 每实体一个仓储文件
    │   ├── service/             # 每实体一个服务文件（含审核引擎）
    │   ├── handler/             # 每实体一个处理器文件
    │   ├── router/              # 每实体一个路由文件
    │   ├── middleware/          # request_id/auth/api_key/dual_auth/rbac/rate_limit/audit/error_handler/cors
    │   ├── constants/           # 枚举、错误码、日志模板、消息
    │   ├── queue/               # Redis Stream 审核队列
    │   ├── worker/              # 异步审核消费者
    │   ├── storage/             # MinIO 客户端
    │   └── util/                # logger/jwt/response/app_error/formatters/pagination/validator
    ├── pkg/
    │   ├── pointerx/            # 指针工具
    │   └── testutil/            # 测试种子数据
    └── output/execution.md      # 执行验证报告
```

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `rxcheck` | Compose 项目名（容器/卷名前缀） |
| `BACKEND_PORT` | `19945` | 后端宿主机端口 |
| `DB_PORT` | `44018` | PostgreSQL 宿主机端口 |
| `REDIS_PORT` | `46318` | Redis 宿主机端口 |
| `MINIO_PORT` / `MINIO_CONSOLE_PORT` | `47020` / `47021` | MinIO API / 控制台端口 |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | `rxcheck_db` / `rxcheck_user` / `rxcheck_pwd` | 数据库连接 |
| `JWT_SECRET` | `change_me_to_a_long_random_string` | JWT 签名密钥 |
| `API_KEY_SECRET` | `change_me_to_a_long_random_string` | 服务间调用 API Key |
| `RATE_LIMIT_PER_SECOND` | `200` | 每 IP 每秒请求上限 |
| `MINIO_ROOT_USER` / `MINIO_ROOT_PASSWORD` | `rxcheck_minio` / `rxcheck_minio_pwd` | MinIO 凭据 |

## 本地开发（备选）

```bash
cd backend
go mod tidy
go run ./cmd/server
go build ./...
go vet ./...
go test ./...
```

本地启动需自行准备 PostgreSQL/Redis/MinIO，并通过环境变量覆盖默认连接（见 `internal/config/config.go`）。

## API 调用示例（curl）

### 1. 登录获取 JWT

```bash
curl -sS -X POST http://localhost:19945/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
# 返回 data.token
```

### 2. 注册医生账号

```bash
curl -sS -X POST http://localhost:19945/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"dr_li","password":"secret123","name":"李医生"}'
```

### 3. 使用 JWT 接收 JSON 处方并自动审核

```bash
TOKEN=<上一步返回的 token>
curl -sS -X POST http://localhost:19945/api/v1/prescriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "prescription_no": "RX202608170001",
    "patient_name": "张三",
    "patient_age": 68,
    "patient_gender": "male",
    "allergies": ["青霉素"],
    "pregnant": false,
    "hepatic_impairment": false,
    "renal_impairment": false,
    "diagnoses": [{"icd10":"I10","name":"原发性高血压"}],
    "items": [
      {"drug_name":"氨氯地平","specification":"5mg","single_dose":5,"frequency":"qd","course_days":30,"route":"口服"},
      {"drug_name":"赖诺普利","specification":"10mg","single_dose":10,"frequency":"qd","course_days":30,"route":"口服"},
      {"drug_name":"阿司匹林","specification":"100mg","single_dose":100,"frequency":"qd","course_days":30,"route":"口服"}
    ]
  }'
```

### 4. 使用 XML 提交处方

```bash
cat > rx.xml <<'XML'
<?xml version="1.0" encoding="UTF-8"?>
<SubmitPrescriptionReq>
  <prescription_no>RX202608170002</prescription_no>
  <patient_name>李四</patient_name>
  <patient_age>30</patient_age>
  <patient_gender>female</patient_gender>
  <allergies><item>青霉素</item></allergies>
  <diagnoses>
    <diagnosis><icd10>J18</icd10><name>肺炎</name></diagnosis>
  </diagnoses>
  <items>
    <item><drug_name>头孢呋辛</drug_name><single_dose>250</single_dose><frequency>bid</frequency><course_days>7</course_days><route>口服</route></item>
  </items>
</SubmitPrescriptionReq>
XML
curl -sS -X POST http://localhost:19945/api/v1/prescriptions \
  -H "Content-Type: application/xml" \
  -H "Authorization: Bearer $TOKEN" \
  --data-binary @rx.xml
```

### 5. 查询处方详情 / 审核报告

```bash
curl -sS http://localhost:19945/api/v1/prescriptions/1 -H "Authorization: Bearer $TOKEN"
curl -sS http://localhost:19945/api/v1/prescriptions/1/report -H "Authorization: Bearer $TOKEN"
curl -sS "http://localhost:19945/api/v1/prescriptions?page=1&page_size=10&status=warned" -H "Authorization: Bearer $TOKEN"
curl -sS http://localhost:19945/api/v1/reports -H "Authorization: Bearer $TOKEN"
```

### 6. 医生强制通过警告类处方

```bash
curl -sS -X POST http://localhost:19945/api/v1/prescriptions/1/override \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"reason":"患者合并房颤，抗栓治疗获益大于风险"}'
```

### 7. 使用 API Key 访问（服务间调用，双认证）

```bash
curl -sS http://localhost:19945/api/v1/drugs \
  -H "X-API-Key: change_me_to_a_long_random_string_apikey_rxcheck_2026"
```

### 8. 导出报告快照到 MinIO

```bash
curl -sS -X POST http://localhost:19945/api/v1/prescriptions/1/report/export \
  -H "Authorization: Bearer $TOKEN"
```

### 9. 统计与审计

```bash
curl -sS http://localhost:19945/api/v1/stats/overview -H "Authorization: Bearer $TOKEN"
curl -sS http://localhost:19945/api/v1/audit-logs -H "Authorization: Bearer $TOKEN"
```

## 核心 API 清单

| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /api/v1/auth/login | 登录签发 JWT | 公开 |
| POST | /api/v1/auth/register | 注册医生账号 | 公开 |
| GET | /api/v1/users | 用户列表 | admin |
| POST | /api/v1/users | 创建用户 | admin |
| PUT | /api/v1/users/:id | 更新用户 | admin |
| GET | /api/v1/drugs | 药品列表 | 登录/API Key |
| GET | /api/v1/drugs/:id | 药品详情 | 登录/API Key |
| POST | /api/v1/drugs | 创建药品 | admin |
| PUT | /api/v1/drugs/:id | 更新药品 | admin |
| POST | /api/v1/drugs/:id/disable · enable | 停用/启用药品 | admin |
| GET | /api/v1/interactions | 相互作用规则列表 | 登录/API Key |
| POST | /api/v1/interactions | 创建规则 | admin |
| DELETE | /api/v1/interactions/:id | 删除规则 | admin |
| POST | /api/v1/prescriptions | 接收处方（JSON/XML）并审核 | 登录/API Key |
| GET | /api/v1/prescriptions | 处方列表 | 登录/API Key |
| GET | /api/v1/prescriptions/:id | 处方详情 | 登录/API Key |
| POST | /api/v1/prescriptions/:id/review | 重新审核 | 登录/API Key |
| POST | /api/v1/prescriptions/:id/override | 强制通过（记录原因） | doctor/pharmacist/admin |
| GET | /api/v1/prescriptions/:id/report | 审核报告 | 登录/API Key |
| POST | /api/v1/prescriptions/:id/report/export | 导出报告到 MinIO | 登录/API Key |
| GET | /api/v1/reports | 报告列表 | 登录/API Key |
| GET | /api/v1/stats/overview | 审核统计概览 | 登录/API Key |
| GET | /api/v1/audit-logs | 审计日志 | admin |
| GET | /healthz | 健康检查 | 公开 |

**复用关系**：`GET /api/v1/prescriptions/:id/report` 与 `GET /api/v1/prescriptions/:id` 复用 `PrescriptionRepository.FindByID`；`POST /api/v1/prescriptions/:id/review` 与处方接收后的异步审核复用 `ReviewService.ReviewPrescription`；列表接口复用 `PrescriptionRepository.List`/`CountByStatus`。

## 枚举出现位置清单

以下业务枚举在后端所有出现位置（**新增枚举值必须同步修改全部位置**）：

### 1. 风险等级 RiskLevel（none/low/medium/high/critical）

- 定义：`internal/constants/risk_level.go`
- 模型：`internal/model/interaction_rule.go`（RiskLevel 字段）、`internal/model/review_report.go`（RiskLevel/ReviewItem.RiskLevel）、`internal/model/prescription.go`（OverallRisk）
- DTO：`internal/dto/interaction_dto.go`（oneof=high medium low）
- service 状态机：`internal/service/review_service.go`（MapRiskToStatus、MaxRiskLevel、RiskRank、各审核规则风险赋值）
- handler 校验：`internal/handler/interaction_handler.go` 经 DTO binding oneof 校验
- 错误码/日志模板：`internal/constants/error_codes.go`、`internal/constants/log_templates.go`（LogPrescriptionReviewed 含 risk）
- 格式化：`internal/util/formatters.go`（RiskLevelText）
- 种子数据：`internal/database/seed.go`、`pkg/testutil/seed.go`

### 2. 审核状态 PrescriptionStatus（pending_review/reviewing/passed/warned/rejected/overridden）

- 定义：`internal/constants/review_status.go`（含状态机映射 MapRiskToStatus、终态判断）
- 模型：`internal/model/prescription.go`（Status 字段）
- DTO：`internal/dto/prescription_dto.go`（PrescriptionListResp.Status）
- service 状态机：`internal/service/prescription_service.go`（Submit/Override）、`internal/service/review_service.go`（ReviewPrescription 落库状态）
- handler 校验：`internal/handler/prescription_handler.go`（Override 仅警告类）
- 错误码：`internal/constants/error_codes.go`（CodePrescriptionLocked、CodeOverrideNotAllowed）、`internal/constants/messages.go`（MsgPrescriptionLocked、MsgOverrideNotAllowed）
- 日志模板：`internal/constants/log_templates.go`（LogPrescriptionSubmitted/Reviewed/Overridden）
- 格式化：`internal/util/formatters.go`（PrescriptionStatusText）
- 路由权限：`internal/router/prescription.go`

### 3. 规则类型 RuleType（indication/dosage/interaction/contraindication/duplication）

- 定义：`internal/constants/rule_type.go`
- 模型：`internal/model/review_report.go`（ReviewItem.RuleType）
- service：`internal/service/review_service.go`（五大审核规则逐一使用）
- 格式化：`internal/util/formatters.go`（RuleTypeText）
- 日志模板：`internal/constants/log_templates.go`

### 4. 用户角色 UserRole（admin/doctor/pharmacist）

- 定义：`internal/constants/user_role.go`
- 模型：`internal/model/user.go`（Role 字段）
- DTO：`internal/dto/user_dto.go`（CreateUserReq.Role oneof）、`internal/dto/common_dto.go`
- service：`internal/service/auth_service.go`（注册默认 doctor）、`internal/service/user_service.go`（ValidateRole）
- handler 校验：`internal/handler/user_handler.go`
- 错误码/消息：`internal/constants/error_codes.go`（CodeForbidden）、`internal/constants/messages.go`（MsgForbidden）
- 日志模板：`internal/constants/log_templates.go`（LogLoginSuccess 等）
- RBAC：`internal/middleware/rbac.go`、`internal/router/*.go`

### 5. 实体状态 EntityStatus（active/disabled、enabled/disabled）

- 定义：`internal/constants/entity_status.go`
- 模型：`internal/model/user.go`、`internal/model/drug.go`、`internal/model/interaction_rule.go`
- service：`internal/service/drug_service.go`（Disable/Enable）、`internal/service/user_service.go`
- 日志模板：`internal/constants/log_templates.go`（LogDrugDisabled/Enabled 等）

## Docker 部署说明

- 端口映射：后端 `${BACKEND_PORT:-19945}:8080`；数据库 `${DB_PORT:-44018}:5432`；Redis `${REDIS_PORT:-46318}:6379`；MinIO `${MINIO_PORT:-47020}:9000`、`${MINIO_CONSOLE_PORT:-47021}:9001`。
- 数据持久化：PostgreSQL/Redis/MinIO 使用命名卷（`rxcheck-db-data` / `rxcheck-redis-data` / `rxcheck-minio-data`），`docker compose down` 不丢数据；`docker compose down -v` 清空数据。
- 容器健康检查：db 使用 `pg_isready`，backend 使用 `/healthz`，启动顺序通过 `depends_on: condition: service_healthy` 保证。
- 任意目录名（含中文）下均可启动：所有卷/容器名带 `${COMPOSE_PROJECT_NAME:-rxcheck}` 前缀，未使用绑定挂载。

### 常见问题

1. **端口被占用**：修改 `.env` 中 `BACKEND_PORT` 等端口后 `docker compose up -d` 重启。
2. **容器不健康**：`docker compose logs backend` 查看日志；多为数据库未就绪或密钥长度不足（JWT 密钥过短会报错）。
3. **审核未触发**：确认 Redis 可用（`docker compose logs backend` 出现「Redis 不可用」则降级为同步审核，不影响功能）。
4. **报告导出失败**：确认 MinIO 已启动且 `MINIO_*` 环境变量一致。

## License

MIT License。本项目仅用于医疗健康主题的工程演示，不构成真实医疗建议；生产使用需经专业药事审核与合规评估。
