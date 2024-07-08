package config

import (
	"fmt"
	"time"

	_ "time/tzdata"

	"github.com/crt379/svc-collector-grpc/internal/flags"
	"github.com/crt379/svc-collector-grpc/internal/util"

	"github.com/spf13/viper"
)

var AppConfig Config

func init() {
	viper.SetConfigFile(*flags.Flags.Cfg)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("读取配置异常: %s, error: %s", *flags.Flags.Cfg, err.Error()))
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Sprintf("解析配置异常: %s, error: %s", *flags.Flags.Cfg, err.Error()))
	}

	fmt.Printf("配置: %v\n", AppConfig)

	if AppConfig.TZ != "" {
		if cst, err := time.LoadLocation(AppConfig.TZ); err != nil {
			fmt.Println(err.Error())
		} else {
			time.Local = cst
		}
	}

	AppConfig.Host = AppConfig.Listen.Host
	if AppConfig.Listen.Host == "" || AppConfig.Listen.Host == "*" {
		AppConfig.Host = util.GetIP()
	}

	AppConfig.Addr = fmt.Sprintf("%s:%s", AppConfig.Host, AppConfig.Listen.Port)
}

type Config struct {
	TZ         string         `toml:"TZ"`
	Host       string         `toml:"host"`
	Addr       string         `toml:"addr"`
	Register   RegisterConfig `toml:"register"`
	Listen     AddrConfig     `toml:"listen"`
	PgSql      PgSqlConfig    `toml:"pgsql"`
	Redis      RedisConfig    `toml:"redis"`
	Log        LogConfig      `toml:"log"`
	Etcd       []AddrConfig   `toml:"etcd"`
	Prometheus AddrConfig     `toml:"prometheus"`
}

type RegisterConfig struct {
	Name string
}

type AddrConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type PgSqlConfig struct {
	Write PgSqlMeta `toml:"write"`
	Read  PgSqlMeta `toml:"read"`
}

type PgSqlMeta struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
}

type RedisConfig struct {
	Write RedisMeta `toml:"write"`
	Read  RedisMeta `toml:"read"`
}

type RedisMeta struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type LogConfig struct {
	File  string       `toml:"file"`
	Level int          `toml:"level"`
	Redis LogRedisMeta `toml:"redis"`
}

type LogRedisMeta struct {
	Enabled   bool   `toml:"enabled"`
	Key       string `toml:"key"`
	RedisMeta `mapstructure:",squash"`
}
