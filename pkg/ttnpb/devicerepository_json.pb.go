// Code generated by protoc-gen-go-json. DO NOT EDIT.
// versions:
// - protoc-gen-go-json v1.1.0
// - protoc             v3.9.1
// source: lorawan-stack/api/devicerepository.proto

package ttnpb

import (
	gogo "github.com/TheThingsIndustries/protoc-gen-go-json/gogo"
	jsonplugin "github.com/TheThingsIndustries/protoc-gen-go-json/jsonplugin"
	types "github.com/gogo/protobuf/types"
)

// KeyProvisioning_customname contains custom string values that override KeyProvisioning_name.
var KeyProvisioning_customname = map[int32]string{
	0: "unknown",
	1: "custom",
	2: "join server",
	3: "manifest",
}

// MarshalProtoJSON marshals the KeyProvisioning to JSON.
func (x KeyProvisioning) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	s.WriteEnumString(int32(x), KeyProvisioning_customname, KeyProvisioning_name)
}

// KeyProvisioning_customvalue contains custom string values that extend KeyProvisioning_value.
var KeyProvisioning_customvalue = map[string]int32{
	"UNKNOWN":     0,
	"unknown":     0,
	"CUSTOM":      1,
	"custom":      1,
	"JOIN_SERVER": 2,
	"join server": 2,
	"MANIFEST":    3,
	"manifest":    3,
}

// UnmarshalProtoJSON unmarshals the KeyProvisioning from JSON.
func (x *KeyProvisioning) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	v := s.ReadEnum(KeyProvisioning_value, KeyProvisioning_customvalue)
	if err := s.Err(); err != nil {
		s.SetErrorf("could not read KeyProvisioning enum: %v", err)
		return
	}
	*x = KeyProvisioning(v)
}

// KeySecurity_customname contains custom string values that override KeySecurity_name.
var KeySecurity_customname = map[int32]string{
	0: "unknown",
	1: "none",
	2: "read protected",
	3: "secure element",
}

// MarshalProtoJSON marshals the KeySecurity to JSON.
func (x KeySecurity) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	s.WriteEnumString(int32(x), KeySecurity_customname, KeySecurity_name)
}

// KeySecurity_customvalue contains custom string values that extend KeySecurity_value.
var KeySecurity_customvalue = map[string]int32{
	"UNKNOWN":        0,
	"unknown":        0,
	"NONE":           1,
	"none":           1,
	"READ_PROTECTED": 2,
	"read protected": 2,
	"SECURE_ELEMENT": 3,
	"secure element": 3,
}

// UnmarshalProtoJSON unmarshals the KeySecurity from JSON.
func (x *KeySecurity) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	v := s.ReadEnum(KeySecurity_value, KeySecurity_customvalue)
	if err := s.Err(); err != nil {
		s.SetErrorf("could not read KeySecurity enum: %v", err)
		return
	}
	*x = KeySecurity(v)
}

