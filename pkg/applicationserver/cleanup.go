// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

package applicationserver

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/cleanup"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

type RegistryCleaner struct {
	DevRegistry    DeviceRegistry
	AppUpsRegistry ApplicationUplinkRegistry
}

func (cleaner *RegistryCleaner) RangeToLocalSet(ctx context.Context) (local map[ttnpb.EndDeviceIdentifiers]struct{}, err error) {
	local = make(map[ttnpb.EndDeviceIdentifiers]struct{})
	err = cleaner.DevRegistry.Range(ctx, []string{"ids"}, func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, dev *ttnpb.EndDevice) bool {
		local[ids] = struct{}{}
		return true
	},
	)
	return local, err
}

func (cleaner *RegistryCleaner) DeleteComplement(ctx context.Context, devSet map[ttnpb.EndDeviceIdentifiers]struct{}) error {
	for ids := range devSet {
		_, err := cleaner.DevRegistry.Set(ctx, ids, nil, func(dev *ttnpb.EndDevice) (*ttnpb.EndDevice, []string, error) {
			return nil, nil, nil
		})
		if err != nil {
			return err
		}
		if err := cleaner.AppUpsRegistry.Clear(ctx, ids); err != nil {
			return err
		}
	}
	return nil
}

func (cleaner *RegistryCleaner) CleanData(ctx context.Context, isSet map[ttnpb.EndDeviceIdentifiers]struct{}) error {
	localSet, err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return err
	}
	complement := cleanup.ComputeDeviceSetComplement(isSet, localSet)
	err = cleaner.DeleteComplement(ctx, complement)
	if err != nil {
		return err
	}
	return nil
}
