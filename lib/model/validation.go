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



func (service Service) IsValid() (valid bool, error string){
	if service.Protocol.Name == "" && service.Protocol.Id == "" {
		return false, "missing service protocol"
	}
	if len(service.Url) == 0 {
		return false, "missing service url"
	}
	if len(service.Name) == 0 {
		return false, "missing service name"
	}
	if len(service.Description) == 0 {
		return false, "missing service description"
	}
	return true, error
}

func (this DeviceType) IsValid() (valid bool, error string) {
	//TODO: check if protocol exists
	for _, service := range this.Services {
		valid, error = service.IsValid()
		if !valid {
			return
		}
	}
	//TODO other fields
	return true, error
}

func (this DeviceInstance) IsValid() (valid bool, error string) {
	//TODO
	return true, error
}