// MarshalProtoJSON marshals the EndDeviceModel message to JSON.
func (x *EndDeviceModel) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.BrandId != "" || s.HasField("brand_id") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("brand_id")
		s.WriteString(x.BrandId)
	}
	if x.ModelId != "" || s.HasField("model_id") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("model_id")
		s.WriteString(x.ModelId)
	}
	if x.Name != "" || s.HasField("name") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("name")
		s.WriteString(x.Name)
	}
	if x.Description != "" || s.HasField("description") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("description")
		s.WriteString(x.Description)
	}
	if len(x.HardwareVersions) > 0 || s.HasField("hardware_versions") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("hardware_versions")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.HardwareVersions {
			s.WriteMoreIf(&wroteElement)
			// NOTE: EndDeviceModel_HardwareVersion does not seem to implement MarshalProtoJSON.
			gogo.MarshalMessage(s, element)
		}
		s.WriteArrayEnd()
	}
	if len(x.FirmwareVersions) > 0 || s.HasField("firmware_versions") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("firmware_versions")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.FirmwareVersions {
			s.WriteMoreIf(&wroteElement)
			// NOTE: EndDeviceModel_FirmwareVersion does not seem to implement MarshalProtoJSON.
			gogo.MarshalMessage(s, element)
		}
		s.WriteArrayEnd()
	}
	if len(x.Sensors) > 0 || s.HasField("sensors") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("sensors")
		s.WriteStringArray(x.Sensors)
	}
	if x.Dimensions != nil || s.HasField("dimensions") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("dimensions")
		// NOTE: EndDeviceModel_Dimensions does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.Dimensions)
	}
	if x.Weight != nil || s.HasField("weight") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("weight")
		if x.Weight == nil {
			s.WriteNil()
		} else {
			s.WriteFloat32(x.Weight.Value)
		}
	}
	if x.Battery != nil || s.HasField("battery") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("battery")
		// NOTE: EndDeviceModel_Battery does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.Battery)
	}
	if x.OperatingConditions != nil || s.HasField("operating_conditions") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("operating_conditions")
		// NOTE: EndDeviceModel_OperatingConditions does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.OperatingConditions)
	}
	if x.IpCode != "" || s.HasField("ip_code") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("ip_code")
		s.WriteString(x.IpCode)
	}
	if len(x.KeyProvisioning) > 0 || s.HasField("key_provisioning") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("key_provisioning")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.KeyProvisioning {
			s.WriteMoreIf(&wroteElement)
			element.MarshalProtoJSON(s)
		}
		s.WriteArrayEnd()
	}
	if x.KeySecurity != 0 || s.HasField("key_security") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("key_security")
		x.KeySecurity.MarshalProtoJSON(s)
	}
	if x.Photos != nil || s.HasField("photos") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("photos")
		// NOTE: EndDeviceModel_Photos does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.Photos)
	}
	if x.Videos != nil || s.HasField("videos") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("videos")
		// NOTE: EndDeviceModel_Videos does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.Videos)
	}
	if x.ProductUrl != "" || s.HasField("product_url") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("product_url")
		s.WriteString(x.ProductUrl)
	}
	if x.DatasheetUrl != "" || s.HasField("datasheet_url") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("datasheet_url")
		s.WriteString(x.DatasheetUrl)
	}
	if len(x.Resellers) > 0 || s.HasField("resellers") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("resellers")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.Resellers {
			s.WriteMoreIf(&wroteElement)
			// NOTE: EndDeviceModel_Reseller does not seem to implement MarshalProtoJSON.
			gogo.MarshalMessage(s, element)
		}
		s.WriteArrayEnd()
	}
	if x.Compliances != nil || s.HasField("compliances") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("compliances")
		// NOTE: EndDeviceModel_Compliances does not seem to implement MarshalProtoJSON.
		gogo.MarshalMessage(s, x.Compliances)
	}
	if len(x.AdditionalRadios) > 0 || s.HasField("additional_radios") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("additional_radios")
		s.WriteStringArray(x.AdditionalRadios)
	}
	s.WriteObjectEnd()
}

