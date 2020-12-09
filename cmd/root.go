package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/make-os/kit/config/chains"
	"github.com/make-os/kit/pkgs/logger"
	"github.com/make-os/kit/util"
	"github.com/make-os/kit/util/colorfmt"
	"github.com/pkg/profile"
	"github.com/thoas/go-funk"

	"github.com/make-os/kit/config"
	tmcfg "github.com/tendermint/tendermint/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// BuildVersion is the build version set by goreleaser
	BuildVersion = ""

	// BuildCommit is the git hash of the build. It is set by goreleaser
	BuildCommit = ""

	// BuildDate is the date the build was created. Its is set by goreleaser
	BuildDate = ""

	// GoVersion is the version of go used to build the client
	GoVersion = ""
)

var (
	log logger.Logger

	// cfg is the application config
	cfg = config.EmptyAppConfig()

	// Get a reference to tendermint's config object
	tmconfig = tmcfg.DefaultConfig()

	// itr is used to inform the stoppage of all modules
	itr = util.Interrupt(make(chan struct{}))

	profiler interface{ Stop() }
)

// Execute the root command or fallback command when command is unknown.
func Execute() {

	// When command is unknown, run the root command PersistentPreRun
	// then run the fallback command
	_, _, err := rootCmd.Find(os.Args[1:])
	if err != nil && strings.Index(err.Error(), "unknown command") != -1 {
		rootCmd.PersistentPreRun(fallbackCmd, os.Args)
		fallbackCmd.Run(&cobra.Command{}, os.Args)
		return
	}

	// Stop the profiler if is running
	defer func() {
		if profiler != nil {
			profiler.Stop()
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func setupProfiler(rootCmd *cobra.Command, cfg *config.AppConfig) {
	profileMode, _ := rootCmd.PersistentFlags().GetString("profile.mode")
	switch profileMode {
	case "cpu":
		profiler = profile.Start(profile.CPUProfile, profile.ProfilePath(cfg.DataDir()))
	case "mem":
		profiler = profile.Start(profile.MemProfile, profile.ProfilePath(cfg.DataDir()))
	case "mutex":
		profiler = profile.Start(profile.MutexProfile, profile.ProfilePath(cfg.DataDir()))
	case "block":
		profiler = profile.Start(profile.BlockProfile, profile.ProfilePath(cfg.DataDir()))
	}
}

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   config.AppName,
	Short: "Kit is the official client for the MakeOS network",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// Set version information
		setVersionInfo()

		// Run pre-run routine if current called command is not in the pre-run ignore list
		preRunIgnoreList := []string{cmd.Root().Name()}
		if !funk.ContainsString(preRunIgnoreList, cmd.CalledAs()) {
			preRun(cmd)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Println("Client:", BuildVersion)
			fmt.Println("Build:", BuildCommit)
			fmt.Println("Go:", GoVersion)
			if cfg.G().NodeKey != nil {
				fmt.Println("NodeID:", cfg.G().NodeKey.ID())
			}
			return
		}
		_ = cmd.Help()
	},
}

func preRun(cmd *cobra.Command) {

	// Override net.version if --v1 network preset flag is provided in an `init` call.
	isInit := cmd.CalledAs() == "init"
	if isInit {
		if v1Flag := cmd.Flags().Lookup("v1"); v1Flag != nil && v1Flag.Changed {
			viper.Set("net.version", chains.TestnetV1.NetVersion)
		}
	}

	// Configure the node
	config.Configure(cfg, tmconfig, isInit, &itr)
	log = cfg.G().Log

	// Setup the profiler
	setupProfiler(cmd.Root(), cfg)

	// Load keys in the config object
	if !isInit {
		cfg.LoadKeys(tmconfig.NodeKeyFile(), tmconfig.PrivValidatorKeyFile(), tmconfig.PrivValidatorStateFile())
	}

	// Skip git exec check for certain commands
	if !funk.ContainsString([]string{"init", "start", "console", "sign", "attach", "config"}, cmd.CalledAs()) {
		return
	}

	// Verify git version compliance
	if yes, version := util.IsGitInstalled(cfg.Node.GitBinPath); yes {
		if semver.New(version).LessThan(*semver.New("2.11.0")) {
			log.Fatal(colorfmt.YellowStringf(`Git version is outdated. Please update git executable.` +
				`Visit https://git-scm.com/downloads to download and install the latest version.`,
			))
		}
	} else {
		log.Fatal(colorfmt.YellowStringf(`Git executable was not found.` +
			`If you already have Git installed, provide the executable's location using --gitpath, otherwise ` +
			`visit https://git-scm.com/downloads to download and install it.`,
		))
	}
}

func setVersionInfo() {
	cfg.VersionInfo = &config.VersionInfo{
		BuildCommit:  BuildCommit,
		BuildDate:    BuildDate,
		GoVersion:    GoVersion,
		BuildVersion: BuildVersion,
	}
}

// fallbackCmd is called any time an unknown command is executed
var fallbackCmd = &cobra.Command{
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Unknown command. Use --help to see commands.\n")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(fallbackCmd)
	rootCmd.Flags().SortFlags = false

	// Register flags
	rootCmd.PersistentFlags().String("home", config.DefaultDataDir, "Set the path to the home directory")
	rootCmd.PersistentFlags().String("home.prefix", "", "Adds a prefix to the home directory in dev mode")
	rootCmd.PersistentFlags().String("gitpath", "git", "Set path to git executable")
	rootCmd.PersistentFlags().Bool("dev", false, "Enables development mode")
	rootCmd.PersistentFlags().Uint64("net", config.DefaultNetVersion, "Set network/chain ID")
	rootCmd.PersistentFlags().Bool("no-log", false, "Disables loggers")
	rootCmd.PersistentFlags().Bool("no-colors", false, "Disables output colors")
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
	rootCmd.PersistentFlags().StringToString("loglevel", map[string]string{}, "Set log level for modules")
	rootCmd.PersistentFlags().String("profile.mode", "", "Enable profiling mode, one of [cpu, mem, mutex, block]")

	// Remote API connection flags
	rootCmd.PersistentFlags().String("rpc.user", "", "Set the RPC username")
	rootCmd.PersistentFlags().String("rpc.password", "", "Set the RPC password")
	rootCmd.PersistentFlags().Bool("rpc.https", false, "Force the client to use https:// protocol")
	rootCmd.PersistentFlags().String("remote.address", config.DefaultRemoteServerAddress, "Set the RPC server address")
	rootCmd.PersistentFlags().String("remote", "origin", "Set the default remote name")

	// Viper bindings
	_ = viper.BindPFlag("node.gitpath", rootCmd.PersistentFlags().Lookup("gitpath"))
	_ = viper.BindPFlag("net.version", rootCmd.PersistentFlags().Lookup("net"))
	_ = viper.BindPFlag("dev", rootCmd.PersistentFlags().Lookup("dev"))
	_ = viper.BindPFlag("home", rootCmd.PersistentFlags().Lookup("home"))
	_ = viper.BindPFlag("home.prefix", rootCmd.PersistentFlags().Lookup("home.prefix"))
	_ = viper.BindPFlag("no-log", rootCmd.PersistentFlags().Lookup("no-log"))
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	_ = viper.BindPFlag("no-colors", rootCmd.PersistentFlags().Lookup("no-colors"))
	_ = viper.BindPFlag("rpc.user", rootCmd.PersistentFlags().Lookup("rpc.user"))
	_ = viper.BindPFlag("rpc.password", rootCmd.PersistentFlags().Lookup("rpc.password"))
	_ = viper.BindPFlag("remote.address", rootCmd.PersistentFlags().Lookup("remote.address"))
	_ = viper.BindPFlag("rpc.https", rootCmd.PersistentFlags().Lookup("rpc.https"))
	_ = viper.BindPFlag("remote.name", rootCmd.PersistentFlags().Lookup("remote"))
}
