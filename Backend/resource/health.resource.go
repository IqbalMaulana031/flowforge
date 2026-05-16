package resource

type HealthResource struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

type PingResource struct {
	Message string `json:"message"`
}
