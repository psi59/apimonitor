package config

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/realsangil/apimonitor/pkg/rserrors"
)

var (
	c              configure
	mux            sync.Mutex
	ConfigFilePath = os.Getenv("AM_CONFIG_PATH")
)

type configure struct {
	Environment string       `mapstructure:"environment"`
	DB          dbConfigure  `mapstructure:"db"`
	Logger      logConfigure `mapstructure:"logger"`
}

func (c *configure) Validate() error {
	if err := c.Logger.Validate(); err != nil {
		return errors.WithStack(err)
	}
	if err := c.DB.Validate(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func Init(configFilePath string) error {
	if configFilePath == "" {
		viper.AddConfigPath("./config")
	}
	viper.AddConfigPath(configFilePath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("server_config")
	viper.SetDefault("environment", "development")
	viper.SetDefault("logger.filepath", "./server.log")
	if err := viper.ReadInConfig(); err != nil {
		return errors.WithStack(err)
	}

	mux.Lock()
	defer mux.Unlock()
	if err := viper.Unmarshal(&c); err != nil {
		return errors.WithStack(err)
	}

	// if err := c.Validate(); err != nil {
	// 	panic(err)
	// }
	return nil
}

type dbConfigure struct {
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Port     uint   `mapstructure:"port"`
	Verbose  bool   `mapstructure:"verbose"`
}

func (c *dbConfigure) GetHost() string {
	return c.Host
}

func (c *dbConfigure) GetPort() uint {
	return c.Port
}

func (c *dbConfigure) GetUsername() string {
	return c.Username
}

func (c *dbConfigure) GetPassword() string {
	return c.Password
}

func (c *dbConfigure) GetDatabaseName() string {
	return c.Name
}

func (c *dbConfigure) GetVerbose() bool {
	return c.Verbose
}

func (c *dbConfigure) Validate() error {
	if c.Host == "" {
		return errors.Wrap(rserrors.ErrInvalidParameter, "db.host")
	}
	if c.Port == 0 {
		return errors.Wrap(rserrors.ErrInvalidParameter, "db.password")
	}
	if c.Username == "" {
		return errors.Wrap(rserrors.ErrInvalidParameter, "db.username")
	}
	if c.Password == "" {
		return errors.Wrap(rserrors.ErrInvalidParameter, "db.password")
	}
	if c.Name == "" {
		return errors.Wrap(rserrors.ErrInvalidParameter, "db.name")
	}
	return nil
}

type logConfigure struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Path   string `mapstructure:"path"`
}

func (c *logConfigure) GetLevel() string {
	return c.Level
}

func (c *logConfigure) GetFormat() string {
	return c.Format
}

func (c *logConfigure) GetOutput() string {
	return c.Output
}

func (c *logConfigure) GetPath() string {
	return c.Path
}

func (c *logConfigure) Validate() error {
	switch c.Level {
	case "", "info", "warn", "debug", "error", "fatal":
	default:
		return errors.Wrap(rserrors.ErrInvalidParameter, "logger.level")
	}

	switch c.Format {
	case "", "json", "text":
	default:
		return errors.Wrap(rserrors.ErrInvalidParameter, "logger.format")
	}

	switch c.Output {
	case "", "file", "console":
	default:
		return errors.Wrap(rserrors.ErrInvalidParameter, "logger.output")
	}

	return nil
}

func GetServerConfig() configure {
	return c
}
