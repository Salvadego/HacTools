package options

import (
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Address  string
	User     string
	Password string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetDefaults(cmd *cobra.Command, conf *Config) {
	defaultAddress := getEnvOrDefault("HYBRIS_HAC_URL", "https://localhost:9002/hac")
	defaultUser := getEnvOrDefault("HYBRIS_USER", "admin")
	defaultPassword := getEnvOrDefault("HYBRIS_PASSWORD", "nimda")

	cmd.PersistentFlags().StringVarP(&conf.Address, "address", "s", defaultAddress, "HAC address (default: $HYBRIS_HAC_URL)")
	cmd.PersistentFlags().StringVarP(&conf.User, "user", "u", defaultUser, "Username for HAC (default: $HYBRIS_USER)")
	cmd.PersistentFlags().StringVarP(&conf.Password, "password", "p", defaultPassword, "Password for HAC (default: $HYBRIS_PASSWORD)")
}
