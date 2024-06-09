package app

import "github.com/warehouse/auth-service/internal/dependencies"

type (
	Application interface {
		Run()
	}

	application struct {
		deps dependencies.Dependencies
	}
)

func NewApplication(cfgPath string) Application {
	deps, err := dependencies.NewDependencies(cfgPath)
	if err != nil {
		panic(err)
	}

	return &application{
		deps: deps,
	}
}

func (app *application) Run() {
	httpServer := app.deps.HttpServer()
	httpServer.Start()

	grpcServer := app.deps.GrpcServer()
	grpcServer.Start()

	app.deps.WaitForInterrupr() // программа будет "стоять" тут пока не придет системный сигнал
	app.deps.Close()
}
