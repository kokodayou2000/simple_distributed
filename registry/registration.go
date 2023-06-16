package registry

type Registration struct {
	ServiceName      ServiceName
	ServiceURL       string
	RequiredServices []ServiceName // 服务名称
	ServiceUpdateURL string        // 客户端的服务要暴露该url，比如老的服务更换成新的服务了
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)

// patchEntry 条目
type patchEntry struct {
	Name ServiceName
	URL  string
}
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
