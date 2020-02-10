package cmd

import (
	"context"
	"os"

	"github.com/solo-io/wasme/pkg/defaults"

	"github.com/solo-io/wasme/pkg/cmd/tag"

	"github.com/solo-io/wasme/pkg/cmd/operator"

	"github.com/sirupsen/logrus"

	"github.com/solo-io/wasme/pkg/cmd/build"
	"github.com/solo-io/wasme/pkg/cmd/deploy"
	"github.com/solo-io/wasme/pkg/cmd/initialize"
	"github.com/solo-io/wasme/pkg/cmd/list"
	"github.com/solo-io/wasme/pkg/version"

	ctxo "github.com/deislabs/oras/pkg/context"
	"github.com/solo-io/wasme/pkg/cmd/cache"
	"github.com/solo-io/wasme/pkg/cmd/login"
	"github.com/solo-io/wasme/pkg/cmd/opts"
	"github.com/solo-io/wasme/pkg/cmd/pull"
	"github.com/solo-io/wasme/pkg/cmd/push"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var opts opts.AuthOptions

	ctx2 := context.Background()
	ctx := &ctx2
	cmd := &cobra.Command{
		Use:     "wasme [command]",
		Version: version.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if opts.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			} else if !opts.Verbose {
				ctx2 := ctxo.WithLoggerDiscarded(*ctx)
				*ctx = ctx2
			}
			// set default auth configs
			if len(opts.CredentialsFiles) == 0 {
				opts.CredentialsFiles = []string{defaults.WasmeCredentialsFile}
			}
		},
	}
	cmd.AddCommand(
		initialize.InitCmd(),
		build.BuildCmd(ctx),
		push.PushCmd(ctx, &opts),
		pull.PullCmd(ctx, &opts),
		cache.CacheCmd(ctx, &opts),
		login.LoginCmd(),
		list.ListCmd(),
		deploy.DeployCmd(ctx),
		deploy.UndeployCmd(ctx),
		operator.OperatorCmd(ctx),
		tag.TagCmd(ctx),
	)
	cmd.PersistentFlags().StringArrayVarP(&opts.CredentialsFiles, "config", "c", nil, "path to auth config")
	cmd.PersistentFlags().StringVarP(&opts.Username, "username", "u", "", "registry username")
	cmd.PersistentFlags().StringVarP(&opts.Password, "password", "p", "", "registry password")
	cmd.PersistentFlags().BoolVarP(&opts.Insecure, "insecure", "", false, "allow connections to SSL registry without certs")
	cmd.PersistentFlags().BoolVarP(&opts.PlainHTTP, "plain-http", "", false, "use plain http and not https")
	cmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVarP(&opts.Debug, "debug", "d", false, "debug mode")
	return cmd
}

func Run() {
	if err := Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
