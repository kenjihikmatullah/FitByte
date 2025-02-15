package httpServer

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rafitanujaya/go-fiber-template/src/config"
	"github.com/rafitanujaya/go-fiber-template/src/di"
	activityController "github.com/rafitanujaya/go-fiber-template/src/http/controllers/activity"
	userController "github.com/rafitanujaya/go-fiber-template/src/http/controllers/user"
	"github.com/rafitanujaya/go-fiber-template/src/http/routes"
	activityroutes "github.com/rafitanujaya/go-fiber-template/src/http/routes/activity"
	userroutes "github.com/rafitanujaya/go-fiber-template/src/http/routes/user"
	"github.com/samber/do/v2"
)

type HttpServer struct{}

func (s *HttpServer) Listen() {
	fmt.Printf("New Fiber\n")
	app := fiber.New(fiber.Config{
		ServerHeader: "TIM-DEBUG",
	})

	fmt.Printf("Inject Controllers\n")
	//? Depedency Injection
	//? UserController
	uc := do.MustInvoke[userController.UserControllerInterface](di.Injector)
	ac := do.MustInvoke[activityController.ActivityControllerInterface](di.Injector)

	routes := routes.SetRoutes(app)
	userroutes.SetRouteUsers(routes, uc)
	activityroutes.SetRouteActivities(routes, ac)

	fmt.Printf("Start Lister\n")
	app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", config.GetPort()))
}
