package initialize

import (
	"github.com/gin-gonic/gin"
	"order-service/global"
	"order-service/internal/constants"
	"order-service/internal/middlewares"
	"order-service/internal/routers"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	r.Use(gin.Logger()) //logging
	r.Use(middlewares.CorsMiddleware())

	manage := routers.RouterGroupApp.Manage
	user := routers.RouterGroupApp.User

	MainGroup := r.Group(constants.API_ENPOINT)
	{
		manage.InitOrderRouter(MainGroup)
		user.InitOrderRouter(MainGroup)
	}

	return r
}
