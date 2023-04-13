package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		HTTP     map[string]HTTPConfig
		Redis    RedisConfig
		Postgres PostgresConfig
		MongoDB  MongoDBConfig
		Custom   map[string]string
		Auth     Auth
	}

	Auth struct {
		Username string `mapstructure:"username" yaml:"username"`
		Password string `mapstructure:"password" yaml:"password"`
	}

	HTTPConfig struct {
		HostString         string        `mapstructure:"hostString"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderMegabytes"`
	}

	RedisConfig struct {
		HostString    string `mapstructure:"hostString"`
		Login         string `mapstructure:"login"`
		Password      string `mapstructure:"password"`
		ConnectsCount int    `mapstructure:"connectionsCount"`
		DBNum         int    `mapstructure:"dbNum"`
	}

	MongoDBConfig struct {
		HostString string `mapstructure:"hostString"`
		Port       int    `mapstructure:"port"`
		Login      string `mapstructure:"login"`
		Password   string `mapstructure:"password"`
		DBName     string `mapstructure:"dbName"`
	}

	PostgresConfig struct {
		HostString    string `mapstructure:"hostString"`
		Port          int    `mapstructure:"port"`
		Login         string `mapstructure:"login"`
		Password      string `mapstructure:"password"`
		DBName        string `mapstructure:"dbName"`
		ConnectsCount int    `mapstructure:"connectionsCount"`
	}
)

func MakeServerList(hostString string) []string {
	return strings.Split(hostString, ",")
}

func (c *Config) GetPostgres() *PostgresConfig {
	return &c.Postgres
}

func (c *Config) GetMongoDB() *MongoDBConfig {
	return &c.MongoDB
}

func (c *Config) GetCustom(key string) string {
	if config, ok := c.Custom[key]; ok {
		return config
	}
	log.Printf("Cannot get Custom config key %v\n", key)
	return ""
}

func (c *Config) GetHTTP(key string) *HTTPConfig {
	if config, ok := c.HTTP[key]; ok {
		return &config
	}
	log.Printf("Cannot get HTTP config key %v\n", key)
	return nil
}

func (c *Config) GetRedis() *RedisConfig {
	return &c.Redis
}

func (c *Config) BuildDSNPostgres() string {
	pg := c.GetPostgres()
	if pg.HostString == "" {
		log.Fatal("postgres hostname not set")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		pg.HostString, pg.Login, pg.Password, pg.DBName, pg.Port,
	)
	return dsn
}

func New(path, serviceName string) (*Config, error) {
	if err := parseConfigFile(path, serviceName); err != nil {
		return nil, err
	}
	// If in os env set GO_CLOUD_CONFIG=true load from config server
	if viper.GetBool("CLOUD_CONFIG") {
		if viper.GetString("CLOUD_URL") == "" {
			return nil, fmt.Errorf("CLOUD_URL not set")
		}
		if serviceName == "" {
			serviceName = viper.GetString("SERVICE_NAME")
		}
		loadConfiguration(viper.GetString("CLOUD_URL"), serviceName, "default")
	} else {
		log.Printf("loading config from %s\n", path)
	}
	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("custom", &cfg.Custom); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("mongodb", &cfg.MongoDB); err != nil {
		return err
	}
	cfg.Auth.Username = viper.GetString("auth.username")
	cfg.Auth.Password = viper.GetString("auth.password")
	return nil
}

// / configs/config
func parseConfigFile(filepath, serviceName string) error {
	path := strings.Split(filepath, "/")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("go")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(strings.Join(path[:len(path)-1], "/")) // folder
	viper.SetConfigName(serviceName)                           // config file name
	return viper.ReadInConfig()
}
