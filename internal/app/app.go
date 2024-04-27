package app

import "auth-service/internal/dependencies"

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
	appServer := app.deps.AppServer()
	appServer.Start()

	app.deps.WaitForInterrupr() // программа будет "стоять" тут пока не придет системный сигнал
	app.deps.Close()
}
