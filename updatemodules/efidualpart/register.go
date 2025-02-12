// SPDX-License-Identifier: Apache-2.0
//
// Copyright (C) 2021 Renesas Electronics Corporation.
// Copyright (C) 2021 EPAM Systems, Inc.
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

package efidualpart

import (
	"encoding/json"

	"github.com/aoscloud/aos_common/aoserrors"

	"github.com/aoscloud/aos_updatemanager/updatehandler"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/controllers/eficontroller"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/modules/dualpartmodule"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/rebooters/systemdrebooter"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/updatechecker/systemdchecker"
)

/*******************************************************************************
 * Types
 ******************************************************************************/

type moduleConfig struct {
	Loader         string                `json:"loader"`
	VersionFile    string                `json:"versionFile"`
	Partitions     []string              `json:"partitions"`
	SystemdChecker systemdchecker.Config `json:"systemdChecker"`
}

/*******************************************************************************
 * Init
 ******************************************************************************/

func init() {
	updatehandler.RegisterPlugin("efidualpart",
		func(id string, configJSON json.RawMessage,
			storage updatehandler.ModuleStorage,
		) (module updatehandler.UpdateModule, err error) {
			if len(configJSON) == 0 {
				return nil, aoserrors.Errorf("config for %s module is required", id)
			}

			var config moduleConfig

			if err = json.Unmarshal(configJSON, &config); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			controller, err := eficontroller.New(config.Partitions, config.Loader)
			if err != nil {
				return nil, aoserrors.Wrap(err)
			}

			if module, err = dualpartmodule.New(id, config.Partitions, config.VersionFile,
				controller, storage, &systemdrebooter.SystemdRebooter{}, systemdchecker.New(config.SystemdChecker)); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			return module, nil
		},
	)
}
