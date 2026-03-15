package datamodel

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/components/infobox"
	"github.com/GoAdminGroup/themes/adminlte/components/smallbox"
)

// GetDashboardContent 返回系统大盘（Dashboard）页面内容
// 根据 PRD 监控模块需求：展示设备基础信息、端口状态概览
func GetDashboardContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()

	/**************************
	 * 系统信息 Small Boxes
	 **************************/
	systemInfo := getSystemInfo()

	smallbox1 := smallbox.New().
		SetColor("blue").
		SetIcon(icon.Server).
		SetUrl("/admin/system/info").
		SetTitle("产品型号").
		SetValue(template.HTML(systemInfo.Model)).
		GetContent()

	smallbox2 := smallbox.New().
		SetColor("purple").
		SetIcon(icon.Cog).
		SetUrl("/admin/system/info").
		SetTitle("软件版本").
		SetValue(template.HTML(systemInfo.SoftwareVersion)).
		GetContent()

	smallbox3 := smallbox.New().
		SetColor("green").
		SetIcon(icon.ClockO).
		SetUrl("/admin/system/info").
		SetTitle("运行时间").
		SetValue(template.HTML(systemInfo.Uptime)).
		GetContent()

	smallbox4 := smallbox.New().
		SetColor("yellow").
		SetIcon(icon.Barcode).
		SetUrl("/admin/system/info").
		SetTitle("序列号").
		SetValue(template.HTML(systemInfo.SerialNumber)).
		GetContent()

	size := types.Size(6, 3, 0).XS(12)
	col1 := colComp.SetSize(size).SetContent(smallbox1).GetContent()
	col2 := colComp.SetSize(size).SetContent(smallbox2).GetContent()
	col3 := colComp.SetSize(size).SetContent(smallbox3).GetContent()
	col4 := colComp.SetSize(size).SetContent(smallbox4).GetContent()
	row1 := components.Row().SetContent(col1 + col2 + col3 + col4).GetContent()

	/**************************
	 * 端口状态概览
	 **************************/
	ports := getPortOverview()

	// 端口统计 InfoBox
	portUpCount := countPortsByStatus(ports, "up")
	portDownCount := countPortsByStatus(ports, "down")
	portTotal := len(ports)
	portUtilization := float64(portUpCount) / float64(portTotal) * 100

	infobox1 := infobox.New().
		SetText("端口总数").
		SetColor("blue").
		SetNumber(template.HTML(formatNumber(portTotal))).
		SetIcon(icon.Th).
		GetContent()

	infobox2 := infobox.New().
		SetText("活跃端口").
		SetColor("green").
		SetNumber(template.HTML(formatNumber(portUpCount))).
		SetIcon(icon.CheckCircle).
		GetContent()

	infobox3 := infobox.New().
		SetText("Down 端口").
		SetColor("red").
		SetNumber(template.HTML(formatNumber(portDownCount))).
		SetIcon(icon.Close).
		GetContent()

	infobox4 := infobox.New().
		SetText("端口利用率").
		SetColor("yellow").
		SetNumber(template.HTML(formatFloat(portUtilization))).
		SetIcon(icon.PieChart).
		GetContent()

	var infoboxSize = types.Size(6, 3, 0).XS(12)
	infoboxCol1 := colComp.SetSize(infoboxSize).SetContent(infobox1).GetContent()
	infoboxCol2 := colComp.SetSize(infoboxSize).SetContent(infobox2).GetContent()
	infoboxCol3 := colComp.SetSize(infoboxSize).SetContent(infobox3).GetContent()
	infoboxCol4 := colComp.SetSize(infoboxSize).SetContent(infobox4).GetContent()
	row2 := components.Row().SetContent(infoboxCol1 + infoboxCol2 + infoboxCol3 + infoboxCol4).GetContent()

	/**************************
	 * 端口状态表格
	 **************************/
	portTable := components.Table().SetInfoList(ports).SetThead(types.Thead{
		{Head: "端口"},
		{Head: "状态"},
		{Head: "链路"},
		{Head: "速率"},
		{Head: "双工"},
		{Head: "描述"},
	}).GetContent()

	portBox := components.Box().
		WithHeadBorder().
		SetHeader("端口状态").
		SetHeadColor("#f7f7f7").
		SetBody(portTable).
		SetFooter(`<div class="clearfix"><a href="/admin/monitor/ports" class="btn btn-sm btn-info btn-flat pull-left">查看所有端口</a></div>`).
		GetContent()

	portCol := colComp.SetSize(types.SizeMD(12)).SetContent(row2 + portBox).GetContent()
	row3 := components.Row().SetContent(portCol).GetContent()

	/**************************
	 * 系统详细信息
	 **************************/
	systemDetail := `<table class="table table-bordered">
		<tr><th width="150">产品型号</th><td>` + systemInfo.Model + `</td></tr>
		<tr><th>序列号</th><td>` + systemInfo.SerialNumber + `</td></tr>
		<tr><th>MAC 地址</th><td>` + systemInfo.MACAddress + `</td></tr>
		<tr><th>软件版本</th><td>` + systemInfo.SoftwareVersion + `</td></tr>
		<tr><th>硬件版本</th><td>` + systemInfo.HardwareVersion + `</td></tr>
		<tr><th>运行时间</th><td>` + systemInfo.Uptime + `</td></tr>
	</table>`

	systemBox := components.Box().
		WithHeadBorder().
		SetHeader("系统详细信息").
		SetHeadColor("#f7f7f7").
		SetBody(template.HTML(systemDetail)).
		SetFooter(`<div class="clearfix"><a href="/admin/system/info" class="btn btn-sm btn-info btn-flat pull-left">查看详细</a></div>`).
		GetContent()

	systemCol := colComp.SetSize(types.SizeMD(12)).SetContent(systemBox).GetContent()
	row4 := components.Row().SetContent(systemCol).GetContent()

	return types.Panel{
		Content:     row1 + row3 + row4,
		Title:       "系统大盘",
		Description: "BroadEdge-S3652 智能交换机管理系统",
	}, nil
}

