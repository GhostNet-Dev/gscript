package commands

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/GhostNet-Dev/gscript/internal/gconfig"
	"github.com/GhostNet-Dev/gscript/repl"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	defaultCfg = gconfig.NewDefaultConfig()
)

// RootCmd root command binding
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gscript",
		Short: "PERL (read, Eval, Print, Loop)",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			user, err := user.Current()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Hello %s! This is the Ghost Lambda Language!\n", user.Username)
			fmt.Printf("Feel free to type in commands\n")
			repl.Start(os.Stdin, os.Stdout)
		},
	}

	return cmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(defaultCfg.DefaultConfigFilename)
	v.AddConfigPath(".")
	v.AddConfigPath("../")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(defaultCfg.EnvPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", defaultCfg.EnvPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
