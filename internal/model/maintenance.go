package model

// SystemConfig 系统配置
type SystemConfig struct {
	Network     NetworkConfig     `json:"network"`
	Temperature TemperatureConfig `json:"temperature"`
	DeviceInfo  DeviceInfo        `json:"device_info"`
	DateTime    DateTimeConfig    `json:"datetime"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	IP      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
	DNS     string `json:"dns"`
}

// TemperatureConfig 温度配置
type TemperatureConfig struct {
	Low  int `json:"low"`
	High int `json:"high"`
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

// DateTimeConfig 日期时间配置
type DateTimeConfig struct {
	Timezone string `json:"timezone"`
	DateTime string `json:"datetime"`
}

// NetworkConfigRequest 网络配置请求
type NetworkConfigRequest struct {
	IP      string `json:"ip"`
	Mask    string `json:"mask"`
	Gateway string `json:"gateway"`
	DNS     string `json:"dns"`
}

// TemperatureConfigRequest 温度配置请求
type TemperatureConfigRequest struct {
	Low  int `json:"low"`
	High int `json:"high"`
}

// DeviceInfoRequest 设备信息请求
type DeviceInfoRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Contact  string `json:"contact"`
}

// DateTimeRequest 日期时间请求
type DateTimeRequest struct {
	Timezone string `json:"timezone"`
	DateTime string `json:"datetime"`
}

// RebootRequest 重启请求
type RebootRequest struct {
	Delay int `json:"delay"` // 延迟秒数
}

// UserConfig 用户配置
type UserConfig struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Role      int    `json:"role"`
	RoleName  string `json:"role_name"`
	CreatedAt string `json:"created_at"`
}

// UserRequest 用户请求
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

// UserDeleteRequest 用户删除请求
type UserDeleteRequest struct {
	Usernames []string `json:"usernames"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Users []UserConfig `json:"users"`
}

// SessionConfig 会话配置
type SessionConfig struct {
	SessionID    string `json:"session_id"`
	Username     string `json:"username"`
	LoginTime    string `json:"login_time"`
	LastActivity string `json:"last_activity"`
	IPAddress    string `json:"ip_address"`
}

// SessionListResponse 会话列表响应
type SessionListResponse struct {
	Sessions []SessionConfig `json:"sessions"`
}

// SessionDeleteRequest 会话删除请求
type SessionDeleteRequest struct {
	SessionIDs []string `json:"session_ids"`
}

// LogEntry 日志条目
type LogEntry struct {
	ID        int    `json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	Logs  []LogEntry `json:"logs"`
	Total int        `json:"total"`
}

// LogClearRequest 日志清除请求
type LogClearRequest struct {
	Levels []string `json:"levels"` // 可选，清除指定级别的日志
}

// FileInfo 文件信息
type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	UpdatedAt string `json:"updated_at"`
	Type      string `json:"type"`
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	Files []FileInfo `json:"files"`
	Total int        `json:"total"`
}

// FileUploadRequest 文件上传请求 (用于服务层接口)
type FileUploadRequest struct {
	Path      string
	Content   []byte
	Overwrite bool
}

// FileDeleteRequest 文件删除请求
type FileDeleteRequest struct {
	Paths []string `json:"paths"`
}

// SNMPConfig SNMP 配置
type SNMPConfig struct {
	Version     string          `json:"version"`
	Community   string          `json:"community"`
	Contact     string          `json:"contact"`
	Location    string          `json:"location"`
	Enabled     bool            `json:"enabled"`
	TrapEnabled bool            `json:"trap_enabled"`
	TrapHosts   []TrapHost      `json:"trap_hosts"`
	Communities []SNMPCommunity `json:"communities"`
}

// TrapHost Trap 目标主机
type TrapHost struct {
	ID      int    `json:"id"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
}

// SNMPCommunity SNMP 团体
type SNMPCommunity struct {
	Name        string `json:"name"`
	Access      string `json:"access"` // read/write
	Description string `json:"description"`
}

// SNMPConfigRequest SNMP 配置请求
type SNMPConfigRequest struct {
	Version   string `json:"version"`
	Community string `json:"community"`
	Contact   string `json:"contact"`
	Location  string `json:"location"`
	Enabled   bool   `json:"enabled"`
}

// TrapHostRequest Trap 主机请求
type TrapHostRequest struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
}

// TrapHostDeleteRequest Trap 主机删除请求
type TrapHostDeleteRequest struct {
	Host string `json:"host"`
}

// WormRule 蠕虫规则
type WormRule struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
	Stats    int    `json:"stats"`
	Enabled  bool   `json:"enabled"`
}

// WormRuleList 蠕虫规则列表
type WormRuleList struct {
	Rules []WormRule `json:"rules"`
}

// WormRuleRequest 蠕虫规则请求
type WormRuleRequest struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Port     string `json:"port"`
	Enabled  bool   `json:"enabled"`
}

// DDoSConfig DDoS 防护配置
type DDoSConfig struct {
	Enabled   bool   `json:"enabled"`
	Threshold int    `json:"threshold"`
	Action    string `json:"action"`
}

// DDoSConfigRequest DDoS 配置请求
type DDoSConfigRequest struct {
	Enabled   bool   `json:"enabled"`
	Threshold int    `json:"threshold"`
	Action    string `json:"action"`
}

// ARPConfig ARP 防护配置
type ARPConfig struct {
	Enabled        bool     `json:"enabled"`
	InspectEnabled bool     `json:"inspect_enabled"`
	TrustPorts     []string `json:"trust_ports"`
}

// ARPConfigRequest ARP 配置请求
type ARPConfigRequest struct {
	Enabled        bool     `json:"enabled"`
	InspectEnabled bool     `json:"inspect_enabled"`
	TrustPorts     []string `json:"trust_ports"`
}

// LoadConfigRequest 加载配置请求
type LoadConfigRequest struct {
	ConfigFile string `json:"config_file"`
}

// LoadConfigFile 配置文件
type LoadConfigFile struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

// LoadConfigListResponse 配置文件列表响应
type LoadConfigListResponse struct {
	Files []LoadConfigFile `json:"files"`
}
