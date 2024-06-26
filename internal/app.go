package internal

import (
	"context"
	"notifications/internal/io/storage"

	"notifications/internal/handler"
	"notifications/internal/io/transports"
	"notifications/internal/io/web"
)

type App struct {
	ctx context.Context
	cfg *Config
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,
		ctx: ctx,
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error)

	tCollection, tCollectionInitErr := transports.NewCollection(a.cfg.TransportsCfgFilePath)
	if tCollectionInitErr != nil {
		return tCollectionInitErr
	}

	store, storeErr := storage.Factory(a.cfg.StorageDsn)
	if storeErr != nil {
		return storeErr
	}

	executor := handler.NewHandler(store, tCollection)

	httpSrv := web.NewServer(&a.cfg.Web, executor.Handle)

	httpSrv.Run(errCh)

	select {
	case <-a.ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
