package main

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases  []Database `yaml:"databases"`
	BackupPath string     `yaml:"backup_path"`
	MySQL      MySQL      `yaml:"mysql"`
	Cron       string     `yaml:"cron"`
	Des        string     `yaml:"des"`
	Clear      int        `yaml:"clear"`
	Init       bool       `yaml:"init"`
	Gzip       bool       `yaml:"gzip"`    // whether to compress backup
	Thread     int        `yaml:"thread"`  // concurrency level
	Command    string     `yaml:"command"` // dump command: mysqldump / mariadb-dump (auto-detect if empty)
}

type Database struct {
	Name string `yaml:"name"`
}

type MySQL struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// detectDumpCommand 自动检测系统中可用的 dump 命令
// 优先顺序: mariadb-dump -> mysqldump
func detectDumpCommand() (string, error) {
	candidates := []string{"mariadb-dump", "mysqldump"}
	for _, cmd := range candidates {
		if path, err := exec.LookPath(cmd); err == nil {
			return path, nil // 返回完整路径，更可靠
		}
	}
	return "", fmt.Errorf("未找到可用的 dump 命令，请安装 mariadb-dump 或 mysqldump，或在 config.yml 中手动指定 command")
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// defaults
	if cfg.Thread <= 0 {
		cfg.Thread = 1
	}
	if cfg.Clear <= 0 {
		cfg.Clear = 7
	}

	// 兼容 mysqldump / mariadb-dump
	if cfg.Command == "" {
		detected, err := detectDumpCommand()
		if err != nil {
			return nil, err
		}
		cfg.Command = detected
	} else {
		// 用户指定了命令，验证是否存在（允许相对名或完整路径）
		if _, err := exec.LookPath(cfg.Command); err != nil {
			return nil, fmt.Errorf("配置的 command %q 不可用: %w", cfg.Command, err)
		}
	}

	return &cfg, nil
}
