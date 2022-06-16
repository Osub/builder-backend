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

package resource

import (
	"github.com/illa-family/builder-backend/internal/repository"
	"go.uber.org/zap"
)

type ResourceService interface {
	CreateResource(resource ResourceDto) (ResourceDto, error)
	DeleteResource(resourceId string) error
	UpdateResource(resource ResourceDto) (ResourceDto, error)
}

type ResourceDto struct {
	id string
}

type ResourceServiceImpl struct {
	logger             *zap.SugaredLogger
	resourceRepository repository.ResourceRepository
}

func NewResourceServiceImpl(logger *zap.SugaredLogger, resourceRepository repository.ResourceRepository) *ResourceServiceImpl {
	return &ResourceServiceImpl{
		logger:             logger,
		resourceRepository: resourceRepository,
	}
}

func (impl *ResourceServiceImpl) CreateResource(resource ResourceDto) (ResourceDto, error) {
	return ResourceDto{}, nil
}

func (impl *ResourceServiceImpl) DeleteResource(resourceId string) error {
	return nil
}

func (impl *ResourceServiceImpl) UpdateResource(resource ResourceDto) (ResourceDto, error) {
	return ResourceDto{}, nil
}
