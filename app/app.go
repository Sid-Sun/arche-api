package app

import (
	"fmt"
	"github.com/sid-sun/arche-api/app/database"
	"github.com/sid-sun/arche-api/app/initializers"
	"github.com/sid-sun/arche-api/app/router"
	"github.com/sid-sun/arche-api/app/service"
	"github.com/sid-sun/arche-api/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Start(cfg *config.Config, lgr *zap.Logger) {
	dbClient, err := initializers.InitDBClient(cfg.DBConfig, lgr)
	if err != nil {
		lgr.Fatal(fmt.Sprintf("[App] [Start] [InitDBClient] %v", err))
	}
	db := database.NewDBInstance(dbClient, lgr)

	mc := initializers.InitMGClient(cfg.EmailConfig)

	svc := service.NewService(db, mc, lgr)
	rtr := router.NewRouter(svc, cfg.JWT, cfg.VECfg, lgr)

	srv := &http.Server{
		Addr:    cfg.HTTP.GetListenAddr(),
		Handler: rtr,
	}

	lgr.Info(fmt.Sprintf("[App] [Start] [ListenAndServe] Going to listening on http://%s", cfg.HTTP.GetListenAddr()))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lgr.Error(fmt.Sprintf("[App] [Start] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
