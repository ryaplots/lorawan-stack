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

package cleanup

import (
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

func ComputeApplicationSetComplement(isSet map[*ttnpb.ApplicationIdentifiers]struct{}, localSet map[*ttnpb.ApplicationIdentifiers]struct{}) (complement map[*ttnpb.ApplicationIdentifiers]struct{}) {
	complement = make(map[*ttnpb.ApplicationIdentifiers]struct{})
	for appIds := range localSet {
		if _, ok := isSet[appIds]; ok {
			continue
		}
		complement[appIds] = struct{}{}
	}
	return complement
}

func ComputeDeviceSetComplement(isSet map[*ttnpb.EndDeviceIdentifiers]struct{}, localSet map[*ttnpb.EndDeviceIdentifiers]struct{}) (complement map[*ttnpb.EndDeviceIdentifiers]struct{}) {
	complement = make(map[*ttnpb.EndDeviceIdentifiers]struct{})
	for devIds := range localSet {
		if _, ok := isSet[devIds]; ok {
			continue
		}
		complement[devIds] = struct{}{}
	}
	return complement
}
