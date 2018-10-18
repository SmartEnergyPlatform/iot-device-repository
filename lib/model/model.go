/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package model

//TODO: User/Owner and time informations

// mgo uses bson and not json --> field names are lowercases (FlowId -> flowid)
// to change the field names in mongodb: `bson:""`

type AuthAction int

const (
	READ AuthAction = iota
	WRITE
	EXECUTE
	ADMINISTRATE
)

type Auth struct {
	ResourceId  string       `json:"resource_id"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Auth"`
	Owner       string       `json:"owner"                               rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#owner"`
	Permissions []Permission `json:"permissions"                         rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#grandsPermissions"`
}

type Permission struct {
	Id      string `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Permission"`
	Role    string `json:"role,omitempty"                      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#role"`
	User    string `json:"user,omitempty"                      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#user"`
	Read    bool   `json:"read"                                rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#read"`
	Write   bool   `json:"write"                               rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#write"`
	Execute bool   `json:"execute"                             rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#execute"`
}

type MsgSegment struct {
	Id          string   `json:"id"            rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#MsgSegment" rdf_root:"true"`
	Name        string   `json:"name"                               rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Constraints []string `json:"constraints"   rdf_ref:"true"       rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasConstraint"`
}

type Protocol struct {
	Id                 string       `json:"id,omitempty"     rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Protocol" rdf_root:"true"`
	ProtocolHandlerUrl string       `json:"protocol_handler_url"          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#protocol_handler_url"`
	Name               string       `json:"name,omitempty"                rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Desc               string       `json:"description"                   rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#description"`
	MsgStructure       []MsgSegment `json:"msg_structure"                 rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasMsgSegment"`
}

type DeviceServiceEntity struct {
	Device   DeviceInstance `json:"device"`
	Services []ShortService `json:"services"`
}

type ShortService struct {
	Id          string `json:"id,omitempty"                           rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Service"`
	ServiceType string `json:"service_type,omitempty"       rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasServiceType"` // "Actuator" || "Sensor"
	Url         string `json:"url,omitempty"                                  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#url"`
}

type ShortDeviceType struct {
	Id       string         `json:"id,omitempty"                  rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceType" rdf_root:"true"`
	Services []ShortService `json:"services,omitempty"                          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasService"`
}

type UltraShortDeviceType struct {
	Id      string `json:"id"                  rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceType" rdf_root:"true"`
	Service string `json:"service"             rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasService"`
}

type Service struct {
	Id             string           `json:"id,omitempty"                           rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Service"`
	ServiceType    string           `json:"service_type,omitempty"       rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasServiceType"` // "Actuator" || "Sensor"
	Name           string           `json:"name,omitempty"                                 rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Description    string           `json:"description,omitempty"                          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#description"`
	Protocol       Protocol         `json:"protocol,omitempty"                             rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasProtocol"`
	Input          []TypeAssignment `json:"input,omitempty"                                rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasInput"`
	Output         []TypeAssignment `json:"output,omitempty"                               rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasOutput"` // list of alternative result types; for example a string if success or a json on error
	Url            string           `json:"url,omitempty"                                  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#url"`
	EndpointFormat string           `json:"endpoint_format,omitempty"                      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#endpoint_format"`
}

type Endpoint struct {
	Id              string `json:"id"                                  rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Endpoint"`
	Endpoint        string `json:"endpoint"                            rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#endpoint"`
	Service         string `json:"service"           rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#refService"`
	Device          string `json:"device"            rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#refDevice"`
	ProtocolHandler string `json:"protocol_handler"                    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#protocol_handler"`
}

type AdditionalFormatInfo struct {
	Id         string    `json:"id,omitempty"                    rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#FormatInfo"`
	Field      FieldType `json:"field"                           rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasFormatinfoForField" rdf_lending:"true"`
	FormatFlag string    `json:"format_flag"                     rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasFormatFlag"`
}

type TypeAssignment struct {
	Id                   string                 `json:"id,omitempty"             rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#MsgSegmentAssignment"`
	Name                 string                 `json:"name"                     rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	MsgSegment           MsgSegment             `json:"msg_segment"              rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasMsgSegment"`
	Type                 ValueType              `json:"type"                     rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasTypeAssigned"`
	Format               string                 `json:"format" rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasFormat"`
	AdditionalFormatinfo []AdditionalFormatInfo `json:"additional_formatinfo"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasFormatInfo"`
}

type DeviceType struct {
	Id          string            `json:"id,omitempty"                  rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceType" rdf_root:"true"`
	Name        string            `json:"name,omitempty"                              rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Description string            `json:"description,omitempty"                       rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#description"`
	Generated   bool              `json:"generated,omitempty"                       rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#generated"`
	Maintenance []string          `json:"maintenance,omitempty"                       rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#maintenance"`
	DeviceClass DeviceClass       `json:"device_class,omitempty"                      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasDeviceClass"`
	Services    []Service         `json:"services,omitempty"                          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasService"`
	Vendor      Vendor            `json:"vendor,omitempty"                            rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#vendor"`
	Config      []ConfigFieldType `json:"config_parameter,omitempty"                  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasConfigParameter"`
	ImgUrl      string            `json:"img,omitempty"             rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#image"`
}

type DeviceInstance struct {
	Id         string        `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceInstance" rdf_root:"true"`
	Name       string        `json:"name,omitempty"                                rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	DeviceType string        `json:"device_type,omitempty"     rdf_ref:"true"      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasDeviceType"`
	Config     []ConfigField `json:"config,omitempty"                              rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasConfig"`
	Url        string        `json:"uri,omitempty"                                 rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#url"`
	Tags       []string      `json:"tags"                                          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasTag"`
	UserTags   []string      `json:"user_tags"                                     rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasUserTag"`
	Gateway    string        `json:"gateway,omitempty"         rdf_ref:"true"      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectedByGateway"`
	ImgUrl     string        `json:"img,omitempty"             rdf_ref:"true"      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#image"`
}

