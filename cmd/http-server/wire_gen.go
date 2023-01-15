// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/illacloud/builder-backend/api/resthandler"
	"github.com/illacloud/builder-backend/api/router"
	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/internal/util"
	"github.com/illacloud/builder-backend/internal/accesscontrol"
	"github.com/illacloud/builder-backend/pkg/action"
	"github.com/illacloud/builder-backend/pkg/builder"
	"github.com/illacloud/builder-backend/pkg/app"
	"github.com/illacloud/builder-backend/pkg/db"
	"github.com/illacloud/builder-backend/pkg/resource"
	"github.com/illacloud/builder-backend/pkg/room"
	"github.com/illacloud/builder-backend/pkg/state"
)

// Injectors from wire.go:

func Initialize() (*Server, error) {
	config, err := GetAppConfig()
	if err != nil {
		return nil, err
	}
	engine := gin.New()
	sugaredLogger := util.NewSugardLogger()
	dbConfig, err := db.GetConfig()
	if err != nil {
		return nil, err
	}
	gormDB, err := db.NewDbConnection(dbConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	
	// init supervisior
	attrg, err := accesscontrol.NewRawAttributeGroup()
	if err != nil {
		return nil, err
	}
	appRepositoryImpl := repository.NewAppRepositoryImpl(sugaredLogger, gormDB)
	kvStateRepositoryImpl := repository.NewKVStateRepositoryImpl(sugaredLogger, gormDB)
	treeStateRepositoryImpl := repository.NewTreeStateRepositoryImpl(sugaredLogger, gormDB)
	setStateRepositoryImpl := repository.NewSetStateRepositoryImpl(sugaredLogger, gormDB)
	actionRepositoryImpl := repository.NewActionRepositoryImpl(sugaredLogger, gormDB)	
	appServiceImpl := app.NewAppServiceImpl(sugaredLogger, appRepositoryImpl, kvStateRepositoryImpl, treeStateRepositoryImpl, setStateRepositoryImpl, actionRepositoryImpl)
	treeStateServiceImpl := state.NewTreeStateServiceImpl(sugaredLogger, treeStateRepositoryImpl)
	appRestHandlerImpl := resthandler.NewAppRestHandlerImpl(sugaredLogger, appServiceImpl, attrg, treeStateServiceImpl)
	appRouterImpl := router.NewAppRouterImpl(appRestHandlerImpl)
	// room
	roomServiceImpl := room.NewRoomServiceImpl(sugaredLogger)
	roomRestHandlerImpl := resthandler.NewRoomRestHandlerImpl(sugaredLogger, roomServiceImpl, attrg)
	roomRouterImpl := router.NewRoomRouterImpl(roomRestHandlerImpl)
	// resource
	resourceRepositoryImpl := repository.NewResourceRepositoryImpl(sugaredLogger, gormDB)
	resourceServiceImpl := resource.NewResourceServiceImpl(sugaredLogger, resourceRepositoryImpl)
	resourceRestHandlerImpl := resthandler.NewResourceRestHandlerImpl(sugaredLogger, resourceServiceImpl, attrg)
	resourceRouterImpl := router.NewResourceRouterImpl(resourceRestHandlerImpl)
	// actions
	actionServiceImpl := action.NewActionServiceImpl(sugaredLogger, appRepositoryImpl, actionRepositoryImpl, resourceRepositoryImpl)
	actionRestHandlerImpl := resthandler.NewActionRestHandlerImpl(sugaredLogger, actionServiceImpl, attrg)
	actionRouterImpl := router.NewActionRouterImpl(actionRestHandlerImpl)
	// builder
	builderServiceImpl := builder.NewBuilderServiceImpl(sugaredLogger, appRepositoryImpl, resourceRepositoryImpl, actionRepositoryImpl)
	builderRestHandlerImpl := resthandler.NewBuilderRestHandlerImpl(sugaredLogger, builderServiceImpl, attrg)
	builderRouterImpl := router.NewBuilderRouterImpl(builderRestHandlerImpl)
	restRouter := router.NewRESTRouter(sugaredLogger, builderRouterImpl, appRouterImpl, roomRouterImpl, actionRouterImpl, resourceRouterImpl)
	server := NewServer(config, engine, restRouter, sugaredLogger)
	return server, nil
}
