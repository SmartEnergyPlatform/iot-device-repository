Provides metadata about protocols, devices, device-types and value-types. Saves data as rdf, publishes changes to amqp and uses permissionsearch for text-search.


# Depth
Results may represent a entity and its relations in different depths. Depth -1 means that a entity will be returned with all its relations recursively.
A depth of 1 means that only the entity without its relations will be returned. A depth of 2 means that a entity with only its direct relationships will be returned.

# Device-Instance

## GET /deviceInstances/:limit/:offset
Lists device instances where the user has reading permissions. Depth -1.


## GET /deviceInstances/:limit/:offset/:action
Lists device instances where the user has permissions corresponding to action. Valid actions are `read`, `write` and `execute`.  Depth -1.


## GET /deviceInstance/:id
Returns device instance if requesting user has reading permission.  Depth -1.

## POST /deviceInstance
Creates a device instance. Request-Body is a device instance without id and can be prepared with `GET /ui/deviceInstance/resourceSkeleton/:deviceTypeId`.

#### consistent
A device is inconsistent if:
* the url is missing
* the name is missing
* the device-type id does not reference a device-type
* the device-config does not match device-type-config
* the generated device endpoint for a service would collide with a existing endpoint

#### cqrs
If the device is consistent it will be published to amqp and asynchronously consumed to save the device to the database.
Through this published message other services will also be able to save a representation in its own databases. The permissionsearch service is one example for such services.

#### endpoints
If the corresponding device-type has a endpoint-format set for a service, it will be used to generate a endpoint pointing to this device and the service.

#### img
If the ImgUrl field is left empty the img from the associated device type will be used.

#### uri
A for the protocol unique identifier of the device instance.

## POST /deviceInstance/:id
Updates a device instance with similar behavior to the creation.

#### gateway
If the device-url or a device-tag is changed and the device is assigned to a gateway, the gateways hash will be reset.
This guaranties that the client-connector (gateway) sees the change and has the opportunity to publish its device-information.


## DELETE /deviceInstance/:id
removes the device, removes corresponding endpoints, resets hash of assigned gateway. works asynchronous.


## GET /ui/deviceInstance/resourceSkeleton/:deviceTypeId
Returns a skeleton of a new device instance which would be of the given device-type. the skeleton can be used for `POST /deviceInstance`.


## GET /byservice/deviceInstances/:serviceid
Returns device instances that are of the device type that contains the given service and the user has read access to. Depth -1.


## GET /byservice/deviceInstances/:serviceid/:action
Returns device instances that are of the device type that contains the given service and the user has the requested access to. Depth -1.


## GET /bydevicetype/deviceInstances/:devicetype
Returns device instances that are of the given device type and the user has read access to. Depth -1.


## GET /bydevicetype/deviceInstances/:devicetype/:action
Returns device instances that are of the given device type and the user has the requested access to. Depth -1.


## GET /url_to_devices/:device_url
Returns a list of objects containing devices and services where the url matches the requested one and the user has read access. Depth -1.


## POST /endpoint/in
Returns list of endpoints with matching endpoint string and protocol where the user has execute access and the referenced service is declared as sensor.


## POST /endpoint/listen/auth/check/:handler/:endpoint
Checks if user has execute access to all devices referenced by the given endpoint and protocol.


## POST /endpoint/out
Returns endpoint of given device and service if user has execute access.


## POST /endpoint/generate
Creates new device to be referenced by the given endpoint. If no device type exists that has a matching service a new device type will be created.
If no value type exists to match the given message, a new one will be created.


## GET /endpoints/:limit/:offset
Lists all endpoints, if the user has the role `"admin"`


# Gateway

## GET /gateways/:limit/:offset
Lists gateways the user has read access to. Depth -1.


## GET /gateway/:id
Returns the gateway if the user has read access. Depth -1.


## DELETE /gateway/:id
Deletes gateway if user has admin access.


## POST /gateway/:id/clear
Removes all associations with devices and sets hash to empty string if the user has write access.


## POST /gateway/:id/name/:name
Sets the name of the gateway if the user has write access.


## GET /gateway/:id/name
Get the name of the gateway if the user has read access.


## POST /gateway/:id/commit
Assigns devices and a hash to the gateway. The commit will fail if a device is already assigned to a different gateway. If this happens you can clear ore delete the other gateway.
Gateway creation and updates are asynchronous, so you may nead to wait a short moment after creation until you can commit.


## POST /gateway/:id/provide
Returns a gateway with the given id if the user has execute rights.
If no id is passed or the id is unknown a new gateway will be created and returned.



# DeviceType


## GET /service/:id
Returns service if user has access rights to related device type.


## GET /deviceTypes/:limit/:offset
Lists device types the user has read access to.


## GET /maintenance/deviceTypes/:maintenance
Lists device types where the given maintenance flag is set and the user has write access.


## GET /maintenance/check/deviceTypes/:maintenance
Checks if any device type has maintenance flag set.