// SystemInfo 系统信息结构体
type SystemInfo struct {
	Model           string // 产品型号
	SerialNumber    string // 序列号
	MACAddress      string // MAC 地址
	SoftwareVersion string // 软件版本
	HardwareVersion string // 硬件版本
	Uptime          string // 运行时间
}

// PortInfo 端口信息结构体
type PortInfo struct {
	Name        string // 端口名称
	AdminStatus string // 管理状态：enable/disable
	LinkStatus  string // 链路状态：up/down
	Speed       string // 速率：10M/100M/1000M/10G
	Duplex      string // 双工模式：Full/Half
	Description string // 端口描述
}

// getSystemInfo 获取系统信息
func getSystemInfo() SystemInfo {
	return SystemInfo{
		Model:           "BroadEdge-S3652",
		SerialNumber:    "E605MT252088",
		MACAddress:      "00:07:30:D2:35:67",
		SoftwareVersion: "OPTEL v7.0.5.15",
		HardwareVersion: "3.0",
		Uptime:          "0 天 18 小时 42 分钟",
	}
}

// getPortOverview 获取端口概览
func getPortOverview() []map[string]types.InfoItem {
	ports := []PortInfo{
		{Name: "GE1/0/1", AdminStatus: "enable", LinkStatus: "up", Speed: "1000M", Duplex: "Full", Description: "Server-A"},
		{Name: "GE1/0/2", AdminStatus: "enable", LinkStatus: "down", Speed: "-", Duplex: "-", Description: ""},
		{Name: "GE1/0/3", AdminStatus: "enable", LinkStatus: "up", Speed: "1000M", Duplex: "Full", Description: "AP-Floor1"},
		{Name: "GE1/0/4", AdminStatus: "disable", LinkStatus: "down", Speed: "-", Duplex: "-", Description: "Unused"},
	}

	var infoList []map[string]types.InfoItem
	for _, port := range ports {
		linkBadge := port.LinkStatus
		if port.LinkStatus == "up" {
			linkBadge = "<span class='label label-success'>" + port.LinkStatus + "</span>"
		} else {
			linkBadge = "<span class='label label-danger'>" + port.LinkStatus + "</span>"
		}

		infoList = append(infoList, map[string]types.InfoItem{
			"端口":   {Content: template.HTML(port.Name)},
			"状态":   {Content: template.HTML(port.AdminStatus)},
			"链路":   {Content: template.HTML(linkBadge)},
			"速率":   {Content: template.HTML(port.Speed)},
			"双工":   {Content: template.HTML(port.Duplex)},
			"描述":   {Content: template.HTML(port.Description)},
		})
	}
	return infoList
}

// countPortsByStatus 统计指定状态的端口数量
func countPortsByStatus(ports []map[string]types.InfoItem, status string) int {
	count := 0
	for _, port := range ports {
		if port["链路"].Content == template.HTML(status) {
			count++
		}
	}
	return count
}

// formatNumber 格式化数字显示
func formatNumber(n int) string {
	return fmt.Sprintf("%d", n)
}

// formatFloat 格式化浮点数显示
func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f%%", f)
}
