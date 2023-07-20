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

package filter

import (
	"errors"

	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/pkg/app"
	"github.com/illacloud/builder-backend/pkg/state"

	"github.com/illacloud/builder-backend/internal/util/builderoperation"
	ws "github.com/illacloud/builder-backend/internal/websocket"
)

func SignalPutState(hub *ws.Hub, message *ws.Message) error {
	// deserialize message
	currentClient, hit := hub.Clients[message.ClientID]
	if !hit {
		return errors.New("[SignalPutState] target client(" + message.ClientID.String() + ") does dot exists.")
	}
	stateType := repository.STATE_TYPE_INVALIED
	teamID := currentClient.TeamID
	appDto := app.NewAppDto()
	appDto.ConstructWithID(currentClient.APPID)
	appDto.ConstructWithUpdateBy(currentClient.MappedUserID)
	appDto.SetTeamID(currentClient.TeamID)
	message.RewriteBroadcast()

	// target switch
	switch message.Target {
	case builderoperation.TARGET_NOTNING:
		return nil
	case builderoperation.TARGET_COMPONENTS:
		return nil

	case builderoperation.TARGET_DEPENDENCIES:
		stateType = repository.KV_STATE_TYPE_DEPENDENCIES
		// delete all
		kvStateDto := state.NewKVStateDto()
		kvStateDto.InitUID()
		kvStateDto.SetTeamID(teamID)
		kvStateDto.ConstructByApp(appDto) // set AppRefID
		kvStateDto.ConstructWithType(stateType)
		if err := hub.KVStateServiceImpl.DeleteAllEditKVStateByStateType(kvStateDto); err != nil {
			return err
		}
		// create k-v state
		for _, v := range message.Payload {
			subv, ok := v.(map[string]interface{})
			if !ok {
				err := errors.New("K-V State reflect failed, please check your input.")
				return err
			}
			for key, depState := range subv {
				// fill KVStateDto
				kvStateDto := state.NewKVStateDto()
				kvStateDto.InitUID()
				kvStateDto.SetTeamID(teamID)
				kvStateDto.ConstructWithKey(key)
				kvStateDto.ConstructForDependenciesState(depState)
				kvStateDto.ConstructByApp(appDto) // set AppRefID
				kvStateDto.ConstructWithType(stateType)

				if _, err := hub.KVStateServiceImpl.CreateKVState(kvStateDto); err != nil {
					currentClient.Feedback(message, ws.ERROR_CREATE_STATE_FAILED, err)
					return err
				}
			}
		}
	case builderoperation.TARGET_DRAG_SHADOW:
		return nil

	case builderoperation.TARGET_DOTTED_LINE_SQUARE:
		return nil

	case builderoperation.TARGET_DISPLAY_NAME:
		return nil

	case builderoperation.TARGET_APPS:
		// serve on HTTP API, this signal only for broadcast
	case builderoperation.TARGET_RESOURCE:
		// serve on HTTP API, this signal only for broadcast
	}

	// the currentClient does not need feedback when operation success

	// change app modify time
	hub.AppServiceImpl.UpdateAppModifyTime(appDto)

	// feedback otherClient
	hub.BroadcastToOtherClients(message, currentClient)
	return nil
}
