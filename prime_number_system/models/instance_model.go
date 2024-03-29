package models

type InstanceModel struct {
	InstanceId     *string              `json:"instanceId"`
	HostName       *string              `json:"hostName"`
	App            *string              `json:"app"`
	IpAddress      *string              `json:"ipAddr"`
	Status         *string              `json:"status"`
	HealthCheckUrl *string              `json:"healthCheckUrl"`
	StatusPageUrl  *string              `json:"statusPageUrl"`
	HomePageUrl    *string              `json:"homePageUrl"`
	Port           *PortModel           `json:"port"`
	DataCenterInfo *DataCenterInfoModel `json:"dataCenterInfo"`
	MetaData       *MetaDataModel       `json:"metadata"`
}

type PortModel struct {
	PortNumber *int64  `json:"$"`
	Enabled    *string `json:"@enabled"`
}

type MetaDataModel struct {
	Role *string `json:"role"`
}

type DataCenterInfoModel struct {
	Class *string `json:"@class"`
	Name  *string `json:"name"`
}

type InstancesModel struct {
	Application *ApplicationModel `json:"application"`
}

type ApplicationModel struct {
	Name     *string          `json:"name"`
	Instance []*InstanceModel `json:"instance"`
}

type NodesModel struct {
	Nodes []*ApplicationModel `json:"application"`
}
