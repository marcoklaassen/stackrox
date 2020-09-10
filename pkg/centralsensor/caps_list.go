package centralsensor

const (
	// PullTelemetryDataCap identifies the capability to pull telemetry data from sensor.
	PullTelemetryDataCap SensorCapability = "PullTelemetryData"

	// CancelTelemetryPullCap identifies the capability to cancel an ongoing telemetry data pull.
	CancelTelemetryPullCap SensorCapability = "CancelTelemetryPull"

	// SensorDetectionCap identifies the capability to run detection from sensor
	SensorDetectionCap SensorCapability = "SensorDetection"

	// ComplianceInNodesCap identifies the capability to run compliance in compliance pods
	ComplianceInNodesCap SensorCapability = "ComplianceInNodes"

	// HealthMonitoringCap identifies the capability to send health information
	HealthMonitoringCap SensorCapability = "HealthMonitoring"

	// NetworkGraphExternalSrcsCap identifies the capability to handle custom network graph external sources.
	NetworkGraphExternalSrcsCap SensorCapability = "NetworkGraphExternalSrcs"
)
