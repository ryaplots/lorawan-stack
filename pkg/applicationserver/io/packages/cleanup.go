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

package packages

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/cleanup"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/unique"
)

type RegistryCleaner struct {
	ApplicationPackagesRegistry Registry
}

func (cleaner *RegistryCleaner) RangeToLocalSet(ctx context.Context) (localDevices map[string]struct{}, localApplications map[string]struct{}, err error) {
	localDevices = make(map[string]struct{})
	localApplications = make(map[string]struct{})
	err = cleaner.ApplicationPackagesRegistry.Range(ctx, []string{"ids"},
		func(ctx context.Context, ids ttnpb.EndDeviceIdentifiers, pb *ttnpb.ApplicationPackageAssociation) bool {
			localDevices[unique.ID(ctx, ids)] = struct{}{}
			return true
		},
		func(ctx context.Context, ids ttnpb.ApplicationIdentifiers, pb *ttnpb.ApplicationPackageDefaultAssociation) bool {
			localApplications[unique.ID(ctx, ids)] = struct{}{}
			return true
		},
	)
	return localDevices, localApplications, err
}

func (cleaner *RegistryCleaner) DeleteComplement(ctx context.Context, deviceSet map[string]struct{}, applicationSet map[string]struct{}) error {
	for ids := range deviceSet {
		devIds, err := unique.ToDeviceID(ids)
		if err != nil {
			return err
		}
		associations, err := cleaner.ApplicationPackagesRegistry.ListAssociations(ctx, devIds, []string{"ids"})
		if err != nil {
			return err
		}
		for _, association := range associations {
			_, err := cleaner.ApplicationPackagesRegistry.SetAssociation(ctx, association.ApplicationPackageAssociationIdentifiers, nil,
				func(assoc *ttnpb.ApplicationPackageAssociation) (*ttnpb.ApplicationPackageAssociation, []string, error) {
					return nil, nil, nil
				},
			)
			if err != nil {
				return err
			}
		}
	}
	for ids := range applicationSet {
		appIds, err := unique.ToApplicationID(ids)
		if err != nil {
			return err
		}
		associations, err := cleaner.ApplicationPackagesRegistry.ListDefaultAssociations(ctx, appIds, []string{"ids"})
		if err != nil {
			return err
		}
		for _, association := range associations {
			_, err := cleaner.ApplicationPackagesRegistry.SetDefaultAssociation(ctx, association.ApplicationPackageDefaultAssociationIdentifiers, nil,
				func(association *ttnpb.ApplicationPackageDefaultAssociation) (*ttnpb.ApplicationPackageDefaultAssociation, []string, error) {
					return nil, nil, nil
				},
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cleaner *RegistryCleaner) CleanData(ctx context.Context, isDeviceSet map[string]struct{}, isApplicationSet map[string]struct{}) error {
	localDeviceSet, localApplicationSet, err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return err
	}
	devComplement := cleanup.ComputeSetComplement(isDeviceSet, localDeviceSet)
	appComplement := cleanup.ComputeSetComplement(isApplicationSet, localApplicationSet)
	err = cleaner.DeleteComplement(ctx, devComplement, appComplement)
	if err != nil {
		return err
	}
	return nil
}