type GatewayRef struct {
	Id      string   `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Gateway" rdf_root:"true"`
	Devices []string `json:"devices,omitempty"     rdf_ref:"true"          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectsDevices"`
	Hash    string   `json:"hash,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hash"`
}

type GatewayFlat struct {
	Id      string   `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Gateway" rdf_root:"true"`
	Devices []string `json:"devices,omitempty"     rdf_ref:"true"          rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectsDevices"`
	Hash    string   `json:"hash,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hash"`
	Name    string   `json:"name,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type Gateway struct {
	Id      string           `json:"id,omitempty"               rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Gateway" rdf_root:"true"`
	Name    string           `json:"name,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Hash    string           `json:"hash,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hash"`
	Devices []DeviceInstance `json:"devices,omitempty"                     rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectsDevices"`
}

type SmartObject struct {
	Id   string `json:"id,omitempty" rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#SmartObject" rdf_root:"true"`
	Name string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type Format struct {
	Id   string `json:"id,omitempty" rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Format" rdf_root:"true"`
	Name string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type Vendor struct {
	Id   string `json:"id,omitempty" rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Vendor" rdf_root:"true"`
	Name string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type DeviceClass struct {
	Id   string `json:"id,omitempty" rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceClass" rdf_root:"true"`
	Name string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type ConfigField struct {
	Id    string `json:"id,omitempty"    rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#ConfigField"`
	Name  string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Value string `json:"value"           rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasValue"`
}

type ConfigFieldType struct {
	Id   string `json:"id,omitempty"    rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#ConfigField"`
	Name string `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

type FieldType struct {
	Id   string    `json:"id,omitempty"    rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#FieldType"`
	Name string    `json:"name,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Type ValueType `json:"type,omitempty"  rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasValueType"`
}

type ValueType struct {
	Id          string      `json:"id,omitempty"             rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#ValueType" rdf_root:"true"`
	Name        string      `json:"name,omitempty"                             rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
	Description string      `json:"description,omitempty"                      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#description"`
	BaseType    string      `json:"base_type,omitempty"      rdf_ref:"true"    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasBaseType"`
	Fields      []FieldType `json:"fields"                           rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasField"`
	Literal     string      `json:"literal"                                    rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#hasLiteral"` //is literal, if not empty
}

type DeviceGatewayRelation struct {
	Id      string `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceInstance" rdf_root:"true"`
	Gateway string `json:"gateway,omitempty"         rdf_ref:"true"      rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectedByGateway"`
}

type DeviceToGateway struct {
	Id      string      `json:"id,omitempty"                        rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#DeviceInstance" rdf_root:"true"`
	Gateway GatewayName `json:"gateway,omitempty"              rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#connectedByGateway"`
}

type GatewayName struct {
	Id   string `json:"id,omitempty"                rdf_entity:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#Gateway" rdf_root:"true"`
	Name string `json:"name,omitempty"                        rdf_field:"http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#name"`
}

const (
	IndexStructBaseType = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#index_structure"
	StructBaseType      = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#structure"
	MapBaseType         = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#map"
	ListBaseType        = "http://www.sepl.wifa.uni-leipzig.de/ontlogies/device-repo#list"

	XsdString = "http://www.w3.org/2001/XMLSchema#string"
	XsdInt    = "http://www.w3.org/2001/XMLSchema#integer"
	XsdFloat  = "http://www.w3.org/2001/XMLSchema#decimal"
	XsdBool   = "http://www.w3.org/2001/XMLSchema#boolean"
)

type AllowedValues struct {
	ServiceTypes []SmartObject `json:"service_types"`
	Formats      []Format      `json:"formats"`
	Primitive    []string      `json:"primitive"`
	Collections  []string      `json:"collections"`
	Structures   []string      `json:"structures"`
	Map          []string      `json:"map"`
	Set          []string      `json:"set"`
}

func GetAllowedValuesBase() AllowedValues {
	return AllowedValues{
		Map: []string{
			MapBaseType,
		},
		Set: []string{
			ListBaseType,
		},
		Collections: []string{
			ListBaseType,
			MapBaseType,
		},
		Structures: []string{
			IndexStructBaseType,
			StructBaseType,
		},
		Primitive: []string{
			XsdString,
			XsdInt,
			XsdFloat,
			XsdBool,
		},
	}
}

func (allowedValues AllowedValues) IsMap(valueType ValueType) bool {
	for _, element := range allowedValues.Map {
		if element == valueType.BaseType {
			return true
		}
	}
	return false
}

func (allowedValues AllowedValues) IsSet(valueType ValueType) bool {
	for _, element := range allowedValues.Set {
		if element == valueType.BaseType {
			return true
		}
	}
	return false
}

func (allowedValues AllowedValues) IsCollection(valueType ValueType) bool {
	for _, element := range allowedValues.Collections {
		if element == valueType.BaseType {
			return true
		}
	}
	return false
}

func (allowedValues AllowedValues) IsStructure(valueType ValueType) bool {
	for _, element := range allowedValues.Structures {
		if element == valueType.BaseType {
			return true
		}
	}
	return false
}

func (allowedValues AllowedValues) IsPrimitive(valueType ValueType) bool {
	for _, element := range allowedValues.Primitive {
		if element == valueType.BaseType {
			return true
		}
	}
	return false
}
