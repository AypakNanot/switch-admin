package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSTPContent STP 配置页面
func getSTPContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	stpContent := `
<style>
	.stp-section { margin-bottom: 30px; }
	.stp-table { font-size: 13px; }
	.stp-table .status-forwarding { color: #28a745; font-weight: bold; }
	.stp-table .status-blocking { color: #dc3545; }
	.stp-table .status-listening { color: #f0ad4e; }
	.stp-table .status-learning { color: #5bc0de; }
	.role-root { color: #28a745; }
	.role-designated { color: #337ab7; }
	.role-alternate { color: #f0ad4e; }
	.role-disabled { color: #dc3545; }
</style>

<div class="stp-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">STP 基本配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveSTPConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">STP 模式</label>
				<div class="col-sm-3">
					<select class="form-control" id="stp-mode">
						<option value="disabled">禁用</option>
						<option value="stp">STP (802.1D)</option>
						<option value="rstp">RSTP (802.1w)</option>
						<option value="mstp">MSTP (802.1s)</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">交换机优先级</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="stp-priority" min="0" max="61440" step="4096" value="32768">
					<small class="text-muted">4096 的倍数，值越小优先级越高</small>
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-warning btn-sm" onclick="resetSTPDefault()">
						<i class="fa fa-undo"></i> 恢复默认
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">Hello Time</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="stp-hello-time" min="1" max="10" value="2">
						<span class="input-group-addon">秒</span>
					</div>
				</div>
				<label class="col-sm-2 control-label">Forward Delay</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="stp-forward-delay" min="4" max="30" value="15">
						<span class="input-group-addon">秒</span>
					</div>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">Max Age</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="stp-max-age" min="6" max="40" value="20">
						<span class="input-group-addon">秒</span>
					</div>
				</div>
				<label class="col-sm-2 control-label">Root 保护</label>
				<div class="col-sm-3">
					<select class="form-control" id="stp-root-protection">
						<option value="disable">禁用</option>
						<option value="enable">启用</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="stp-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">STP 状态信息</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshSTPStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-trophy"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">根桥 ID</span>
							<span class="info-box-number mono-font" id="root-bridge">32768.001A.2BFF.3C4D</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-arrow-right"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">根端口</span>
							<span class="info-box-number" id="root-port">-</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-random"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">到根桥路径开销</span>
							<span class="info-box-number" id="root-cost">0</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="stp-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">STP 端口状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshSTPStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-warning btn-sm" onclick="clearSTPException()">
					<i class="fa fa-exclamation-triangle"></i> 清除异常
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover stp-table">
				<thead>
					<tr>
						<th width="100">端口</th>
						<th width="120">角色</th>
						<th width="100">状态</th>
						<th width="100">优先级</th>
						<th width="120">路径开销</th>
						<th width="150">设计桥 ID</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="stp-port-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="role-designated">指定端口</span></td>
						<td><span class="status-forwarding">Forwarding</span></td>
						<td>128</td>
						<td>20000</td>
						<td class="mono-font">32768.001A.2BFF.3C4D</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortSTPDetail('GE1/0/1')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="disablePortSTP('GE1/0/1')">禁用</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="role-designated">指定端口</span></td>
						<td><span class="status-forwarding">Forwarding</span></td>
						<td>128</td>
						<td>20000</td>
						<td class="mono-font">32768.001A.2BFF.3C4D</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortSTPDetail('GE1/0/2')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="disablePortSTP('GE1/0/2')">禁用</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/23</td>
						<td><span class="role-root">根端口</span></td>
						<td><span class="status-forwarding">Forwarding</span></td>
						<td>128</td>
						<td>10000</td>
						<td class="mono-font">4096.0000.11FF.2233</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortSTPDetail('GE1/0/23')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="disablePortSTP('GE1/0/23')">禁用</button>
						</td>
					</tr>
					<tr>
						<td>GE1/0/24</td>
						<td><span class="role-alternate">Alternate</span></td>
						<td><span class="status-blocking">Blocking</span></td>
						<td>128</td>
						<td>20000</td>
						<td class="mono-font">32768.001A.2BFF.3C4D</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortSTPDetail('GE1/0/24')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="disablePortSTP('GE1/0/24')">禁用</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口 STP 详情弹窗 -->
<div class="modal fade" id="portSTPModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">STP 端口详情 - <span id="port-stp-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="port-stp-form">
					<div class="box box-default">
						<div class="box-header">
							<h4 class="box-title">端口 STP 参数</h4>
						</div>
						<div class="box-body">
							<table class="table">
								<tr><td width="150">端口 ID</td><td id="detail-port-id">-</td></tr>
								<tr><td>角色</td><td id="detail-port-role">-</td></tr>
								<tr><td>状态</td><td id="detail-port-state">-</td></tr>
								<tr><td>优先级</td><td id="detail-port-priority">-</td></tr>
								<tr><td>路径开销</td><td id="detail-port-cost">-</td></tr>
								<tr><td>设计桥 ID</td><td id="detail-port-designated-bridge">-</td></tr>
								<tr><td>设计端口 ID</td><td id="detail-port-designated-port">-</td></tr>
								<tr><td>TC 收到次数</td><td id="detail-port-tc-received">-</td></tr>
								<tr><td>TC 发送次数</td><td id="detail-port-tc-sent">-</td></tr>
							</table>
						</div>
					</div>
					<div class="box box-warning">
						<div class="box-header">
							<h4 class="box-title">端口 STP 特性</h4>
						</div>
						<div class="box-body">
							<div class="form-group">
								<label>边缘端口</label>
								<select class="form-control" id="port-edge">
									<option value="disable">禁用</option>
									<option value="enable">启用</option>
								</select>
								<small class="text-muted">边缘端口直接跳转到 Forwarding 状态，用于连接终端设备</small>
							</div>
							<div class="form-group">
								<label>BPDU 保护</label>
								<select class="form-control" id="port-bpdu-protection">
									<option value="disable">禁用</option>
									<option value="enable">启用</option>
								</select>
								<small class="text-muted">边缘端口收到 BPDU 时关闭端口，防止环路</small>
							</div>
							<div class="form-group">
								<label>根保护</label>
								<select class="form-control" id="port-root-protection">
									<option value="disable">禁用</option>
									<option value="enable">启用</option>
								</select>
								<small class="text-muted">防止该端口成为根端口</small>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortSTPConfig()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshSTPStatus() {
	location.reload();
}

function resetSTPDefault() {
	if (confirm('确定要恢复 STP 默认配置吗？')) {
		$.ajax({
			url: '/api/v1/config/stp/reset',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('已恢复默认配置', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
				} else {
					layer.msg(res.message || '恢复失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('恢复失败，请重试', {icon: 2});
			}
		});
	}
}

function saveSTPConfig() {
	var data = {
		mode: $('#stp-mode').val(),
		priority: parseInt($('#stp-priority').val()),
		hello_time: parseInt($('#stp-hello-time').val()),
		forward_delay: parseInt($('#stp-forward-delay').val()),
		max_age: parseInt($('#stp-max-age').val()),
		root_protection: $('#stp-root-protection').val() === 'enable'
	};

	if (data.priority % 4096 !== 0) {
		layer.msg('交换机优先级必须是 4096 的倍数', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/stp/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('STP 配置保存成功', {icon: 1});
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

function clearSTPException() {
	if (confirm('确定要清除 STP 异常状态吗？')) {
		$.ajax({
			url: '/api/v1/config/stp/exception',
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('异常清除成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
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

function viewPortSTPDetail(portId) {
	$('#port-stp-id').text(portId);
	$('#portSTPModal').modal('show');
	// TODO: 加载端口 STP 详情
	$('#detail-port-id').text(portId);
	$('#detail-port-role').text('指定端口');
	$('#detail-port-state').text('Forwarding');
	$('#detail-port-priority').text('128');
	$('#detail-port-cost').text('20000');
	$('#detail-port-designated-bridge').text('32768.001A.2BFF.3C4D');
	$('#detail-port-designated-port').text('128');
	$('#detail-port-tc-received').text('5');
	$('#detail-port-tc-sent').text('12');
}

function savePortSTPConfig() {
	var portId = $('#port-stp-id').text();
	var data = {
		edge: $('#port-edge').val() === 'enable',
		bpdu_protection: $('#port-bpdu-protection').val() === 'enable',
		root_protection: $('#port-root-protection').val() === 'enable'
	};

	$.ajax({
		url: '/api/v1/config/stp/port/' + portId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('端口 STP 配置保存成功', {icon: 1});
				$('#portSTPModal').modal('hide');
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function disablePortSTP(portId) {
	if (confirm('确定要禁用端口 ' + portId + ' 的 STP 功能吗？这可能导致网络环路！')) {
		$.ajax({
			url: '/api/v1/config/stp/port/' + portId + '/disable',
			type: 'PUT',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('STP 已禁用', {icon: 1});
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
}
</script>
`

	boxContent := template.HTML(stpContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-random"></i> STP 配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "STP 配置",
		Description: "配置生成树协议 (STP/RSTP/MSTP)，防止网络环路，支持根桥选举、端口角色管理、边缘端口等特性",
	}

	return panel, nil
}
