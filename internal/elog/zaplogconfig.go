package zaplog

const (
	LogPath = ""
)

type (
	zaplogConfig struct {
		logPath    string
		isDevelop  bool
		maxSize    int
		maxAge     int
		maxBackups int
		compress   bool
		logApiPath string
		listenAddr string
		service    string
		instance   string
	}

	ConfigOption interface {
		apply(config *zaplogConfig)
	}
)

var defaultConfig = &zaplogConfig{
	logPath:    LogPath,
	isDevelop:  false,
	maxSize:    100,
	maxAge:     7,
	compress:   false,
	logApiPath: "/log",
	listenAddr: "0.0.0.0:0",
}

type zapLogOptionFunc func(config *zaplogConfig)

func (f zapLogOptionFunc) apply(config *zaplogConfig) {
	f(config)
}

func SetLogPath(path string) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.logPath = path
	})

}

func SetIsDevelop(flag bool) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.isDevelop = flag
	})
}

func SetMaxSize(size int) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.maxSize = size
	})
}

func SetMaxAge(age int) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.maxAge = age
	})
}

func SetMaxBackups(backups int) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.maxBackups = backups
	})
}

func SetCompress(compress bool) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.compress = compress
	})
}

func SetLogApiPath(logApiPath string) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.logApiPath = logApiPath
	})
}

func SetListenAddr(listenAddr string) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.listenAddr = listenAddr
	})
}

func SetService(service string) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.service = service
	})
}

func SetInstance(instance string) ConfigOption {
	return zapLogOptionFunc(func(config *zaplogConfig) {
		config.instance = instance
	})
}
