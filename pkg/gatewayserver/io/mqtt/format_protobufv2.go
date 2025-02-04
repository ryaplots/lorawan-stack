// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
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

package mqtt

import (
	"context"
	"fmt"
	"strconv"
	"time"

	ttnpbv2 "go.thethings.network/lorawan-stack-legacy/v2/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io/mqtt/topics"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/datarate"
)

// eirpDelta is the delta between EIRP and ERP.
const eirpDelta = 2.15

var (
	sourceToV3 = map[ttnpbv2.LocationMetadata_LocationSource]ttnpb.LocationSource{
		ttnpbv2.LocationMetadata_GPS:            ttnpb.SOURCE_GPS,
		ttnpbv2.LocationMetadata_CONFIG:         ttnpb.SOURCE_REGISTRY,
		ttnpbv2.LocationMetadata_REGISTRY:       ttnpb.SOURCE_REGISTRY,
		ttnpbv2.LocationMetadata_IP_GEOLOCATION: ttnpb.SOURCE_IP_GEOLOCATION,
	}

	frequencyPlanToBand = map[ttnpbv2.FrequencyPlan]string{
		ttnpbv2.FrequencyPlan_EU_863_870: "EU_863_870",
		ttnpbv2.FrequencyPlan_US_902_928: "US_902_928",
		ttnpbv2.FrequencyPlan_CN_779_787: "CN_779_787",
		ttnpbv2.FrequencyPlan_EU_433:     "EU_433",
		ttnpbv2.FrequencyPlan_AU_915_928: "AU_915_928",
		ttnpbv2.FrequencyPlan_CN_470_510: "CN_470_510",
		ttnpbv2.FrequencyPlan_AS_923:     "AS_923",
		ttnpbv2.FrequencyPlan_AS_920_923: "AS_923",
		ttnpbv2.FrequencyPlan_AS_923_925: "AS_923",
		ttnpbv2.FrequencyPlan_KR_920_923: "KR_920_923",
		ttnpbv2.FrequencyPlan_IN_865_867: "IN_865_867",
		ttnpbv2.FrequencyPlan_RU_864_870: "RU_864_870",
	}

	errNotScheduled    = errors.DefineInvalidArgument("not_scheduled", "not scheduled")
	errLoRaWANPayload  = errors.DefineInvalidArgument("lorawan_payload", "invalid LoRaWAN payload")
	errLoRaWANMetadata = errors.DefineInvalidArgument("lorawan_metadata", "missing LoRaWAN metadata")
	errDataRate        = errors.DefineInvalidArgument("data_rate", "unknown data rate `{data_rate}`")
	errModulation      = errors.DefineInvalidArgument("modulation", "unknown modulation `{modulation}`")
	errFrequencyPlan   = errors.DefineNotFound("frequency_plan", "unknown frequency plan `{frequency_plan}`")
)

type protobufv2 struct {
	topics.Layout
}

// TimePtr converts time to a pointer.
func TimePtr(x time.Time) *time.Time {
	return &x
}

func (protobufv2) FromDownlink(down *ttnpb.DownlinkMessage, _ ttnpb.GatewayIdentifiers) ([]byte, error) {
	settings := down.GetScheduled()
	if settings == nil {
		return nil, errNotScheduled.New()
	}
	lorawan := &ttnpbv2.LoRaWANTxConfiguration{}
	if pld, ok := down.GetPayload().GetPayload().(*ttnpb.Message_MacPayload); ok && pld != nil {
		lorawan.FCnt = pld.MacPayload.FHDR.FCnt
	}
	switch dr := settings.DataRate.Modulation.(type) {
	case *ttnpb.DataRate_Lora:
		lorawan.Modulation = ttnpbv2.Modulation_LORA
		lorawan.CodingRate = settings.CodingRate
		lorawan.DataRate = fmt.Sprintf("SF%dBW%d", dr.Lora.SpreadingFactor, dr.Lora.Bandwidth/1000)
	case *ttnpb.DataRate_Fsk:
		lorawan.Modulation = ttnpbv2.Modulation_FSK
		lorawan.BitRate = dr.Fsk.BitRate
	default:
		return nil, errModulation.New()
	}

	v2downlink := &ttnpbv2.DownlinkMessage{
		Payload: down.RawPayload,
		GatewayConfiguration: &ttnpbv2.GatewayTxConfiguration{
			Frequency:             settings.Frequency,
			Power:                 int32(settings.Downlink.TxPower - eirpDelta),
			PolarizationInversion: true,
			RfChain:               0,
			Timestamp:             settings.Timestamp,
		},
		ProtocolConfiguration: &ttnpbv2.ProtocolTxConfiguration{
			Lorawan: lorawan,
		},
	}
	return v2downlink.Marshal()
}

