package configmodel

type ConfigTracingModel struct {
	TraceSessionList []ConfigTraceSessionModel `json:"traceSessionList"`
}

type ConfigTraceSessionModel struct {
	Disable        bool                         `json:"disable"`
	TraceSessionId string                       `json:"traceSessionId"`
	EndTime        string                       `json:"endTime"`
	SubscriberList []ConfigTraceSubscriberModel `json:"subscriberList"`
	InterfaceList  []ConfigTraceInterfaceModel  `json:"interfaceList"`
}

type ConfigTraceSubscriberModel struct {
	IMSI  string `json:"imsi"`
	MSIDN string `json:"msisdn"`
	IMEI  string `json:"imei"`
}

type ConfigTraceInterfaceModel struct {
	PodName     string   `json:"podName"` // If podName is absent, means this applies to all pods of a micro-service
	ServiceName string   `json:"serviceName"`
	Interfaces  []string `json:"interfaces"`
}

func GetTracingConfig() *ConfigTracingModel {
	return confInst.tracingModel
}
