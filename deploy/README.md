# 部署说明

- Docker Compose 一键部署：`docker compose up -d --build`
- 后端监听：`http://localhost:19945`
- 健康检查：`http://localhost:19945/healthz`
- Swagger：`http://localhost:19945/swagger/index.html`
- 数据库：PostgreSQL 16（`localhost:44018`），Redis 7（`localhost:46318`），MinIO（`localhost:47020`）

## 端口映射

| 服务 | 容器内端口 | 宿主机端口（默认） |
| --- | --- | --- |
| backend | 8080 | 19945 |
| db | 5432 | 44018 |
| redis | 6379 | 46318 |
| minio | 9000（控制台 9001 仅容器内） | 47020 |

## 常见问题

1. 端口被占用：修改 `.env` 中的 `BACKEND_PORT` 后重启。
2. 数据持久化：PostgreSQL/Redis/MinIO 均使用命名卷，`docker compose down` 不丢数据；`docker compose down -v` 会清空。
