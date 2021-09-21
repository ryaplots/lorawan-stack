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

package web

import (
	"context"

	"go.thethings.network/lorawan-stack/v3/pkg/cleanup"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

type RegistryCleaner struct {
	WebRegistry WebhookRegistry
}

func (cleaner *RegistryCleaner) RangeToLocalSet(ctx context.Context) (local map[*ttnpb.ApplicationIdentifiers]struct{}, err error) {
	err = cleaner.WebRegistry.Range(ctx, []string{"ids"},
		func(ctx context.Context, ids ttnpb.ApplicationIdentifiers, wh *ttnpb.ApplicationWebhook) bool {
			local[&ids] = struct{}{}
			return true
		},
	)
	return local, err
}

func (cleaner *RegistryCleaner) DeleteComplement(ctx context.Context, applicationSet map[*ttnpb.ApplicationIdentifiers]struct{}) error {
	for ids := range applicationSet {
		webhooks, err := cleaner.WebRegistry.List(ctx, *ids, []string{"ids"})
		if err != nil {
			return err
		}
		for _, webhook := range webhooks {
			_, err := cleaner.WebRegistry.Set(ctx, webhook.ApplicationWebhookIdentifiers, nil,
				func(webhook *ttnpb.ApplicationWebhook) (*ttnpb.ApplicationWebhook, []string, error) {
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

func (cleaner *RegistryCleaner) CleanData(ctx context.Context, isSet map[*ttnpb.ApplicationIdentifiers]struct{}) error {
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
