package main

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"os"
)

type (
	Store interface {
		Get(string) string
	}

	MysqlStore struct{}

	Service struct {
		store Store
	}
)

func (s *MysqlStore) Get(data string) string {
	return "Mysql: " + data
}

func NewMysqlStore() Store {
	return &MysqlStore{}
}

func NewService(
	store Store,
) *Service {
	return &Service{
		store: store,
	}
}

var ServiceModule = fx.Module("service",
	fx.Provide(NewMysqlStore),
	fx.Provide(NewService),
	fx.Invoke(func(s *Service) {
		rs := s.store.Get("hello")
		fmt.Println(rs)
	}),
)

func main() {
	app := fx.New(ServiceModule)

	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Error starting the application: %v\n", err)
		os.Exit(1)
	}

	<-app.Done()
}
