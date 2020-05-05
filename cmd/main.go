package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/bamboo/internal/bamboo"
	"github.com/IgorRybak2055/bamboo/internal/storage"
	"github.com/IgorRybak2055/bamboo/pkg/config"
	"github.com/IgorRybak2055/bamboo/pkg/services"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	app := bamboo.New(cfg.HTTP)

	// TODO: we don't need somethi
	srvs := services.NewServices(log)

	srvs.Run(2, 10, func() error {
		return cfg.DB.MakeMigrations(log)
	})

	srvs.Run(10, 10, func() error {
		conn, err := storage.Connect(cfg.DB.Postgres(), log)
		if err != nil {
			log.Error(err)
			return err
		}
		app.DBC = conn
		log.Info("database connection established")

		go func() {
			<-ctx.Done()
			log.Println("started gracefully closing db connection")

			if err := conn.Close(); err != nil {
				log.Error(err)
			}
		}()
		return nil
	})

	srvs.Run(2, 10, func() error {
		return app.Start()
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	<-stop

	if err := app.Srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutting down server %s", err)
	} else {
		log.Println("server gracefully stopped")
	}
}
