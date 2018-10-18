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

import "sort"

type OrderDeviceType []DeviceType
func (a OrderDeviceType) Len() int           { return len(a) }
func (a OrderDeviceType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderDeviceType) Less(i, j int) bool { return a[i].Id < a[j].Id }


type OrderDeviceInstance []DeviceInstance
func (a OrderDeviceInstance) Len() int           { return len(a) }
func (a OrderDeviceInstance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderDeviceInstance) Less(i, j int) bool { return a[i].Id < a[j].Id }


type OrderValueType []ValueType
func (a OrderValueType) Len() int           { return len(a) }
func (a OrderValueType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OrderValueType) Less(i, j int) bool { return a[i].Id < a[j].Id }


func SortDeviceTypes(deviceTypes *[]DeviceType){
	sort.Sort(OrderDeviceType(*deviceTypes))
}

func SortDeviceInstance(instances *[]DeviceInstance){
	sort.Sort(OrderDeviceInstance(*instances))
}

func SortValueTypes(instances *[]ValueType){
	sort.Sort(OrderValueType(*instances))
}