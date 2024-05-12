package internal

import (
	"APIGateWay/config"
	"APIGateWay/internal/proxy/usecase"
	"APIGateWay/internal/service"
	service_repo "APIGateWay/internal/service/repository"
	"APIGateWay/pkg/db"
	"APIGateWay/pkg/guzzle_logger"
	"APIGateWay/pkg/rabbit/publisher"
	"github.com/jmoiron/sqlx"
)

type App struct {
	cfg          *config.Config
	UC           map[string]interface{}
	Repo         map[string]interface{}
	API          map[string]interface{}
	dbConnection map[string]*sqlx.DB
	GuzzleLogger guzzle_logger.API
}

func NewApp(cfg *config.Config) *App {
	logPublisher, err := publisher.NewPublisher(cfg.Rabbit.LogPublisher.Url, cfg.Rabbit.LogPublisher.QueueName)
	if err != nil {
		panic(err)
	}
	guzzleLogger := guzzle_logger.New("gateway", "Proxy", logPublisher)
	return &App{
		UC:           make(map[string]interface{}),
		Repo:         make(map[string]interface{}),
		dbConnection: make(map[string]*sqlx.DB),
		API:          make(map[string]interface{}),
		cfg:          cfg,
		GuzzleLogger: guzzleLogger,
	}
}

func (a *App) Init() error {
	var err error
	a.dbConnection["postgres"], err = db.InitPsqlDB(a.cfg)
	if err != nil {
		return err
	}
	a.Repo["service"] = service_repo.NewPostgresRepository(a.dbConnection["postgres"])

	a.UC["proxy"] = proxy_uc.NewProxyUC(a.Repo["service"].(service.Repository), a.GuzzleLogger)
	return nil
}
