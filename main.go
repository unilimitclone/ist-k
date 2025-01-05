package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/alist-org/alist/v3/pkg/utils/random"
)

// 数据库配置结构体
type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASS"`
	Name        string `json:"name" env:"NAME"`
	DBFile      string `json:"db_file" env:"FILE"`
	TablePrefix string `json:"table_prefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"ssl_mode" env:"SSL_MODE"`
	DSN         string `json:"dsn" env:"DSN"`
}

// Meilisearch 结构体
type Meilisearch struct {
	Host        string `json:"host" env:"HOST"`
	APIKey      string `json:"api_key" env:"API_KEY"`
	IndexPrefix string `json:"index_prefix" env:"INDEX_PREFIX"`
}

// 服务器配置结构体
type Scheme struct {
	Address      string `json:"address" env:"ADDR"`
	HttpPort     int    `json:"http_port" env:"HTTP_PORT"`
	HttpsPort    int    `json:"https_port" env:"HTTPS_PORT"`
	ForceHttps   bool   `json:"force_https" env:"FORCE_HTTPS"`
	CertFile     string `json:"cert_file" env:"CERT_FILE"`
	KeyFile      string `json:"key_file" env:"KEY_FILE"`
	UnixFile     string `json:"unix_file" env:"UNIX_FILE"`
	UnixFilePerm string `json:"unix_file_perm" env:"UNIX_FILE_PERM"`
}

// 日志配置结构体
type LogConfig struct {
	Enable     bool   `json:"enable" env:"LOG_ENABLE"`
	Name       string `json:"name" env:"LOG_NAME"`
	MaxSize    int    `json:"max_size" env:"MAX_SIZE"`
	MaxBackups int    `json:"max_backups" env:"MAX_BACKUPS"`
	MaxAge     int    `json:"max_age" env:"MAX_AGE"`
	Compress   bool   `json:"compress" env:"COMPRESS"`
}

// TaskConfig 结构体
type TaskConfig struct {
	Workers        int  `json:"workers" env:"WORKERS"`
	MaxRetry       int  `json:"max_retry" env:"MAX_RETRY"`
	TaskPersistant bool `json:"task_persistant" env:"TASK_PERSISTANT"`
}

// TasksConfig 结构体
type TasksConfig struct {
	Download TaskConfig `json:"download" envPrefix:"DOWNLOAD_"`
	Transfer TaskConfig `json:"transfer" envPrefix:"TRANSFER_"`
	Upload   TaskConfig `json:"upload" envPrefix:"UPLOAD_"`
	Copy     TaskConfig `json:"copy" envPrefix:"COPY_"`
}

// Cors 结构体
type Cors struct {
	AllowOrigins []string `json:"allow_origins" env:"ALLOW_ORIGINS"`
	AllowMethods []string `json:"allow_methods" env:"ALLOW_METHODS"`
	AllowHeaders []string `json:"allow_headers" env:"ALLOW_HEADERS"`
}

// S3 结构体
type S3 struct {
	Enable bool `json:"enable" env:"ENABLE"`
	Port   int  `json:"port" env:"PORT"`
	SSL    bool `json:"ssl" env:"SSL"`
}

// FTP 结构体
type FTP struct {
	Enable                  bool   `json:"enable" env:"ENABLE"`
	Listen                  string `json:"listen" env:"LISTEN"`
	FindPasvPortAttempts    int    `json:"find_pasv_port_attempts" env:"FIND_PASV_PORT_ATTEMPTS"`
	ActiveTransferPortNon20 bool   `json:"active_transfer_port_non_20" env:"ACTIVE_TRANSFER_PORT_NON_20"`
	IdleTimeout             int    `json:"idle_timeout" env:"IDLE_TIMEOUT"`
	ConnectionTimeout       int    `json:"connection_timeout" env:"CONNECTION_TIMEOUT"`
	DisableActiveMode       bool   `json:"disable_active_mode" env:"DISABLE_ACTIVE_MODE"`
	DefaultTransferBinary   bool   `json:"default_transfer_binary" env:"DEFAULT_TRANSFER_BINARY"`
	EnableActiveConnIPCheck bool   `json:"enable_active_conn_ip_check" env:"ENABLE_ACTIVE_CONN_IP_CHECK"`
	EnablePasvConnIPCheck   bool   `json:"enable_pasv_conn_ip_check" env:"ENABLE_PASV_CONN_IP_CHECK"`
}

// SFTP 结构体
type SFTP struct {
	Enable bool   `json:"enable" env:"ENABLE"`
	Listen string `json:"listen" env:"LISTEN"`
}

