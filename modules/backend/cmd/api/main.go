package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mstrYoda/maasanketi.co/handler"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/mstrYoda/maasanketi.co/swagger" // Import swagger docs
)

type Application struct {
	app  *fiber.App
	repo *repository.Repository
}

func (a *Application) Register() {
	// Health and monitoring
	a.app.Get("/healthcheck", handler.HealthCheck)
	a.app.Get("/monitor", monitor.New())
	a.app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// Swagger
	a.app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1 routes
	v1 := a.app.Group("/api/v1")

	// Survey routes
	surveys := v1.Group("/surveys")
	surveys.Get("/", handler.GetSurveys)
	surveys.Get("/:id", handler.GetSurvey)
}

// @title Software Engineer Salary Survey API
// @version 1.0
// @description API for the Software Engineer Salary Survey platform
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.maasanketi.co/support
// @contact.email support@maasanketi.co
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:4004
// @BasePath /api/v1
// @schemes http https
func main() {
	repo, err := repository.New()
	if err != nil {
		log.Logger().Panic(fmt.Sprintf("failed to create repository: %s", err.Error()))
	} else {
		log.Logger().Info("repository created successfully")
	}
	defer repo.Close()

	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		JSONDecoder: jsoniter.Unmarshal,
	})
	app.Use(cors.New())
	app.Use(recover.New())

	application := &Application{app, repo}
	application.Register()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		log.Logger().Info("application gracefully shutting down..")
		_ = app.Shutdown()
	}()

	if err := app.Listen("0.0.0.0:4004"); err != nil {
		log.Logger().Panic(fmt.Sprintf("app error: %s", err.Error()))
	}
}
