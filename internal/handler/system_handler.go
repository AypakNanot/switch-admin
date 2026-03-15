package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/gin-gonic/gin"
	"switch-admin/internal/dao"
)

// SystemHandler 系统配置处理器
type SystemHandler struct {
	configDAO *dao.ConfigDAO
}

// NewSystemHandler 创建系统配置处理器
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		configDAO: dao.NewConfigDAO(),
	}
}

// GetRunMode 获取运行模式（用于 datamodel）
func (h *SystemHandler) GetRunMode() (string, error) {
	return h.configDAO.GetRunMode()
}

// GetSystemConfigPage 获取系统配置页面
func (h *SystemHandler) GetSystemConfigPage(ctx *context.Context) (types.Panel, error) {
	// 获取当前运行模式
	currentMode, _ := h.configDAO.GetRunMode()
	modeDesc := "离线测试模式"
	if currentMode == "switch" {
		modeDesc = "交换机模式"
	}

	// 获取适配器配置
	adapterConfigs, _ := h.getAdapterConfigsHTML()

	// 系统配置信息卡片
	sysInfoHTML := `
	<div class="row">
		<div class="col-md-12">
			<div class="box box-primary">
				<div class="box-header with-border">
					<h3 class="box-title">系统配置</h3>
				</div>
				<div class="box-body">
					<table class="table table-bordered">
						<tr>
							<th width="150">运行模式</th>
							<td>
								<span class="label label-info">` + modeDesc + `</span>
								<span class="text-muted">(` + currentMode + `)</span>
							</td>
						</tr>
						<tr>
							<th>数据库</th>
							<td>SQLite3 (data/admin.db)</td>
						</tr>
						<tr>
							<th>GoAdmin 版本</th>
							<td>v1.2.26</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	</div>
	`

	// 适配器配置表格
	adapterHTML := `
	<div class="row">
		<div class="col-md-12">
			<div class="box box-success">
				<div class="box-header with-border">
					<h3 class="box-title">适配器配置</h3>
				</div>
				<div class="box-body">
					<table class="table table-bordered table-hover">
						<thead>
							<tr>
								<th>功能</th>
								<th>适配器类型</th>
								<th>优先级</th>
								<th>状态</th>
								<th>配置</th>
							</tr>
						</thead>
						<tbody>
							` + adapterConfigs + `
						</tbody>
					</table>
				</div>
			</div>
		</div>
	</div>
	`

	// 模式切换说明
	helpHTML := `
	<div class="row">
		<div class="col-md-12">
			<div class="box box-warning">
				<div class="box-header with-border">
					<h3 class="box-title">💡 使用说明</h3>
				</div>
				<div class="box-body">
					<div class="alert alert-info">
						<h4>运行模式说明：</h4>
						<ul>
							<li><strong>离线测试模式 (mock)</strong>：使用数据库模拟数据，不依赖真实交换机硬件</li>
							<li><strong>交换机模式 (switch)</strong>：连接真实交换机，通过 CLI/Netconf 等方式操作硬件</li>
						</ul>
						<h4>切换模式：</h4>
						<p>可以通过 API 切换运行模式，无需重启服务：</p>
						<pre style="background:#f5f5f5;padding:10px;border-radius:4px;"><code>POST /admin/api/mode
Content-Type: application/json

{"mode": "mock"}  // 或 {"mode": "switch"}</code></pre>
					</div>
				</div>
			</div>
		</div>
	</div>
	`

	allHTML := template.HTML(sysInfoHTML + adapterHTML + helpHTML)

	return types.Panel{
		Content:     allHTML,
		Title:       "系统配置",
		Description: "查看和管理系统配置、适配器配置",
	}, nil
}