// 配置结构体
type Config struct {
	Force                 bool        `json:"force" env:"FORCE"`
	SiteURL               string      `json:"site_url" env:"SITE_URL"`
	Cdn                   string      `json:"cdn" env:"CDN"`
	JwtSecret             string      `json:"jwt_secret" env:"JWT_SECRET"`
	TokenExpiresIn        int         `json:"token_expires_in" env:"TOKEN_EXPIRES_IN"`
	Database              Database    `json:"database" envPrefix:"DB_"`
	Meilisearch           Meilisearch `json:"meilisearch" envPrefix:"MEILISEARCH_"`
	Scheme                Scheme      `json:"scheme"`
	TempDir               string      `json:"temp_dir" env:"TEMP_DIR"`
	BleveDir              string      `json:"bleve_dir" env:"BLEVE_DIR"`
	DistDir               string      `json:"dist_dir"`
	Log                   LogConfig   `json:"log"`
	DelayedStart          int         `json:"delayed_start" env:"DELAYED_START"`
	MaxConnections        int         `json:"max_connections" env:"MAX_CONNECTIONS"`
	TlsInsecureSkipVerify bool        `json:"tls_insecure_skip_verify" env:"TLS_INSECURE_SKIP_VERIFY"`
	Tasks                 TasksConfig `json:"tasks" envPrefix:"TASKS_"`
	Cors                  Cors        `json:"cors" envPrefix:"CORS_"`
	S3                    S3          `json:"s3" envPrefix:"S3_"`
	FTP                   FTP         `json:"ftp" envPrefix:"FTP_"`
	SFTP                  SFTP        `json:"sftp" envPrefix:"SFTP_"`
}

// 初始化配置
func initConfig(configFilePath string) *Config {
	DATABASE_URL := os.Getenv("DATABASE_URL")
	fmt.Println("DatabaseUrl", DATABASE_URL)
	// DATABASE_URL = "postgres://hfhgpvbymdzusj:39d7f6f3ee4288103e382d5dec22ce668c4e5cb65120f64d574b808775674eb4@ec2-3-218-171-44.compute-1.amazonaws.com:5432/d4o07n33pf6ot7"
	u, err := url.Parse(DATABASE_URL)
	if err != nil {
		fmt.Println(err)
	}
	user := u.User.Username()
	pass, _ := u.User.Password()
	host := u.Hostname()
	port, _ := strconv.Atoi(u.Port())
	name := u.Path[1:]
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if _, err := os.Stat(configFilePath); err == nil {
			data, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				log.Fatalf("Failed to read config file: %s", err)
			}
			var existingConfig Config
			err = json.Unmarshal(data, &existingConfig)
			if err != nil {
				log.Fatalf("Failed to unmarshal existing config: %s", err)
			}
			jwtSecret = existingConfig.JwtSecret
		}
	}
	if jwtSecret == "" {
		jwtSecret = random.String(16)
	}
	return &Config{
		Scheme: Scheme{
			Address:    "0.0.0.0",
			UnixFile:   "",
			HttpPort:   5244,
			HttpsPort:  -1,
			ForceHttps: false,
			CertFile:   "",
			KeyFile:    "",
		},
		JwtSecret: jwtSecret,
		TokenExpiresIn: 48,
		TempDir: "data/temp",
		Database: Database{
			User:        user,
			Password:    pass,
			Host:        host,
			Port:        port,
			Name:        name,
			TablePrefix: "x_",
			DBFile:      "data/data.db",
		},
		Meilisearch: Meilisearch{
			Host: "http://localhost:7700",
		},
		BleveDir: "data/bleve",
		Log: LogConfig{
			Enable:     true,
			Name:       "data/log/log.log",
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     28,
		},
		MaxConnections:        0,
		TlsInsecureSkipVerify: true,
		Tasks: TasksConfig{
			Download: TaskConfig{
				Workers:  5,
				MaxRetry: 1,
				// TaskPersistant: true,
			},
			Transfer: TaskConfig{
				Workers:  5,
				MaxRetry: 2,
				// TaskPersistant: true,
			},
			Upload: TaskConfig{
				Workers: 5,
			},
			Copy: TaskConfig{
				Workers:  5,
				MaxRetry: 2,
				// TaskPersistant: true,
			},
		},
		Cors: Cors{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		},
		S3: S3{
			Enable: false,
			Port:   5246,
			SSL:    false,
		},
		FTP: FTP{
			Enable:                  false,
			Listen:                  ":5221",
			FindPasvPortAttempts:    50,
			ActiveTransferPortNon20: false,
			IdleTimeout:             900,
			ConnectionTimeout:       30,
			DisableActiveMode:       false,
			DefaultTransferBinary:   false,
			EnableActiveConnIPCheck: true,
			EnablePasvConnIPCheck:   true,
		},
		SFTP: SFTP{
			Enable: false,
			Listen: ":5222",
		},
	}
}

// 写入配置文件
func main() {
	configFilePath := "/opt/alist/data/config.json"
	config := initConfig(configFilePath)
	confBody, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal JSON: %s", err.Error())
	}
	err = ioutil.WriteFile(configFilePath, confBody, 0644)
	if err != nil {
		log.Fatalf("failed to write JSON: %s", err.Error())
	}
}