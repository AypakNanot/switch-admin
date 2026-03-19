package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPoEContent PoE 配置页面
func getPoEContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	poeContent := `
<style>
	.poe-section { margin-bottom: 30px; }
	.poe-table { font-size: 13px; }
	.poe-table .status-delivering { color: #28a745; font-weight: bold; }
	.poe-table .status-searching { color: #f0ad4e; }
	.poe-table .status-fault { color: #dc3545; }
	.poe-table .status-disabled { color: #999; }
	.power-bar { width: 100%; height: 20px; background: #f0f0f0; border-radius: 3px; overflow: hidden; }
	.power-bar-fill { height: 100%; background: linear-gradient(to right, #28a745, #ffc107, #dc3545); transition: width 0.3s; }
	.info-box-power { min-height: 80px; }
</style>

<div class="poe-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">PoE 全局配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="savePoEConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">PoE 状态</label>
				<div class="col-sm-3">
					<select class="form-control" id="poe-enabled">
						<option value="true">启用</option>
						<option value="false">关闭</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">PoE 标准</label>
				<div class="col-sm-3">
					<select class="form-control" id="poe-standard">
						<option value="802.3af">802.3af (PoE)</option>
						<option value="802.3at">802.3at (PoE+)</option>
						<option value="802.3bt">802.3bt (PoE++)</option>
						<option value="legacy">Legacy (24V)</option>
					</select>
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-default btn-sm" onclick="resetPoEDefault()">
						<i class="fa fa-undo"></i> 恢复默认
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">管理方式</label>
				<div class="col-sm-3">
					<select class="form-control" id="poe-mode">
						<option value="auto">自动模式</option>
						<option value="static">静态模式</option>
					</select>
					<small class="text-muted">自动模式：按端口顺序分配功率；静态模式：按优先级分配</small>
				</div>
				<label class="col-sm-2 control-label">功率分配优先级</label>
				<div class="col-sm-3">
					<select class="form-control" id="poe-priority">
						<option value="critical">Critical (最高)</option>
						<option value="high">High (高)</option>
						<option value="low" selected>Low (低)</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">最大功耗限制</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="poe-max-power" min="154" max="30000" value="30000">
						<span class="input-group-addon">mW</span>
					</div>
					<small class="text-muted">整机最大 PoE 功率预算</small>
				</div>
				<label class="col-sm-2 control-label">保留功率</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="poe-reserved-power" min="0" max="30000" value="0">
						<span class="input-group-addon">mW</span>
					</div>
					<small class="text-muted">为高优先级端口保留的功率</small>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="poe-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">PoE 功率预算</h3>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-3">
					<div class="info-box info-box-power">
						<span class="info-box-icon bg-blue"><i class="fa fa-bolt"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">总功率预算</span>
							<span class="info-box-number">30 W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box info-box-power">
						<span class="info-box-icon bg-green"><i class="fa fa-plug"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">已用功率</span>
							<span class="info-box-number">18.5 W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box info-box-power">
						<span class="info-box-icon bg-yellow"><i class="fa fa-battery-half"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">剩余功率</span>
							<span class="info-box-number">11.5 W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box info-box-power">
						<span class="info-box-icon bg-red"><i class="fa fa-tachometer"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">使用率</span>
							<span class="info-box-number">61.7%</span>
						</div>
					</div>
				</div>
			</div>
			<div class="power-bar" style="margin-top: 10px;">
				<div class="power-bar-fill" style="width: 61.7%;"></div>
			</div>
			<div style="text-align: center; margin-top: 5px; color: #666; font-size: 12px;">
				功率使用率：61.7% (18.5W / 30W)
			</div>
		</div>
	</div>
</div>

<div class="poe-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">PoE 端口状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshPoEStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-warning btn-sm" onclick="resetPoEPorts()">
					<i class="fa fa-power-off"></i> 批量重启
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover poe-table">
				<thead>
					<tr>
						<th width="40"><input type="checkbox" id="select-all"></th>
						<th width="100">端口</th>
						<th width="100">状态</th>
						<th width="100">PD 类型</th>
						<th width="100">功率等级</th>
						<th width="100">供电电压</th>
						<th width="100">供电电流</th>
						<th width="100">实际功率</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="poe-port-body">
					<tr>
						<td><input type="checkbox" class="poe-checkbox" value="GE1/0/1"></td>
						<td>GE1/0/1</td>
						<td><span class="status-delivering">供电中</span></td>
						<td>AP</td>
						<td>Class 3</td>
						<td>53.2 V</td>
						<td>0.28 A</td>
						<td>14.9 W</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPoEPortDetail('GE1/0/1')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="cyclePoEPort('GE1/0/1')">重启</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="poe-checkbox" value="GE1/0/2"></td>
						<td>GE1/0/2</td>
						<td><span class="status-delivering">供电中</span></td>
						<td>IP Camera</td>
						<td>Class 2</td>
						<td>53.5 V</td>
						<td>0.15 A</td>
						<td>8.0 W</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPoEPortDetail('GE1/0/2')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="cyclePoEPort('GE1/0/2')">重启</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="poe-checkbox" value="GE1/0/3"></td>
						<td>GE1/0/3</td>
						<td><span class="status-searching">检测中</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>0 W</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPoEPortDetail('GE1/0/3')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="cyclePoEPort('GE1/0/3')">重启</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="poe-checkbox" value="GE1/0/4"></td>
						<td>GE1/0/4</td>
						<td><span class="status-disabled">未供电</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>0 W</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPoEPortDetail('GE1/0/4')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="cyclePoEPort('GE1/0/4')">重启</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="poe-checkbox" value="GE1/0/5"></td>
						<td>GE1/0/5</td>
						<td><span class="status-fault">故障</span></td>
						<td>Unknown</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>0 W</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPoEPortDetail('GE1/0/5')">详情</button>
							<button class="btn btn-sm btn-warning" onclick="cyclePoEPort('GE1/0/5')">重启</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- PoE 端口配置弹窗 -->
<div class="modal fade" id="poePortModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">PoE 端口配置 - <span id="poe-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="poe-port-form">
					<div class="box box-default">
						<div class="box-body">
							<div class="form-group">
								<label>端口供电状态</label>
								<select class="form-control" id="port-poe-enabled">
									<option value="true">启用</option>
									<option value="false">禁用</option>
								</select>
							</div>
							<div class="form-group">
								<label>最大供电功率</label>
								<div class="input-group">
									<input type="number" class="form-control" id="port-max-power" min="0" max="30000" value="30000">
									<span class="input-group-addon">mW</span>
								</div>
							</div>
							<div class="form-group">
								<label>功率优先级</label>
								<select class="form-control" id="port-priority">
									<option value="critical">Critical (最高)</option>
									<option value="high">High (高)</option>
									<option value="low">Low (低)</option>
								</select>
							</div>
							<div class="form-group">
								<label>PD 名称</label>
								<input type="text" class="form-control" id="pd-name" maxlength="32" placeholder="可选">
								<small class="text-muted">用于标识连接的受电设备</small>
							</div>
							<div class="alert alert-info">
								<i class="fa fa-info-circle"></i> 端口统计:
								<ul style="margin: 5px 0 0 20px;">
									<li>累计供电时间：15 天 8 小时 30 分</li>
									<li>累计消耗电量：2.5 kWh</li>
									<li>平均功率：12.3 W</li>
								</ul>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortPoEConfig()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.poe-checkbox').prop('checked', this.checked);
});

function refreshPoEStatus() {
	location.reload();
}

function resetPoEDefault() {
	if (confirm('确定要恢复 PoE 默认配置吗？')) {
		$.ajax({
			url: '/api/v1/config/poe/reset',
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

function savePoEConfig() {
	var data = {
		enabled: $('#poe-enabled').val() === 'true',
		standard: $('#poe-standard').val(),
		mode: $('#poe-mode').val(),
		priority: $('#poe-priority').val(),
		max_power: parseInt($('#poe-max-power').val()),
		reserved_power: parseInt($('#poe-reserved-power').val())
	};

	$.ajax({
		url: '/api/v1/config/poe/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('PoE 配置保存成功', {icon: 1});
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

function viewPoEPortDetail(portId) {
	$('#poe-port-id').text(portId);
	$('#poePortModal').modal('show');
	// TODO: 加载端口 PoE 详情
}

function savePortPoEConfig() {
	var portId = $('#poe-port-id').text();
	var data = {
		enabled: $('#port-poe-enabled').val() === 'true',
		max_power: parseInt($('#port-max-power').val()),
		priority: $('#port-priority').val(),
		pd_name: $('#pd-name').val()
	};

	$.ajax({
		url: '/api/v1/config/poe/port/' + portId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('端口配置保存成功', {icon: 1});
				$('#poePortModal').modal('hide');
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

function cyclePoEPort(portId) {
	if (confirm('确定要重启端口 ' + portId + ' 的 PoE 供电吗？这将导致 PD 设备短暂断电。')) {
		$.ajax({
			url: '/api/v1/config/poe/port/' + portId + '/cycle',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('PoE 重启成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
				} else {
					layer.msg(res.message || '重启失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('重启失败，请重试', {icon: 2});
			}
		});
	}
}

function resetPoEPorts() {
	var selected = [];
	$('.poe-checkbox:checked').each(function() {
		selected.push($(this).val());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要重启的端口', {icon: 2});
		return;
	}
	if (confirm('确定要重启选中的 ' + selected.length + ' 个端口的 PoE 供电吗？')) {
		layer.msg('批量重启功能待实现', {icon: 3});
	}
}
</script>
`

	boxContent := template.HTML(poeContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-plug"></i> PoE 配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "PoE 配置",
		Description: "配置以太网供电 (PoE/PoE+/PoE++) 功能，管理功率预算、端口供电状态、PD 设备信息，支持 802.3af/at/bt 标准",
	}

	return panel, nil
}