// UnmarshalProtoJSON unmarshals the EndDeviceModel message from JSON.
func (x *EndDeviceModel) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "brand_id", "brandId":
			s.AddField("brand_id")
			x.BrandId = s.ReadString()
		case "model_id", "modelId":
			s.AddField("model_id")
			x.ModelId = s.ReadString()
		case "name":
			s.AddField("name")
			x.Name = s.ReadString()
		case "description":
			s.AddField("description")
			x.Description = s.ReadString()
		case "hardware_versions", "hardwareVersions":
			s.AddField("hardware_versions")
			s.ReadArray(func() {
				// NOTE: EndDeviceModel_HardwareVersion does not seem to implement UnmarshalProtoJSON.
				var v EndDeviceModel_HardwareVersion
				gogo.UnmarshalMessage(s, &v)
				x.HardwareVersions = append(x.HardwareVersions, &v)
			})
		case "firmware_versions", "firmwareVersions":
			s.AddField("firmware_versions")
			s.ReadArray(func() {
				// NOTE: EndDeviceModel_FirmwareVersion does not seem to implement UnmarshalProtoJSON.
				var v EndDeviceModel_FirmwareVersion
				gogo.UnmarshalMessage(s, &v)
				x.FirmwareVersions = append(x.FirmwareVersions, &v)
			})
		case "sensors":
			s.AddField("sensors")
			x.Sensors = s.ReadStringArray()
		case "dimensions":
			s.AddField("dimensions")
			// NOTE: EndDeviceModel_Dimensions does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_Dimensions
			gogo.UnmarshalMessage(s, &v)
			x.Dimensions = &v
		case "weight":
			s.AddField("weight")
			if !s.ReadNil() {
				v := s.ReadFloat32()
				if s.Err() != nil {
					return
				}
				x.Weight = &types.FloatValue{Value: v}
			}
		case "battery":
			s.AddField("battery")
			// NOTE: EndDeviceModel_Battery does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_Battery
			gogo.UnmarshalMessage(s, &v)
			x.Battery = &v
		case "operating_conditions", "operatingConditions":
			s.AddField("operating_conditions")
			// NOTE: EndDeviceModel_OperatingConditions does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_OperatingConditions
			gogo.UnmarshalMessage(s, &v)
			x.OperatingConditions = &v
		case "ip_code", "ipCode":
			s.AddField("ip_code")
			x.IpCode = s.ReadString()
		case "key_provisioning", "keyProvisioning":
			s.AddField("key_provisioning")
			s.ReadArray(func() {
				var v KeyProvisioning
				v.UnmarshalProtoJSON(s)
				x.KeyProvisioning = append(x.KeyProvisioning, v)
			})
		case "key_security", "keySecurity":
			s.AddField("key_security")
			x.KeySecurity.UnmarshalProtoJSON(s)
		case "photos":
			s.AddField("photos")
			// NOTE: EndDeviceModel_Photos does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_Photos
			gogo.UnmarshalMessage(s, &v)
			x.Photos = &v
		case "videos":
			s.AddField("videos")
			// NOTE: EndDeviceModel_Videos does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_Videos
			gogo.UnmarshalMessage(s, &v)
			x.Videos = &v
		case "product_url", "productUrl":
			s.AddField("product_url")
			x.ProductUrl = s.ReadString()
		case "datasheet_url", "datasheetUrl":
			s.AddField("datasheet_url")
			x.DatasheetUrl = s.ReadString()
		case "resellers":
			s.AddField("resellers")
			s.ReadArray(func() {
				// NOTE: EndDeviceModel_Reseller does not seem to implement UnmarshalProtoJSON.
				var v EndDeviceModel_Reseller
				gogo.UnmarshalMessage(s, &v)
				x.Resellers = append(x.Resellers, &v)
			})
		case "compliances":
			s.AddField("compliances")
			// NOTE: EndDeviceModel_Compliances does not seem to implement UnmarshalProtoJSON.
			var v EndDeviceModel_Compliances
			gogo.UnmarshalMessage(s, &v)
			x.Compliances = &v
		case "additional_radios", "additionalRadios":
			s.AddField("additional_radios")
			x.AdditionalRadios = s.ReadStringArray()
		}
	})
}

// MarshalProtoJSON marshals the ListEndDeviceModelsResponse message to JSON.
func (x *ListEndDeviceModelsResponse) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if len(x.Models) > 0 || s.HasField("models") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("models")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.Models {
			s.WriteMoreIf(&wroteElement)
			element.MarshalProtoJSON(s.WithField("models"))
		}
		s.WriteArrayEnd()
	}
	s.WriteObjectEnd()
}

// UnmarshalProtoJSON unmarshals the ListEndDeviceModelsResponse message from JSON.
func (x *ListEndDeviceModelsResponse) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "models":
			s.AddField("models")
			s.ReadArray(func() {
				if s.ReadNil() {
					x.Models = append(x.Models, nil)
					return
				}
				v := &EndDeviceModel{}
				v.UnmarshalProtoJSON(s.WithField("models", false))
				if s.Err() != nil {
					return
				}
				x.Models = append(x.Models, v)
			})
		}
	})
}

