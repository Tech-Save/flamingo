package flamingo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
)

type appmodule struct {
	root   *config.Area
	router *web.Router
	server *http.Server
	logger flamingo.Logger
}

// Inject basic application dependencies
func (a *appmodule) Inject(root *config.Area, router *web.Router, logger flamingo.Logger) {
	a.root = root
	a.router = router
	a.logger = logger
	a.server = &http.Server{
		Addr:    ":3322",
		Handler: a.router,
	}
}

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).ToInstance(a)

	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

func serveProvider(a *appmodule, logger flamingo.Logger) *cobra.Command {
	var addr string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			a.router.Init(a.root)
			logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", addr))
			err := a.server.ListenAndServe()
			if err != nil {
				if err == http.ErrServerClosed {
					logger.Error(err)
				} else {
					logger.Fatal("unexpected error in serving:", err)
				}
			}
		},
	}
	cmd.Flags().StringVarP(&a.server.Addr, "addr", "a", ":3322", "addr on which flamingo runs")

	return cmd
}

// Notify upon flamingo Shutdown event
func (a *appmodule) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		a.logger.Info("Shutdown server on ", a.server.Addr)

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Error("unexpected error on server shutdown: ", err)
		}
	}
}

type option func(config *appconfig)

// ConfigDir configuration option
func ConfigDir(configdir string) func(config *appconfig) {
	return func(config *appconfig) {
		config.configDir = configdir
	}
}

type appconfig struct {
	configDir string
}

// App is a simple app-runner for flamingo
func App(modules []dingo.Module, options ...option) {
	app := new(appmodule)
	root := config.NewArea("root", modules)

	root.Modules = append([]dingo.Module{
		new(framework.InitModule),
		new(zap.Module),
		new(cmd.Module),
	}, root.Modules...)

	root.Modules = append(root.Modules, app)
	cfg := &appconfig{
		configDir: "config",
	}
	for _, option := range options {
		option(cfg)
	}
	config.Load(root, cfg.configDir)

	cmd := root.Injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command)
	root.Injector.GetInstance(new(web.EventRouterProvider)).(web.EventRouterProvider)().Dispatch(context.Background(), new(flamingo.StartupEvent))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}