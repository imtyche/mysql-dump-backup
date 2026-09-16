# MySQL Dump Backup

🚀 **A powerful MySQL backup tool with automated scheduling and restoration capabilities.**

## 📋 Overview

MySQL Dump Backup is a lightweight, efficient tool built in Go for automating MySQL database backups. It supports scheduled backups using cron expressions, automatic cleanup of old backups, and Docker deployment for easy integration into your infrastructure.

## ✨ Features

- **⏰ Scheduled Backups** - Use cron expressions to schedule backups at any interval
- **🗂️ Organized Storage** - Backups automatically organized by date in subdirectories
- **🧹 Automatic Cleanup** - Automatically removes backups older than a specified number of days
- **🐳 Docker Support** - Deploy easily with Docker and Docker Compose
- **🔒 Transaction Support** - Uses `--single-transaction` for consistent backups without locking
- **🎯 Batch Operations** - Back up multiple databases in a single run
- **✅ Graceful Shutdown** - Safely handles signals (SIGINT, SIGTERM)

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Scheduling:** [robfig/cron/v3](https://github.com/robfig/cron)
- **Config:** YAML format via [gopkg.in/yaml.v3](https://github.com/go-yaml/yaml)
- **Containerization:** Docker & Docker Compose

## 📦 Installation

### Prerequisites
- MySQL/MariaDB server
- Go 1.21+ (for source compilation)
- Docker (for containerized deployment)

### From Source
```bash
git clone https://github.com/wayhelper/mysql-dump-backup.git
cd mysql-dump-backup
go build -o mysql-dump-backup
```

### Docker
```bash
docker pull wayhelper/mysql-dump-backup:latest
```

## ⚙️ Configuration

Create a `config.yml` file in your working directory:

```yaml
# List of databases to backup
databases:
  - name: "database1"
  - name: "database2"
  - name: "database3"
  - name: "--all-databases" # 备份整个数据库

# Local path where backups will be stored
backup_path: "./backups"

# MySQL connection settings
mysql:
  host: "localhost"
  port: 3306
  user: "root"
  password: "your_password"

# Cron expression for backup scheduling
# Format: second minute hour day month day-of-week
# Example: "0 0 2 * * *" = Daily at 2:00 AM
cron: "0 0 2 * * *"

# Description of the backup schedule
des: "Daily backup at 2:00 AM"

# Number of days to retain backups (old backups are automatically deleted)
clear: 30
```

### Environment Variables
You can override the config path using an environment variable:
```bash
echo "export CONFIG_PATH=/path/to/config.yml"
```

## 🚀 Usage

### Running Locally
```bash
./mysql-dump-backup
```

The tool will:
1. Load the configuration from `config.yml`
2. Initialize the cron scheduler
3. Execute backups at the scheduled time
4. Clean up old backups automatically

### Running with Docker

Using Docker Run:
```bash
docker run -d \
  -v /path/to/config.yml:/app/config.yml \
  -v /path/to/backups:/app/backups \
  -e TZ=Asia/Shanghai \
  -e CONFIG_PATH=/app/config.yml \
  --restart unless-stopped \
  wayhelper/mysql-dump-backup:latest
```

Using Docker Compose:
```bash
docker-compose up -d
```

Edit `docker-compose.yml` to set your configuration and volume paths.

## 📁 Backup Structure

Backups are organized by date:
```
backups/
├── 2026-04-24/
│   ├── database1_020534.sql
│   ├── database2_020534.sql
│   └── database3_020534.sql
├── 2026-04-25/
│   ├── database1_020534.sql
│   ├── database2_020534.sql
│   └── database3_020534.sql
```

Each directory is named with the backup date (YYYY-MM-DD), and files are named with the database name and timestamp (HHMMSS).

## 🔧 Development

### Project Structure
- `main.go` - Application entry point and cron scheduler
- `backup.go` - Core backup and cleanup logic
- `config.go` - Configuration loading and parsing
- `Dockerfile` - Multi-stage Docker build configuration
- `docker-compose.yml` - Docker Compose setup example

### Build for Different Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o mysql-dump-backup-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o mysql-dump-backup-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o mysql-dump-backup.exe
go build -ldflags "-s -w" -o mysql-dump-backup.exe .
```

## 🐛 Troubleshooting

### Backup fails with "command not found: mysqldump"
- Ensure MySQL client tools are installed on your system or in the Docker image
- For Docker, the image already includes MySQL client tools

### Permissions denied when accessing config or backups
- Check file permissions: `chmod 755 config.yml backups/`
- Ensure the container user has proper access

### Cron expression not working as expected
- Verify your cron format (includes seconds): `second minute hour day month day-of-week`
- Use an online cron expression validator for testing

## 📋 Cron Expression Examples

| Expression | Description |
|-----------|-------------|
| `0 0 2 * * *` | Daily at 2:00 AM |
| `0 0 0 * * 0` | Weekly (Sunday) at midnight |
| `0 0 0 1 * *` | Monthly (1st) at midnight |
| `0 */6 * * * *` | Every 6 hours |
| `0 0 * * * 1-5` | Weekdays at midnight |

## 📄 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 👤 Author

[wayhelper](https://github.com/wayhelper)

---

**Support:** If you have any questions or issues, please open an [issue](https://github.com/wayhelper/mysql-dump-backup/issues) on GitHub.