// MarshalProtoJSON marshals the MessagePayloadDecoder message to JSON.
func (x *MessagePayloadDecoder) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.Formatter != 0 || s.HasField("formatter") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("formatter")
		x.Formatter.MarshalProtoJSON(s)
	}
	if x.FormatterParameter != "" || s.HasField("formatter_parameter") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("formatter_parameter")
		s.WriteString(x.FormatterParameter)
	}
	if x.CodecId != "" || s.HasField("codec_id") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("codec_id")
		s.WriteString(x.CodecId)
	}
	if len(x.Examples) > 0 || s.HasField("examples") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("examples")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.Examples {
			s.WriteMoreIf(&wroteElement)
			// NOTE: MessagePayloadDecoder_Example does not seem to implement MarshalProtoJSON.
			gogo.MarshalMessage(s, element)
		}
		s.WriteArrayEnd()
	}
	s.WriteObjectEnd()
}

// UnmarshalProtoJSON unmarshals the MessagePayloadDecoder message from JSON.
func (x *MessagePayloadDecoder) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "formatter":
			s.AddField("formatter")
			x.Formatter.UnmarshalProtoJSON(s)
		case "formatter_parameter", "formatterParameter":
			s.AddField("formatter_parameter")
			x.FormatterParameter = s.ReadString()
		case "codec_id", "codecId":
			s.AddField("codec_id")
			x.CodecId = s.ReadString()
		case "examples":
			s.AddField("examples")
			s.ReadArray(func() {
				// NOTE: MessagePayloadDecoder_Example does not seem to implement UnmarshalProtoJSON.
				var v MessagePayloadDecoder_Example
				gogo.UnmarshalMessage(s, &v)
				x.Examples = append(x.Examples, &v)
			})
		}
	})
}

// MarshalProtoJSON marshals the MessagePayloadEncoder message to JSON.
func (x *MessagePayloadEncoder) MarshalProtoJSON(s *jsonplugin.MarshalState) {
	if x == nil {
		s.WriteNil()
		return
	}
	s.WriteObjectStart()
	var wroteField bool
	if x.Formatter != 0 || s.HasField("formatter") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("formatter")
		x.Formatter.MarshalProtoJSON(s)
	}
	if x.FormatterParameter != "" || s.HasField("formatter_parameter") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("formatter_parameter")
		s.WriteString(x.FormatterParameter)
	}
	if x.CodecId != "" || s.HasField("codec_id") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("codec_id")
		s.WriteString(x.CodecId)
	}
	if len(x.Examples) > 0 || s.HasField("examples") {
		s.WriteMoreIf(&wroteField)
		s.WriteObjectField("examples")
		s.WriteArrayStart()
		var wroteElement bool
		for _, element := range x.Examples {
			s.WriteMoreIf(&wroteElement)
			// NOTE: MessagePayloadEncoder_Example does not seem to implement MarshalProtoJSON.
			gogo.MarshalMessage(s, element)
		}
		s.WriteArrayEnd()
	}
	s.WriteObjectEnd()
}

// UnmarshalProtoJSON unmarshals the MessagePayloadEncoder message from JSON.
func (x *MessagePayloadEncoder) UnmarshalProtoJSON(s *jsonplugin.UnmarshalState) {
	if s.ReadNil() {
		return
	}
	s.ReadObject(func(key string) {
		switch key {
		default:
			s.ReadAny() // ignore unknown field
		case "formatter":
			s.AddField("formatter")
			x.Formatter.UnmarshalProtoJSON(s)
		case "formatter_parameter", "formatterParameter":
			s.AddField("formatter_parameter")
			x.FormatterParameter = s.ReadString()
		case "codec_id", "codecId":
			s.AddField("codec_id")
			x.CodecId = s.ReadString()
		case "examples":
			s.AddField("examples")
			s.ReadArray(func() {
				// NOTE: MessagePayloadEncoder_Example does not seem to implement UnmarshalProtoJSON.
				var v MessagePayloadEncoder_Example
				gogo.UnmarshalMessage(s, &v)
				x.Examples = append(x.Examples, &v)
			})
		}
	})
}
