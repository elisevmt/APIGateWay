package repository

import (
	"APIGateWay/internal/service"
	"errors"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) service.Repository {
	return &PostgresRepository{
		db: db,
	}
}

func (p *PostgresRepository) GetServiceProxyInstanceById(serviceId *int64) (*service.ServiceProxyInstance, error) {
	var data []service.ServiceProxyInstance
	err := p.db.Select(&data, "select * from service inner join instance i on service.id = i.service_id inner join proxy p on p.id = service.proxy_id where service_id=$1", *serviceId)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("no such service proxy instance")
	}
	return &data[0], nil
}

func (p *PostgresRepository) GetServiceProxyInstance() ([]*service.ServiceProxyInstance, error) {
	var data []*service.ServiceProxyInstance
	err := p.db.Select(&data, "select * from service inner join instance i on service.id = i.service_id inner join proxy p on p.id = service.proxy_id")
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("no such service proxy instance")
	}
	return data, nil
}
