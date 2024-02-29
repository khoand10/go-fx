package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net"
	"net/http"
	"os"
)

type (
	Connection string
	Store      interface {
		Get(string) string
	}

	MysqlStore struct {
		connection Connection
	}

	Service struct {
		store Store
	}
)

func NewConnection() Connection {
	return "mysql"
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
	return r
}

func (s *MysqlStore) Get(data string) string {
	return "Mysql: " + data
}

func NewMysqlStore(con Connection) Store {
	return &MysqlStore{
		connection: con,
	}
}

func NewService(
	store Store,
) *Service {
	return &Service{
		store: store,
	}
}

func NewHTTPServer(lc fx.Lifecycle, gin *gin.Engine) *http.Server {
	srv := &http.Server{Addr: ":8080", Handler: gin}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

var ServiceModule = fx.Module("service",
	fx.Provide(NewConnection),
	fx.Provide(NewMysqlStore),
	fx.Provide(NewService),
	fx.Provide(NewHTTPServer),
	fx.Provide(NewRouter),
	fx.Invoke(func(s *Service) {
		rs := s.store.Get("hello")
		fmt.Println(rs)
	}),
	fx.Invoke(func(server *http.Server) {}),
)

func main() {
	app := fx.New(ServiceModule)

	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Error starting the application: %v\n", err)
		os.Exit(1)
	}

	<-app.Done()
}
