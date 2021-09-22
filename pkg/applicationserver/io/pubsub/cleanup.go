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

package pubsub

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/cleanup"
	"go.thethings.network/lorawan-stack/v3/pkg/events"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

type RegistryCleaner struct {
	PubSubRegistry Registry
}

func (cleaner *RegistryCleaner) RangeToLocalSet(ctx context.Context) (local map[ttnpb.ApplicationIdentifiers]struct{}, err error) {
	local = make(map[ttnpb.ApplicationIdentifiers]struct{})
	err = cleaner.PubSubRegistry.Range(ctx, []string{"ids"},
		func(ctx context.Context, ids ttnpb.ApplicationIdentifiers, pb *ttnpb.ApplicationPubSub) bool {
			local[ids] = struct{}{}
			return true
		},
	)
	return local, err
}

func (cleaner *RegistryCleaner) DeleteComplement(ctx context.Context, applicationSet map[ttnpb.ApplicationIdentifiers]struct{}) error {
	for ids := range applicationSet {
		pubsubs, err := cleaner.PubSubRegistry.List(ctx, ids, []string{"ids"})
		if err != nil {
			return err
		}
		for _, pubsub := range pubsubs {
			_, err := cleaner.PubSubRegistry.Set(ctx, pubsub.ApplicationPubSubIdentifiers, nil,
				func(pubsub *ttnpb.ApplicationPubSub) (*ttnpb.ApplicationPubSub, []string, error) {
					return nil, nil, nil
				},
			)
			if err != nil {
				return err
			}
			events.Publish(evtDeletePubSub.NewWithIdentifiersAndData(ctx, &pubsub.ApplicationPubSubIdentifiers.ApplicationIdentifiers, pubsub.ApplicationPubSubIdentifiers))
		}
	}
	return nil
}

func (cleaner *RegistryCleaner) CleanData(ctx context.Context, isSet map[ttnpb.ApplicationIdentifiers]struct{}) error {
	localSet, err := cleaner.RangeToLocalSet(ctx)
	if err != nil {
		return err
	}
	complement := cleanup.ComputeApplicationSetComplement(isSet, localSet)
	err = cleaner.DeleteComplement(ctx, complement)
	if err != nil {
		return err
	}
	return nil
}
