# 安装 Goose 工具
install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

# 启动数据库迁移
up:
	goose mysql "root:zl020613@tcp(localhost:3309)/user?parseTime=true" up

# 回滚数据库迁移
down:
	goose mysql "root:zl020613@tcp(localhost:3309)/user?parseTime=true" down

# 创建新的迁移文件
create:
	goose mysql "root:zl020613@tcp(localhost:3309)/user?parseTime=true" create ${name} sql
