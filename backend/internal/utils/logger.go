package utils

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

func InitLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return
	}

	// 日志文件按日期命名
	today := time.Now().Format("2006-01-02")
	infoFile, err := os.OpenFile(filepath.Join(logDir, "info-"+today+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开信息日志文件失败: %v", err)
		return
	}

	errorFile, err := os.OpenFile(filepath.Join(logDir, "error-"+today+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开错误日志文件失败: %v", err)
		return
	}

	warnFile, err := os.OpenFile(filepath.Join(logDir, "warn-"+today+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("打开警告日志文件失败: %v", err)
		return
	}

	InfoLogger = log.New(infoFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(warnFile, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	if InfoLogger != nil {
		InfoLogger.Printf(format, v...)
	}
	log.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Printf(format, v...)
	}
	log.Printf("[ERROR] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	if WarnLogger != nil {
		WarnLogger.Printf(format, v...)
	}
	log.Printf("[WARN] "+format, v...)
}