func (protobufv2) ToUplink(message []byte, ids ttnpb.GatewayIdentifiers) (*ttnpb.UplinkMessage, error) {
	v2uplink := &ttnpbv2.UplinkMessage{}
	err := v2uplink.Unmarshal(message)
	if err != nil {
		return nil, err
	}

	lorawanMetadata := v2uplink.GetProtocolMetadata().GetLorawan()
	if lorawanMetadata == nil {
		return nil, errLoRaWANMetadata.New()
	}
	gwMetadata := v2uplink.GatewayMetadata
	uplink := &ttnpb.UplinkMessage{
		RawPayload: v2uplink.Payload,
	}

	settings := ttnpb.TxSettings{
		Frequency: gwMetadata.Frequency,
		Timestamp: gwMetadata.Timestamp,
	}
	switch lorawanMetadata.Modulation {
	case ttnpbv2.Modulation_LORA:
		bandID, ok := frequencyPlanToBand[lorawanMetadata.FrequencyPlan]
		if !ok {
			return nil, errFrequencyPlan.WithAttributes("frequency_plan", lorawanMetadata.FrequencyPlan)
		}
		phy, err := band.GetLatest(bandID)
		if err != nil {
			return nil, err
		}
		var drIndex ttnpb.DataRateIndex
		var found bool
		loraDr, err := datarate.ParseLoRa(lorawanMetadata.DataRate)
		if err != nil {
			return nil, err
		}
		for bandDRIndex, bandDR := range phy.DataRates {
			if bandDR.Rate.Equal(loraDr.DataRate) {
				found = true
				drIndex = bandDRIndex
				break
			}
		}
		if !found {
			return nil, errDataRate.WithAttributes("data_rate", lorawanMetadata.DataRate)
		}
		settings.DataRate = loraDr.DataRate
		settings.CodingRate = lorawanMetadata.CodingRate
		settings.DataRateIndex = drIndex
	case ttnpbv2.Modulation_FSK:
		settings.DataRate = ttnpb.DataRate{
			Modulation: &ttnpb.DataRate_Fsk{
				Fsk: &ttnpb.FSKDataRate{
					BitRate: lorawanMetadata.BitRate,
				},
			},
		}
	default:
		return nil, errModulation.WithAttributes("modulation", lorawanMetadata.Modulation)
	}

	mdTime := time.Unix(0, gwMetadata.Time)
	if antennas := gwMetadata.Antennas; len(antennas) > 0 {
		for _, antenna := range antennas {
			rssi := antenna.ChannelRssi
			if rssi == 0 {
				rssi = antenna.Rssi
			}
			uplink.RxMetadata = append(uplink.RxMetadata, &ttnpb.RxMetadata{
				GatewayIdentifiers:    ids,
				AntennaIndex:          antenna.Antenna,
				ChannelRssi:           rssi,
				FrequencyOffset:       antenna.FrequencyOffset,
				Rssi:                  rssi,
				RssiStandardDeviation: antenna.RssiStandardDeviation,
				Snr:                   antenna.Snr,
				Time:                  &mdTime,
				Timestamp:             gwMetadata.Timestamp,
			})
		}
	} else {
		uplink.RxMetadata = append(uplink.RxMetadata, &ttnpb.RxMetadata{
			GatewayIdentifiers: ids,
			AntennaIndex:       0,
			ChannelRssi:        gwMetadata.Rssi,
			Rssi:               gwMetadata.Rssi,
			Snr:                gwMetadata.Snr,
			Time:               &mdTime,
			Timestamp:          gwMetadata.Timestamp,
		})
	}
	uplink.Settings = settings

	return uplink, nil
}

func (protobufv2) ToStatus(message []byte, _ ttnpb.GatewayIdentifiers) (*ttnpb.GatewayStatus, error) {
	v2status := &ttnpbv2.StatusMessage{}
	err := v2status.Unmarshal(message)
	if err != nil {
		return nil, err
	}
	metrics := map[string]float32{
		"lmnw": float32(v2status.LmNw),
		"lmst": float32(v2status.LmSt),
		"lmok": float32(v2status.LmOk),
		"lpps": float32(v2status.LPps),
		"rxin": float32(v2status.RxIn),
		"rxok": float32(v2status.RxOk),
		"txin": float32(v2status.TxIn),
		"txok": float32(v2status.TxOk),
	}
	if os := v2status.Os; os != nil {
		metrics["cpu_percentage"] = os.CpuPercentage
		metrics["load_1"] = os.Load_1
		metrics["load_5"] = os.Load_5
		metrics["load_15"] = os.Load_15
		metrics["memory_percentage"] = os.MemoryPercentage
		metrics["temp"] = os.Temperature
	}
	if v2status.Rtt != 0 {
		metrics["rtt_ms"] = float32(v2status.Rtt)
	}
	versions := make(map[string]string)
	if v2status.Dsp > 0 {
		versions["dsp"] = strconv.Itoa(int(v2status.Dsp))
	}
	if v2status.Fpga > 0 {
		versions["fpga"] = strconv.Itoa(int(v2status.Fpga))
	}
	if v2status.Hal != "" {
		versions["hal"] = v2status.Hal
	}
	if v2status.Platform != "" {
		versions["platform"] = v2status.Platform
	}
	var antennasLocation []*ttnpb.Location
	if loc := v2status.Location; loc.Validate() {
		antennasLocation = []*ttnpb.Location{
			{
				Accuracy:  loc.Accuracy,
				Altitude:  loc.Altitude,
				Latitude:  float64(loc.Latitude),
				Longitude: float64(loc.Longitude),
				Source:    sourceToV3[loc.Source],
			},
		}
	}
	return &ttnpb.GatewayStatus{
		AntennaLocations: antennasLocation,
		BootTime:         TimePtr(time.Unix(0, v2status.BootTime)),
		Ip:               v2status.Ip,
		Metrics:          metrics,
		Time:             TimePtr(time.Unix(0, v2status.Time)),
		Versions:         versions,
	}, nil
}

func (protobufv2) ToTxAck(message []byte, _ ttnpb.GatewayIdentifiers) (*ttnpb.TxAcknowledgment, error) {
	return nil, errNotSupported.New()
}

// NewProtobufV2 returns a format that uses the legacy The Things Stack V2 Protocol Buffers marshaling and unmarshaling.
func NewProtobufV2(ctx context.Context) Format {
	return &protobufv2{
		Layout: topics.NewV2(ctx),
	}
}
