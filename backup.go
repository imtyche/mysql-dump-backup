package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func Backup(cfg *Config) {
	// 1. 获取当前日期作为子目录名 (例如: 2026-01-30)
	dateDir := time.Now().Format("2006-01-02")
	targetPath := filepath.Join(cfg.BackupPath, dateDir)

	// 创建当天的备份文件夹
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		fmt.Println("创建日期备份目录失败:", err)
		return
	}

	// 2. 执行备份（并发）
	backupWithDump(cfg, targetPath)

	// 3. 清理指定天数前的旧备份
	cleanOldBackups(cfg.BackupPath, cfg.Clear)
}

// backupWithDump 使用 mariadb-dump，支持并发 + 可选 gzip
func backupWithDump(cfg *Config, targetPath string) {
	sem := make(chan struct{}, cfg.Thread)
	var wg sync.WaitGroup

	for _, db := range cfg.Databases {
		wg.Add(1)
		go func(dbName string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			backupDatabase(cfg, dbName, targetPath)
		}(db.Name)
	}
	wg.Wait()
}

func backupDatabase(cfg *Config, dbName string, targetPath string) {
	timestamp := time.Now().Format("150405")
	ext := ".sql"
	if cfg.Gzip {
		ext = ".sql.gz"
	}
	fileName := fmt.Sprintf("%s_%s%s", dbName, timestamp, ext)
	filePath := filepath.Join(targetPath, fileName)

	args := []string{
		"-h", cfg.MySQL.Host,
		"-P", fmt.Sprintf("%d", cfg.MySQL.Port),
		"-u", cfg.MySQL.User,
		fmt.Sprintf("-p%s", cfg.MySQL.Password),
		"--ssl=0",
		"--single-transaction",
		"--flush-privileges",
		"--quick",
		"--routines",
		"--events",
		"--skip-comments",  // 去掉 dump 文件中的注释
		"--skip-dump-date", // 去掉文件头的 dump 日期注释
		"--skip-set-charset",
		dbName,
	}

	fmt.Printf("[%s] 开始备份数据库: %s (gzip=%v)\n",
		time.Now().Format("15:04:05"), dbName, cfg.Gzip)

	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("创建备份文件失败:", err)
		return
	}
	defer outFile.Close()

	cmd := exec.Command("mariadb-dump", args...)
	cmd.Stderr = os.Stderr

	if cfg.Gzip {
		// 使用纯 Go compress/gzip，不依赖系统 gzip 命令
		gzWriter := gzip.NewWriter(outFile)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("创建 stdout pipe 失败:", err)
			return
		}

		if err := cmd.Start(); err != nil {
			fmt.Println("启动 mariadb-dump 失败:", err)
			return
		}

		// 把 dump 输出流式写入 gzip
		_, copyErr := io.Copy(gzWriter, stdout)
		waitErr := cmd.Wait()
		closeErr := gzWriter.Close() // 必须 Close 才能写出完整 gzip footer

		if copyErr != nil {
			fmt.Println("写入压缩数据失败:", dbName, copyErr)
			return
		}
		if waitErr != nil {
			fmt.Println("备份失败:", dbName, waitErr)
			return
		}
		if closeErr != nil {
			fmt.Println("关闭 gzip 失败:", dbName, closeErr)
			return
		}
	} else {
		cmd.Stdout = outFile
		if err := cmd.Run(); err != nil {
			fmt.Println("备份失败:", dbName, err)
			return
		}
	}

	fmt.Printf("[%s] 备份数据库成功: %s\n", time.Now().Format("15:04:05"), filePath)
}

// cleanOldBackups 用于删除指定目录下超过指定天数的文件或目录(根据mod时间而不是文件夹名字)
func cleanOldBackups(backupRoot string, clear int) {
	fmt.Println("检查并清理过期备份...")
	files, err := os.ReadDir(backupRoot)
	if err != nil {
		fmt.Println("读取备份根目录失败:", err)
		return
	}

	now := time.Now()
	clearDays := now.AddDate(0, 0, -clear)

	for _, file := range files {
		fullPath := filepath.Join(backupRoot, file.Name())
		info, err := file.Info()
		if err != nil {
			continue
		}

		// 如果文件的修改时间早于 clear 天前，则删除
		if info.ModTime().Before(clearDays) {
			err := os.RemoveAll(fullPath) // 使用 RemoveAll 可以递归删除文件夹
			if err != nil {
				fmt.Printf("删除过期备份失败 [%s]: %v\n", fullPath, err)
			} else {
				fmt.Printf("已成功清理过期备份: %s\n", fullPath)
			}
		}
	}
}