// getAdapterConfigsHTML 获取适配器配置 HTML
func (h *SystemHandler) getAdapterConfigsHTML() (string, error) {
	// 获取所有适配器配置
	configs, err := h.configDAO.GetAllAdapterConfigs()
	if err != nil {
		return "", err
	}

	html := ""
	for _, config := range configs {
		functionName := config.FunctionName
		adapterType := config.AdapterType
		priority := config.Priority
		enabled := config.Enabled

		// 状态标签
		statusLabel := `<span class="label label-default">禁用</span>`
		if enabled {
			statusLabel = `<span class="label label-success">启用</span>`
		}

		// 适配器类型标签
		typeLabel := "label-default"
		switch adapterType {
		case "cli":
			typeLabel = "label-primary"
		case "netconf":
			typeLabel = "label-info"
		case "rest":
			typeLabel = "label-warning"
		}

		// 解析配置 JSON 并格式化
		configDisplay := "-"
		if config.Config != "" {
			var configMap map[string]interface{}
			if err := json.Unmarshal([]byte(config.Config), &configMap); err == nil {
				// 格式化显示
				configDisplay = "<pre style='margin:0;font-size:12px'>" + formatConfig(configMap) + "</pre>"
			} else {
				configDisplay = config.Config
			}
		}

		html += `
			<tr>
				<td><strong>` + functionName + `</strong></td>
				<td><span class="label ` + typeLabel + `">` + adapterType + `</span></td>
				<td>` + string(rune(priority+48)) + `</td>
				<td>` + statusLabel + `</td>
				<td style="font-size:12px">` + configDisplay + `</td>
			</tr>
		`
	}

	return html, nil
}

// formatConfig 格式化配置显示
func formatConfig(config map[string]interface{}) string {
	result := ""
	for k, v := range config {
		result += k + ": " + formatValue(v) + "\n"
	}
	return result
}

// formatValue 格式化值
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return string(rune(int(val) + 48))
	default:
		return ""
	}
}

// APISwitchMode API: 切换运行模式
func (h *SystemHandler) APISwitchMode(ctx *context.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	var req struct {
		Mode string `json:"mode"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid JSON format",
		})
		return
	}

	// 验证模式
	if req.Mode != "mock" && req.Mode != "switch" {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid mode. Must be 'mock' or 'switch'",
		})
		return
	}

	// 获取当前模式
	oldMode, _ := h.configDAO.GetRunMode()

	// 切换模式
	if err := h.configDAO.SetRunMode(req.Mode); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to switch mode: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": map[string]string{
			"previous_mode": oldMode,
			"current_mode":  req.Mode,
		},
		"message": "模式切换成功",
	})
}

// APIGetMode API: 获取当前模式
func (h *SystemHandler) APIGetMode(ctx *context.Context) {
	currentMode, err := h.configDAO.GetRunMode()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	modeDesc := "离线测试模式"
	if currentMode == "switch" {
		modeDesc = "交换机模式"
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": map[string]string{
			"mode":        currentMode,
			"description": modeDesc,
		},
	})
}

// GinAPISwitchMode Gin 版本：切换运行模式
func (h *SystemHandler) GinAPISwitchMode(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	var req struct {
		Mode string `json:"mode"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid JSON format",
		})
		return
	}

	// 验证模式
	if req.Mode != "mock" && req.Mode != "switch" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid mode. Must be 'mock' or 'switch'",
		})
		return
	}

	// 获取当前模式
	oldMode, _ := h.configDAO.GetRunMode()

	// 切换模式
	if err := h.configDAO.SetRunMode(req.Mode); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to switch mode: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": map[string]string{
			"previous_mode": oldMode,
			"current_mode":  req.Mode,
		},
		"message": "模式切换成功",
	})
}

// GinAPIGetMode Gin 版本：获取当前模式
func (h *SystemHandler) GinAPIGetMode(c *gin.Context) {
	currentMode, err := h.configDAO.GetRunMode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	modeDesc := "离线测试模式"
	if currentMode == "switch" {
		modeDesc = "交换机模式"
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"data": map[string]string{
			"mode":        currentMode,
			"description": modeDesc,
		},
	})
}