## GET /deviceType/:id
Returns device type if user has read access. Depth -1.


## GET /deviceType/:id/:depth
Returns devicetype if user has read access, with modifiable depth.


## POST /deviceType
Creates asynchronously a new device type.

#### Consistency
Device type is inconsistent or invalid if:
* the device type has no name
* the device type has no description
* the device type has no assigned vendor
* the device type has no assigned device class
* a config field has no name
* a input/output type assignment has no name
* a input/output type assignment has no assigned known format
* a input/output type assignment has no assigned message segment matching with the protocol
* a input/output type assignment has no assigned valid value type
* a service of the type has no url
* a service of the type has no name
* a service of the type has no description
* a service of the type has no assigned protocol

#### endpoints 
A endpoint references a device-specific instantiation of iot-device-repository.Service.EndpointFormat. The endpoint identifies device-service combinations. 
The iot-device-repository.Service.EndpointFormat is a mustache template which may reference device and service uri.
For example `/some/path/{{device_uri}}/foo/{{service_uri}}/bar`.
The use of the device uri is highly advised, to ensure that for a devicetype each device has its own service endpoints.
The last example used a structure reminiscent of a http-path, but it can be everything. It is advisable to use one that prevents ambiguous endpoints but it may in fact be somthing like `iuhansdcansdc{{device_uri}}asidauisdha`.

If iot-device-repository.Service.EndpointFormat is  not empty the iot-repository will create a endpoint for new device instances created of this type.

#### img
The image link will be used for display in the web-ui. device instances of this type will use this image as default, unless changed.

#### service.url
A for the device type unique identifier

## POST /deviceType/:id
Updates the device type. Endpoints will be updated. device instance images will be updated, if they still use the default image.


## POST /import/deviceType
Creates a new device type without consistency check and with user predefined ids.
WARNING: can harm your database.


## DELETE /deviceType/:id
Deletes the device type if user has administration access and no device instance with this type exists.


## GET /ui/deviceType/allowedvalues
Returns informations about valid values to create new device types and value types.


## POST /query/deviceType
Use request-body as partial device type to find matching device type, similar to a mongodb filter.
Returns document containing existence boolean and device type id.


## POST /query/service
Use request-body as partial service to find device types with matching services, similar to a mongodb filter.
Returns a list of device type ids.


# Intern
this endpoints should be only accessible from internal services, they will not check for the rigths of the requester.

## POST /intern/gatewaynames
Maps names to the given gateway ids


## POST /intern/device/gatewaynames
Maps names of associated gateways to the given device ids


# Other

## POST /other/vendor
Creates a new vendor


## POST /other/protocol
Creates a new protocol. The protocol url will be used as topic for kafka messages.


## POST /other/deviceclass
Creates a new deviceclass


## POST /other/valueType
Creates a new valueType. 

#### consistency
a value type is inconsistent if:
* value type has no name
* value type has no description
* value type has no known base type
* value type has not exactly one field if base type is list or map
* value type has not primitive base type and has more than 0 fields


## GET /skeleton/:instance_id/:service_id
Returns input/output example for device and service as received by the bpmn-process.


## GET /devicetype/skeleton/:type_id/:service_id
like /skeleton/:instance_id/:service_id but with device type.


## POST /format/example
Returns message example in the selected format as potentially send to a device.


# Search

## GET /ui/search/deviceTypes/:query/:limit/:offset
Searches for partial device type matches of name, description etc. (permissionsearch config defines which fields ar searchable). Depth 1.
Searches will be performed by the permissionsearch service, the data will be completed by the rdf triple store.


## GET /ui/search/deviceInstances/:query/:limit/:offset
Searches for partial device instance matches of name, description etc. (permissionsearch config defines which fields ar searchable), where the user has read access. Depth -1.
Searches will be performed by the permissionsearch service, the data will be completed by the rdf triple store.


## GET /ui/search/deviceInstances/:query/:limit/:offset/:action
Searches for partial device instance matches of name, description etc. (permissionsearch config defines which fields ar searchable), where the user has requested (:action) access. Depth -1.
Searches will be performed by the permissionsearch service, the data will be completed by the rdf triple store.


## GET /ui/search/gateways/:query/:limit/:offset
Searches for partial gateway matches of name, description etc. (permissionsearch config defines which fields ar searchable), where the user has read access. Depth -1.
Searches will be performed by the permissionsearch service, the data will be completed by the rdf triple store.


## GET /ui/search/valueTypes/:query/:limit/:offset
Searches for value types


## GET /ui/search/others/:type/:query/:limit/:offset
Searches for miscellaneous entities. Valid types are `vendors`, `deviceClasses` and `protocols` 


# ValueType

## GET /valueTypes/:limit/:offset
Lists value types


## GET /valueType/:id
Returns value type


## POST /valueType/generate
Generates a value type representing message given by the post body. value type will be returned with the used format, but will not be saved.