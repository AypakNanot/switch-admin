package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSystemConfigContent 系统配置页面
func getSystemConfigContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	systemConfigContent := `
	<style>
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 300px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.section-title { border-bottom: 1px solid #ddd; padding-bottom: 10px; margin: 20px 0 15px 0; font-size: 16px; font-weight: bold; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">网络配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>设备 IP 地址</label>
				<input type="text" class="form-control" id="device-ip" placeholder="例如：192.168.1.1">
			</div>
			<div class="form-group">
				<label>子网掩码</label>
				<input type="text" class="form-control" id="subnet-mask" placeholder="例如：255.255.255.0">
			</div>
			<div class="form-group">
				<label>默认网关</label>
				<input type="text" class="form-control" id="default-gateway" placeholder="例如：192.168.1.254">
			</div>
			<div class="form-group">
				<label>DNS 服务器</label>
				<input type="text" class="form-control" id="dns-server" placeholder="例如：8.8.8.8">
			</div>
			<button type="button" class="btn btn-primary" onclick="saveNetworkConfig()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">温度阈值配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>低温告警阈值 (°C)</label>
				<input type="number" class="form-control" id="temp-low" value="-10" min="-40" max="0">
			</div>
			<div class="form-group">
				<label>高温告警阈值 (°C)</label>
				<input type="number" class="form-control" id="temp-high" value="60" min="40" max="85">
			</div>
			<button type="button" class="btn btn-primary" onclick="saveTemperatureConfig()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">设备信息</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>设备名称</label>
				<input type="text" class="form-control" id="device-name" placeholder="例如：Switch-Core-01">
			</div>
			<div class="form-group">
				<label>设备位置</label>
				<input type="text" class="form-control" id="device-location" placeholder="例如：机房 A-机柜 3-位置 15">
			</div>
			<div class="form-group">
				<label>联系人</label>
				<input type="text" class="form-control" id="contact-person" placeholder="例如：张三 -13800138000">
			</div>
			<button type="button" class="btn btn-primary" onclick="saveDeviceInfo()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">日期时间配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>时区</label>
				<select class="form-control" id="timezone">
					<option value="UTC+8">中国标准时间 (UTC+8)</option>
					<option value="UTC+0">UTC</option>
					<option value="UTC-5">美国东部时间 (UTC-5)</option>
				</select>
			</div>
			<div class="form-group">
				<label>日期时间</label>
				<input type="datetime-local" class="form-control" id="datetime">
			</div>
			<button type="button" class="btn btn-primary" onclick="saveDateTime()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<script>
	function loadConfig() {
		fetch('/api/v1/system/config')
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data) {
					var cfg = data.data;
					if (cfg.network) {
						document.getElementById('device-ip').value = cfg.network.ip || '';
						document.getElementById('subnet-mask').value = cfg.network.mask || '';
						document.getElementById('default-gateway').value = cfg.network.gateway || '';
						document.getElementById('dns-server').value = cfg.network.dns || '';
					}
					if (cfg.temperature) {
						document.getElementById('temp-low').value = cfg.temperature.low || -10;
						document.getElementById('temp-high').value = cfg.temperature.high || 60;
					}
					if (cfg.device_info) {
						document.getElementById('device-name').value = cfg.device_info.name || '';
						document.getElementById('device-location').value = cfg.device_info.location || '';
						document.getElementById('contact-person').value = cfg.device_info.contact || '';
					}
					if (cfg.datetime) {
						document.getElementById('timezone').value = cfg.datetime.timezone || 'UTC+8';
						document.getElementById('datetime').value = cfg.datetime.datetime || '';
					}
				}
			});
	}

	function saveNetworkConfig() {
		var payload = {
			ip: document.getElementById('device-ip').value,
			mask: document.getElementById('subnet-mask').value,
			gateway: document.getElementById('default-gateway').value,
			dns: document.getElementById('dns-server').value
		};
		fetch('/api/v1/system/network', {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(payload)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('网络配置保存成功');
			} else {
				alert(data.message || '保存失败');
			}
		});
	}

	function saveTemperatureConfig() {
		var payload = {
			low: parseInt(document.getElementById('temp-low').value),
			high: parseInt(document.getElementById('temp-high').value)
		};
		fetch('/api/v1/system/temperature', {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(payload)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('温度配置保存成功');
			} else {
				alert(data.message || '保存失败');
			}
		});
	}

	function saveDeviceInfo() {
		var payload = {
			name: document.getElementById('device-name').value,
			location: document.getElementById('device-location').value,
			contact: document.getElementById('contact-person').value
		};
		fetch('/api/v1/system/info', {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(payload)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('设备信息保存成功');
			} else {
				alert(data.message || '保存失败');
			}
		});
	}

	function saveDateTime() {
		var payload = {
			timezone: document.getElementById('timezone').value,
			datetime: document.getElementById('datetime').value
		};
		fetch('/api/v1/system/datetime', {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(payload)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('日期时间配置保存成功');
			} else {
				alert(data.message || '保存失败');
			}
		});
	}

	// 页面加载时获取配置
	loadConfig();
	</script>
	`

	boxContent := template.HTML(systemConfigContent)

	systemConfigBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-cogs"></i> 系统配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(systemConfigBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "系统配置",
		Description: "维护 → 系统配置",
	}, nil
}
