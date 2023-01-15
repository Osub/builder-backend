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
	"github.com/illacloud/builder-backend/pkg/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RESTRouter struct {
	logger         *zap.SugaredLogger
	Router         *gin.RouterGroup
	BuilderRouter  BuilderRouter
	AppRouter      AppRouter
	RoomRouter     RoomRouter
	ActionRouter   ActionRouter
	ResourceRouter ResourceRouter
}

func NewRESTRouter(logger *zap.SugaredLogger, builderRouter BuilderRouter, appRouter AppRouter, roomRouter RoomRouter,
	actionRouter ActionRouter, resourceRouter ResourceRouter) *RESTRouter {
	return &RESTRouter{
		logger:         logger,
		BuilderRouter:  builderRouter,
		AppRouter:      appRouter,
		RoomRouter:     roomRouter,
		ActionRouter:   actionRouter,
		ResourceRouter: resourceRouter,
	}
}

func (r RESTRouter) InitRouter(router *gin.RouterGroup) {
	v1 := router.Group("/v1")

	builderRouter := v1.Group("/teams/:teamID/builder")
	appRouter := v1.Group("/teams/:teamID/apps")
	resourceRouter := v1.Group("/teams/:teamID/resources")
	actionRouter := v1.Group("/teams/:teamID/actions")
	roomRouter := v1.Group("/teams/:teamID/room")

	builderRouter.Use(user.RemoteJWTAuth())
	appRouter.Use(user.RemoteJWTAuth())
	roomRouter.Use(user.RemoteJWTAuth())
	actionRouter.Use(user.RemoteJWTAuth())
	resourceRouter.Use(user.RemoteJWTAuth())

	r.BuilderRouter.InitBuilderRouter(builderRouter)
	r.AppRouter.InitAppRouter(appRouter)
	r.RoomRouter.InitRoomRouter(roomRouter)
	r.ActionRouter.InitActionRouter(actionRouter)
	r.ResourceRouter.InitResourceRouter(resourceRouter)
}
