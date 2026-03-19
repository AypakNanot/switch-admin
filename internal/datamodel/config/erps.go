package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getERPSContent ERPS 配置页面
func getERPSContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	erpsContent := `
<style>
	.erps-section { margin-bottom: 30px; }
	.erps-table { font-size: 13px; }
	.erps-table .status-active { color: #28a745; font-weight: bold; }
	.erps-table .status-standby { color: #f0ad4e; }
	.erps-table .status-down { color: #dc3545; }
	.topo-diagram { background: #f5f5f5; padding: 20px; border-radius: 4px; text-align: center; }
	.node { display: inline-block; width: 80px; height: 80px; border-radius: 50%; background: #337ab7; color: white; line-height: 80px; margin: 10px; }
	.node.master { background: #28a745; }
	.node.slave { background: #f0ad4e; }
	.link { display: inline-block; width: 100px; height: 3px; background: #999; vertical-align: middle; margin: 0 10px; }
	.link.rpl { background: #dc3545; }
	.link.active { background: #28a745; }
</style>

<div class="erps-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">ERPS 环网配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveERPSConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
				<button type="button" class="btn btn-warning btn-sm" onclick="forceSwitch()">
					<i class="fa fa-exchange"></i> 强制切换
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">ERPS 状态</label>
				<div class="col-sm-3">
					<select class="form-control" id="erps-enabled">
						<option value="false">禁用</option>
						<option value="true">启用</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">环网 ID</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="erps-ring-id" min="1" max="16" value="1">
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-default btn-sm" onclick="testERPS()">
						<i class="fa fa-flask"></i> 测试
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">控制 VLAN</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="erps-control-vlan" min="1" max="4094" value="4000">
					<small class="text-muted">用于传输 ERPS 协议报文的 VLAN</small>
				</div>
				<label class="col-sm-2 control-label">数据 VLAN</label>
				<div class="col-sm-3">
					<input type="text" class="form-control" id="erps-data-vlan" placeholder="多个用逗号分隔，如：10,20,30">
					<small class="text-muted">受环网保护的业务 VLAN</small>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">节点角色</label>
				<div class="col-sm-3">
					<select class="form-control" id="erps-role">
						<option value="auto">自动协商</option>
						<option value="master">Master 节点</option>
						<option value="slave">Slave 节点</option>
						<option value="transit">Transit 节点</option>
					</select>
					<small class="text-muted">Master/Slave 为 RPL 所有者节点</small>
				</div>
				<label class="col-sm-2 control-label">WTR 时间</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="erps-wtr" min="1" max="12" value="5">
						<span class="input-group-addon">分钟</span>
					</div>
					<small class="text-muted">等待恢复时间，故障恢复后等待切换回原链路的时间</small>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="erps-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">环网拓扑</h3>
		</div>
		<div class="box-body">
			<div class="topo-diagram">
				<div class="node master">节点 A<br><small>(Master)</small></div>
				<div class="link active"></div>
				<div class="node">节点 B</div>
				<div class="link active"></div>
				<div class="node">节点 C</div>
				<div class="link active"></div>
				<div class="node slave">节点 D<br><small>(Slave)</small></div>
				<div class="link rpl" id="rpl-link">RPL<br><small>(阻塞)</small></div>
				<div style="margin-top: 20px; color: #666;">
					<i class="fa fa-circle" style="color: #28a745;"></i> 活动链路
					<i class="fa fa-circle" style="color: #dc3545; margin-left: 15px;"></i> RPL 阻塞链路
					<i class="fa fa-circle" style="color: #f0ad4e; margin-left: 15px;"></i> 备用链路
				</div>
			</div>
		</div>
	</div>
</div>

<div class="erps-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">ERPS 环网端口</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshERPSStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover erps-table">
				<thead>
					<tr>
						<th width="100">端口</th>
						<th width="100">角色</th>
						<th width="100">状态</th>
						<th width="120">邻居</th>
						<th width="100">信号失效</th>
						<th width="100">信号劣化</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="erps-port-body">
					<tr>
						<td>GE1/0/23</td>
						<td>东向端口</td>
						<td><span class="status-active">活动</span></td>
						<td>00:1A:2B:3C:4D:5E</td>
						<td>正常</td>
						<td>正常</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortERPSDetail('GE1/0/23')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="simulateFault('GE1/0/23')">模拟故障</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/24</td>
						<td>西向端口</td>
						<td><span class="status-standby">RPL 阻塞</span></td>
						<td>00:1A:2B:3C:4D:5F</td>
						<td>正常</td>
						<td>正常</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortERPSDetail('GE1/0/24')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="simulateFault('GE1/0/24')">模拟故障</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="erps-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">告警与事件</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="clearEvents()">
					<i class="fa fa-trash"></i> 清除事件
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">事件通知</label>
				<div class="col-sm-9">
					<label class="checkbox-inline"><input type="checkbox" class="event-checkbox" value="log" checked> 系统日志</label>
					<label class="checkbox-inline"><input type="checkbox" class="event-checkbox" value="snmp" checked> SNMP Trap</label>
					<label class="checkbox-inline"><input type="checkbox" class="event-checkbox" value="email"> 邮件通知</label>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">信号劣化阈值</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="sd-threshold" min="1" max="100" value="30">
						<span class="input-group-addon">%</span>
					</div>
					<small class="text-muted">误码率超过此阈值时触发信号劣化告警</small>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- 端口 ERPS 详情弹窗 -->
<div class="modal fade" id="portERPSModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">ERPS 端口详情 - <span id="port-erps-id"></span></h4>
			</div>
			<div class="modal-body">
				<table class="table">
					<tr><td width="150">端口</td><td id="detail-port">-</td></tr>
					<tr><td>方向</td><td id="detail-direction">-</td></tr>
					<tr><td>状态</td><td id="detail-state">-</td></tr>
					<tr><td>邻居 MAC</td><td id="detail-neighbor">-</td></tr>
					<tr><td>收到 SF 次数</td><td id="detail-sf-count">-</td></tr>
					<tr><td>收到 SD 次数</td><td id="detail-sd-count">-</td></tr>
					<tr><td>最后一次 SF 时间</td><td id="detail-last-sf">-</td></tr>
					<tr><td>最后一次 SD 时间</td><td id="detail-last-sd">-</td></tr>
				</table>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshERPSStatus() {
	location.reload();
}

function saveERPSConfig() {
	var data = {
		enabled: $('#erps-enabled').val() === 'true',
		ring_id: parseInt($('#erps-ring-id').val()),
		control_vlan: parseInt($('#erps-control-vlan').val()),
		data_vlans: $('#erps-data-vlan').val().split(',').map(function(v) { return parseInt(v.trim()); }).filter(function(v) { return !isNaN(v); }),
		role: $('#erps-role').val(),
		wtr: parseInt($('#erps-wtr').val())
	};

	$.ajax({
		url: '/api/v1/config/erps/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('ERPS 配置保存成功', {icon: 1});
				setTimeout(function() { location.reload(); }, 1000);
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function forceSwitch() {
	if (confirm('确定要执行 ERPS 强制切换吗？这将导致网络短暂中断。')) {
		$.ajax({
			url: '/api/v1/config/erps/force-switch',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('强制切换成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
				} else {
					layer.msg(res.message || '切换失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('切换失败，请重试', {icon: 2});
			}
		});
	}
}

function testERPS() {
	layer.msg('ERPS 测试功能待实现', {icon: 3});
}

function viewPortERPSDetail(portId) {
	$('#port-erps-id').text(portId);
	$('#portERPSModal').modal('show');
	// TODO: 加载端口 ERPS 详情
	$('#detail-port').text(portId);
	$('#detail-direction').text('东向端口');
	$('#detail-state').text('活动');
	$('#detail-neighbor').text('00:1A:2B:3C:4D:5E');
	$('#detail-sf-count').text('0');
	$('#detail-sd-count').text('0');
	$('#detail-last-sf').text('-');
	$('#detail-last-sd').text('-');
}

function simulateFault(portId) {
	if (confirm('确定要模拟端口 ' + portId + ' 的链路故障吗？这将触发 ERPS 保护切换。')) {
		$.ajax({
			url: '/api/v1/config/erps/port/' + portId + '/simulate-fault',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('故障模拟成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 2000);
				} else {
					layer.msg(res.message || '模拟失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('模拟失败，请重试', {icon: 2});
			}
		});
	}
}

function clearEvents() {
	if (confirm('确定要清除所有 ERPS 事件记录吗？')) {
		$.ajax({
			url: '/api/v1/config/erps/events',
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('事件清除成功', {icon: 1});
				} else {
					layer.msg(res.message || '清除失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('清除失败，请重试', {icon: 2});
			}
		});
	}
}
</script>
`

	boxContent := template.HTML(erpsContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-retweet"></i> ERPS 配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "ERPS 配置",
		Description: "配置以太环网保护切换 (ERPS)，实现小于 50ms 的快速环网保护。支持 Master/Slave 节点、RPL 链路管理、信号失效/劣化检测",
	}

	return panel, nil
}
