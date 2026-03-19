package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getResourceContent 资源管理页面
func getResourceContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	resourceContent := `
<style>
	.resource-section { margin-bottom: 30px; }
	.resource-table { font-size: 13px; }
	.resource-bar { width: 100%; height: 25px; background: #f0f0f0; border-radius: 3px; overflow: hidden; margin: 5px 0; }
	.resource-bar-fill { height: 100%; transition: width 0.3s; }
	.bar-green { background: linear-gradient(to right, #28a745, #5cb85c); }
	.bar-yellow { background: linear-gradient(to right, #ffc107, #ff9800); }
	.bar-red { background: linear-gradient(to right, #f44336, #d32f2f); }
	.resource-item { padding: 15px; border: 1px solid #ddd; border-radius: 4px; margin-bottom: 15px; }
	.resource-label { font-size: 14px; color: #666; margin-bottom: 5px; }
	.resource-value { font-size: 20px; font-weight: bold; color: #333; }
	.alarm-normal { color: #28a745; }
	.alarm-warning { color: #ffc107; }
	.alarm-critical { color: #dc3545; }
</style>

<div class="resource-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">系统资源概览</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshResource()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-3">
					<div class="resource-item" style="text-align: center;">
						<div class="resource-label">CPU 使用率</div>
						<div class="resource-value" style="color: #28a745;">23%</div>
						<div class="resource-bar">
							<div class="resource-bar-fill bar-green" style="width: 23%;"></div>
						</div>
						<small class="text-muted">1 分钟平均</small>
					</div>
				</div>
				<div class="col-md-3">
					<div class="resource-item" style="text-align: center;">
						<div class="resource-label">内存使用率</div>
						<div class="resource-value" style="color: #ffc107;">58%</div>
						<div class="resource-bar">
							<div class="resource-bar-fill bar-yellow" style="width: 58%;"></div>
						</div>
						<small class="text-muted">已用 348MB / 600MB</small>
					</div>
				</div>
				<div class="col-md-3">
					<div class="resource-item" style="text-align: center;">
						<div class="resource-label">Flash 使用率</div>
						<div class="resource-value" style="color: #28a745;">35%</div>
						<div class="resource-bar">
							<div class="resource-bar-fill bar-green" style="width: 35%;"></div>
						</div>
						<small class="text-muted">已用 1.4GB / 4GB</small>
					</div>
				</div>
				<div class="col-md-3">
					<div class="resource-item" style="text-align: center;">
						<div class="resource-label">TCAM 使用率</div>
						<div class="resource-value" style="color: #ffc107;">67%</div>
						<div class="resource-bar">
							<div class="resource-bar-fill bar-yellow" style="width: 67%;"></div>
						</div>
						<small class="text-muted">ACL + QoS 规则</small>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="resource-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">端口带宽利用率</h3>
			<div class="box-tools">
				<select class="form-control" style="display: inline-block; width: 150px;" id="port-bandwidth-filter">
					<option value="all">全部端口</option>
					<option value="ge">千兆端口</option>
					<option value="xe">万兆端口</option>
				</select>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover resource-table">
				<thead>
					<tr>
						<th width="100">端口</th>
						<th>带宽利用率</th>
						<th width="150">入方向</th>
						<th width="150">出方向</th>
						<th width="100">状态</th>
					</tr>
				</thead>
				<tbody id="port-bandwidth-body">
					<tr>
						<td>GE1/0/1</td>
						<td>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 25%;"></div>
							</div>
						</td>
						<td>250 Mbps (25%)</td>
						<td>180 Mbps (18%)</td>
						<td><span class="alarm-normal">正常</span></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-yellow" style="width: 75%;"></div>
							</div>
						</td>
						<td>650 Mbps (65%)</td>
						<td>750 Mbps (75%)</td>
						<td><span class="alarm-warning">繁忙</span></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 12%;"></div>
							</div>
						</td>
						<td>100 Mbps (10%)</td>
						<td>120 Mbps (12%)</td>
						<td><span class="alarm-normal">正常</span></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-red" style="width: 92%;"></div>
							</div>
						</td>
						<td>880 Mbps (88%)</td>
						<td>920 Mbps (92%)</td>
						<td><span class="alarm-critical">拥塞</span></td>
					</tr>
					<tr>
						<td>GE1/0/24</td>
						<td>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 45%;"></div>
							</div>
						</td>
						<td>400 Mbps (40%)</td>
						<td>450 Mbps (45%)</td>
						<td><span class="alarm-normal">正常</span></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="resource-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">表项资源</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshTableResource()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">MAC 地址表</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">1,245 / 16,384</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 8%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">动态：1,200</span>
								<span class="label label-info">静态：45</span>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">ARP 表</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">856 / 8,192</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 10%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">动态：820</span>
								<span class="label label-info">静态：36</span>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">路由表</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">128 / 4,096</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 3%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">直连：12</span>
								<span class="label label-info">静态：25</span>
								<span class="label label-primary">OSPF：91</span>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="row" style="margin-top: 15px;">
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">ACL 规则</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">45 / 1,024</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 4%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">IPv4: 35</span>
								<span class="label label-info">IPv6: 10</span>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">QoS 策略</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">28 / 512</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 5%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">队列：16</span>
								<span class="label label-info">整形：12</span>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">VLAN 资源</h4>
						</div>
						<div class="box-body">
							<div class="resource-label">已用 / 总计</div>
							<div class="resource-value">15 / 4,094</div>
							<div class="resource-bar">
								<div class="resource-bar-fill bar-green" style="width: 0.4%;"></div>
							</div>
							<div style="margin-top: 10px;">
								<span class="label label-success">Active: 12</span>
								<span class="label label-warning">Suspended: 3</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="resource-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">告警阈值配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveAlarmThreshold()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="form-group">
						<label>CPU 告警阈值 (%)</label>
						<div class="form-group">
							<label>警告</label>
							<input type="number" class="form-control" id="cpu-warning" min="1" max="100" value="70">
						</div>
						<div class="form-group">
							<label>严重</label>
							<input type="number" class="form-control" id="cpu-critical" min="1" max="100" value="90">
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="form-group">
						<label>内存告警阈值 (%)</label>
						<div class="form-group">
							<label>警告</label>
							<input type="number" class="form-control" id="mem-warning" min="1" max="100" value="75">
						</div>
						<div class="form-group">
							<label>严重</label>
							<input type="number" class="form-control" id="mem-critical" min="1" max="100" value="90">
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="form-group">
						<label>端口带宽告警阈值 (%)</label>
						<div class="form-group">
							<label>警告</label>
							<input type="number" class="form-control" id="bw-warning" min="1" max="100" value="80">
						</div>
						<div class="form-group">
							<label>严重</label>
							<input type="number" class="form-control" id="bw-critical" min="1" max="100" value="95">
						</div>
					</div>
				</div>
			</div>
			<div class="form-group">
				<label>告警通知方式</label>
				<div>
					<label class="checkbox-inline"><input type="checkbox" class="alarm-notify-checkbox" value="log" checked> 系统日志</label>
					<label class="checkbox-inline"><input type="checkbox" class="alarm-notify-checkbox" value="snmp" checked> SNMP Trap</label>
					<label class="checkbox-inline"><input type="checkbox" class="alarm-notify-checkbox" value="email"> 邮件通知</label>
				</div>
			</div>
		</div>
	</div>
</div>

<script>
function refreshResource() {
	location.reload();
}

function refreshTableResource() {
	location.reload();
}

function saveAlarmThreshold() {
	var data = {
		cpu_warning: parseInt($('#cpu-warning').val()),
		cpu_critical: parseInt($('#cpu-critical').val()),
		mem_warning: parseInt($('#mem-warning').val()),
		mem_critical: parseInt($('#mem-critical').val()),
		bw_warning: parseInt($('#bw-warning').val()),
		bw_critical: parseInt($('#bw-critical').val()),
		notifications: []
	};

	$('.alarm-notify-checkbox:checked').each(function() {
		data.notifications.push($(this).val());
	});

	$.ajax({
		url: '/api/v1/config/resource/alarm-threshold',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('告警阈值保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}
</script>
`

	boxContent := template.HTML(resourceContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-tachometer"></i> 资源管理`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "资源管理",
		Description: "监控系统资源使用情况，包括 CPU、内存、Flash、TCAM、端口带宽等，支持表项资源统计和告警阈值配置",
	}

	return panel, nil
}
