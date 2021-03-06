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
)

func Start(cfg *config.Config, lgr *zap.Logger) {
	dbClient, err := initializers.InitDBClient(cfg.DBConfig, lgr)
	if err != nil {
		lgr.Fatal(fmt.Sprintf("[App] [Start] [InitDBClient] %v", err))
	}
	db := database.NewDBInstance(dbClient, lgr)
	svc := service.NewDBService(db, lgr)
	rtr := router.NewRouter(svc, cfg.JWT, lgr)

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

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
