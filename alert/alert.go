package alert

type AlertLevel = string

var (
	WARNING  = "warning"
	CRITICAL = "critical"
)

type AppAlert struct {
	Level             AlertLevel
	App               string
	AppState          string
	Host              string
	Name              string
	HealthyInstances  int32
	RequiredInstances int32
	FailureReason     string
	MarathonAppURL    string
	MesosSandboxUrl   string
}

type Alertor interface {
	Alert(alert AppAlert) error
}
