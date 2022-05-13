package internal

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cqroot/s3-active-exporter/logger"
)

var (
	defaultListenAddress = ":9180"
	defaultTelemetryPath = "/metrics"
	defaultMaxRequests   = 30
	defaultLogDebug      = false
	defaultLogVerbose    = false

	listenAddress  = flag.String("web.listen-address", defaultListenAddress, "Address on which to expose metrics and web interface.")
	metricPath     = flag.String("web.telemetry-path", defaultTelemetryPath, "Path under which to expose metrics.")
	maxRequests    = flag.Int("web.max-requests", defaultMaxRequests, "Maximum number of parallel scrape requests. Use 0 to disable.")
	debug          = flag.Bool("log.debug", defaultLogDebug, "Output debug information.")
	verbose        = flag.Bool("log.verbose", defaultLogVerbose, "Output file name and line number.")
	config         = flag.StringP("config", "c", ".", "Specify the configuration file.")
	defaultFilters = []string{"server"}
)

func InitConfig() {
	// Viper set default
	viper.SetDefault("web.listen-address", defaultListenAddress)
	viper.SetDefault("web.telemetry-path", defaultTelemetryPath)
	viper.SetDefault("web.max-requests", defaultMaxRequests)
	viper.SetDefault("log.debug", defaultLogDebug)
	viper.SetDefault("log.verbose", defaultLogVerbose)

	// Pflag parse
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	// Viper read config
	if *config != "." {
		viper.SetConfigFile(*config)
	} else {
		viper.SetConfigName("s3-active-exporter")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/s3-active-exporter/")
	}
	err := viper.ReadInConfig()

	// init logger
	logger.Init(viper.GetBool("log.debug"), viper.GetBool("log.verbose"))
	if viper.GetBool("log.debug") {
		logger.Debug("Enabling debug output")
	}

	if err != nil {
		logger.Warn(err.Error())
	}
}
