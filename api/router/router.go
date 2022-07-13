// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"github.com/illa-family/builder-backend/pkg/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RESTRouter struct {
	logger         *zap.SugaredLogger
	Router         *gin.RouterGroup
	UserRouter     UserRouter
	AppRouter      AppRouter
	ActionRouter   ActionRouter
	ResourceRouter ResourceRouter
}

func NewRESTRouter(logger *zap.SugaredLogger, userRouter UserRouter, appRouter AppRouter,
	actionRouter ActionRouter, resourceRouter ResourceRouter) *RESTRouter {
	return &RESTRouter{
		logger:         logger,
		UserRouter:     userRouter,
		AppRouter:      appRouter,
		ActionRouter:   actionRouter,
		ResourceRouter: resourceRouter,
	}
}

func (r RESTRouter) InitRouter(router *gin.RouterGroup) {
	v1 := router.Group("/v1")

	authRouter := v1.Group("/auth")
	userRouter := v1.Group("/users")
	appRouter := v1.Group("/apps")
	actionRouter := v1.Group("/apps/:app/versions/:version")
	resourceRouter := v1.Group("/resources")

	userRouter.Use(user.JWTAuth())
	appRouter.Use(user.JWTAuth())
	actionRouter.Use(user.JWTAuth())
	resourceRouter.Use(user.JWTAuth())

	r.UserRouter.InitAuthRouter(authRouter)
	r.UserRouter.InitUserRouter(userRouter)
	r.AppRouter.InitAppRouter(appRouter)
	r.ActionRouter.InitActionRouter(actionRouter)
	r.ResourceRouter.InitResourceRouter(resourceRouter)
}
