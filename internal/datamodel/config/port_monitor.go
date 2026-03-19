package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPortMonitorContent 端口监控页面
func getPortMonitorContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	monitorContent := `
<style>
	.monitor-section { margin-bottom: 30px; }
	.monitor-table { font-size: 13px; }
	.monitor-table .status-enabled { color: #28a745; font-weight: bold; }
	.monitor-table .status-disabled { color: #dc3545; }
	.config-box { padding: 15px; background: #f9f9f9; border-radius: 4px; margin-bottom: 15px; }
	.rate-info { font-size: 12px; color: #666; }
</style>

<div class="monitor-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口监控配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveMonitorConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">监控模式</label>
				<div class="col-sm-4">
					<select class="form-control" id="monitor-mode">
						<option value="disabled">禁用</option>
						<option value="local">本地监控</option>
						<option value="remote">远程监控 (RMON)</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">采样间隔</label>
				<div class="col-sm-4">
					<select class="form-control" id="sample-interval">
						<option value="10">10 秒</option>
						<option value="30">30 秒</option>
						<option value="60" selected>60 秒</option>
						<option value="300">5 分钟</option>
						<option value="1800">30 分钟</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">监控 statistic 类型</label>
				<div class="col-sm-9">
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="packets" checked> 数据包数</label>
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="bytes" checked> 字节数</label>
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="errors"> 错误包</label>
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="discards"> 丢弃包</label>
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="broadcast"> 广播包</label>
					<label class="checkbox-inline"><input type="checkbox" class="stat-checkbox" value="multicast"> 组播包</label>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="monitor-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">端口监控状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshMonitorStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-success btn-sm" onclick="exportMonitorData()">
					<i class="fa fa-download"></i> 导出数据
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover monitor-table">
				<thead>
					<tr>
						<th width="100">端口</th>
						<th width="80">状态</th>
						<th width="120">入方向速率</th>
						<th width="120">出方向速率</th>
						<th width="100">入包数</th>
						<th width="100">出包数</th>
						<th width="80">错误</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="monitor-table-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-enabled">监控中</span></td>
						<td>125.5 Mbps</td>
						<td>89.2 Mbps</td>
						<td>1523456</td>
						<td>987654</td>
						<td>0</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="viewPortDetail('GE1/0/1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-enabled">监控中</span></td>
						<td>45.8 Mbps</td>
						<td>32.1 Mbps</td>
						<td>654321</td>
						<td>456789</td>
						<td>2</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="viewPortDetail('GE1/0/2')">详情</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td><span class="status-disabled">未监控</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="enablePortMonitor('GE1/0/3')">启用</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td><span class="status-enabled">监控中</span></td>
						<td>78.3 Mbps</td>
						<td>56.7 Mbps</td>
						<td>891234</td>
						<td>567890</td>
						<td>0</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="viewPortDetail('GE1/0/4')">详情</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="monitor-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">告警阈值配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveAlarmConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">带宽利用率告警</label>
				<div class="col-sm-4">
					<div class="input-group">
						<input type="number" class="form-control" id="bandwidth-alarm-threshold" value="80" min="1" max="100">
						<span class="input-group-addon">%</span>
					</div>
					<small class="text-muted">当端口带宽利用率超过此阈值时触发告警</small>
				</div>
				<label class="col-sm-2 control-label">错误包率告警</label>
				<div class="col-sm-4">
					<div class="input-group">
						<input type="number" class="form-control" id="error-alarm-threshold" value="1" min="0" max="100">
						<span class="input-group-addon">%</span>
					</div>
					<small class="text-muted">当错误包占比超过此阈值时触发告警</small>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">告警通知方式</label>
				<div class="col-sm-9">
					<label class="checkbox-inline"><input type="checkbox" class="alarm-checkbox" value="log" checked> 系统日志</label>
					<label class="checkbox-inline"><input type="checkbox" class="alarm-checkbox" value="snmp" checked> SNMP Trap</label>
					<label class="checkbox-inline"><input type="checkbox" class="alarm-checkbox" value="email"> 邮件通知</label>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- 端口详情弹窗 -->
<div class="modal fade" id="portDetailModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">端口监控详情 - <span id="detail-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<div class="row">
					<div class="col-md-6">
						<div class="box box-default">
							<div class="box-header">
								<h4 class="box-title">实时速率</h4>
							</div>
							<div class="box-body">
								<div id="rate-chart" style="height: 200px; background: #f5f5f5; display: flex; align-items: center; justify-content: center;">
									图表区域 (需集成图表库)
								</div>
							</div>
						</div>
					</div>
					<div class="col-md-6">
						<div class="box box-default">
							<div class="box-header">
								<h4 class="box-title">统计信息</h4>
							</div>
							<div class="box-body">
								<table class="table table-sm">
									<tr><td>入方向字节数</td><td id="stat-in-bytes">0</td></tr>
									<tr><td>出方向字节数</td><td id="stat-out-bytes">0</td></tr>
									<tr><td>入方向包数</td><td id="stat-in-packets">0</td></tr>
									<tr><td>出方向包数</td><td id="stat-out-packets">0</td></tr>
									<tr><td>广播包数</td><td id="stat-broadcast">0</td></tr>
									<tr><td>组播包数</td><td id="stat-multicast">0</td></tr>
									<tr><td>错误包数</td><td id="stat-errors">0</td></tr>
									<tr><td>丢弃包数</td><td id="stat-discards">0</td></tr>
								</table>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshMonitorStatus() {
	location.reload();
}

function viewPortDetail(portId) {
	$('#detail-port-id').text(portId);
	$('#portDetailModal').modal('show');
	// TODO: 加载端口详细统计信息
	loadPortStats(portId);
}

function loadPortStats(portId) {
	// 模拟数据，实际应从 API 获取
	$('#stat-in-bytes').text('1,523,456,789');
	$('#stat-out-bytes').text('987,654,321');
	$('#stat-in-packets').text('1,523,456');
	$('#stat-out-packets').text('987,654');
	$('#stat-broadcast').text('15,234');
	$('#stat-multicast').text('45,678');
	$('#stat-errors').text('0');
	$('#stat-discards').text('12');
}

function enablePortMonitor(portId) {
	$.ajax({
		url: '/api/v1/config/port-monitor/' + portId + '/enable',
		type: 'PUT',
		success: function(res) {
			if (res.code === 200) {
				layer.msg('已启用监控', {icon: 1});
				setTimeout(function() { location.reload(); }, 1000);
			} else {
				layer.msg(res.message || '操作失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('操作失败，请重试', {icon: 2});
		}
	});
}

function saveMonitorConfig() {
	var stats = [];
	$('.stat-checkbox:checked').each(function() {
		stats.push($(this).val());
	});

	var data = {
		mode: $('#monitor-mode').val(),
		sample_interval: parseInt($('#sample-interval').val()),
		stat_types: stats
	};

	$.ajax({
		url: '/api/v1/config/port-monitor/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('监控配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function saveAlarmConfig() {
	var notifications = [];
	$('.alarm-checkbox:checked').each(function() {
		notifications.push($(this).val());
	});

	var data = {
		bandwidth_threshold: parseInt($('#bandwidth-alarm-threshold').val()),
		error_threshold: parseInt($('#error-alarm-threshold').val()),
		notifications: notifications
	};

	$.ajax({
		url: '/api/v1/config/port-monitor/alarm',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('告警配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function exportMonitorData() {
	layer.msg('正在导出数据...', {icon: 1});
	// TODO: 实现数据导出功能
	window.location.href = '/api/v1/config/port-monitor/export';
}
</script>
`

	boxContent := template.HTML(monitorContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-line-chart"></i> 端口监控`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口监控",
		Description: "监控端口流量统计信息，包括速率、包数、错误等，支持告警阈值配置和数据导出",
	}

	return panel, nil
}
