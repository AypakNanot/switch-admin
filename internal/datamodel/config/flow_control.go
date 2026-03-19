package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getFlowControlContent 流量控制页面
func getFlowControlContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	flowContent := `
<style>
	.flow-section { margin-bottom: 30px; }
	.flow-table { font-size: 13px; }
	.flow-table .status-enabled { color: #28a745; font-weight: bold; }
	.flow-table .status-disabled { color: #dc3545; }
	.config-box { padding: 15px; background: #f9f9f9; border-radius: 4px; margin-bottom: 15px; }
</style>

<div class="flow-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">全局流量控制配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveGlobalFlowConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">流量控制开关</label>
				<div class="col-sm-4">
					<select class="form-control" id="global-flow-enabled">
						<option value="true">启用</option>
						<option value="false">关闭</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">控制模式</label>
				<div class="col-sm-4">
					<select class="form-control" id="global-flow-mode">
						<option value="auto">自动协商</option>
						<option value="force">强制开启</option>
						<option value="disable">强制关闭</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">背压模式</label>
				<div class="col-sm-4">
					<select class="form-control" id="backpressure-mode">
						<option value="disable">禁用</option>
						<option value="enable">启用 (半双工)</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">Pause 帧类型</label>
				<div class="col-sm-4">
					<select class="form-control" id="pause-frame-type">
						<option value="symmetric">对称流控</option>
						<option value="asymmetric">非对称流控</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="flow-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">端口流量控制状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshFlowStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-success btn-sm" onclick="batchEnableFlow()">
					<i class="fa fa-check"></i> 批量启用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover flow-table">
				<thead>
					<tr>
						<th width="40"><input type="checkbox" id="select-all"></th>
						<th width="100">端口</th>
						<th width="100">管理状态</th>
						<th width="100">运行状态</th>
						<th width="120">协商结果</th>
						<th width="100">Pause 帧</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="flow-table-body">
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/1"></td>
						<td>GE1/0/1</td>
						<td><span class="status-enabled">启用</span></td>
						<td><span class="status-enabled">Up</span></td>
						<td>Full/On</td>
						<td>RX/TX</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortFlow('GE1/0/1')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/2"></td>
						<td>GE1/0/2</td>
						<td><span class="status-disabled">关闭</span></td>
						<td><span class="status-down">Down</span></td>
						<td>-</td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortFlow('GE1/0/2')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/3"></td>
						<td>GE1/0/3</td>
						<td><span class="status-enabled">启用</span></td>
						<td><span class="status-enabled">Up</span></td>
						<td>Full/Off</td>
						<td>None</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortFlow('GE1/0/3')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/4"></td>
						<td>GE1/0/4</td>
						<td><span class="status-enabled">启用</span></td>
						<td><span class="status-enabled">Up</span></td>
						<td>Half/Backpressure</td>
						<td>Backpressure</td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortFlow('GE1/0/4')">配置</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口流量控制配置弹窗 -->
<div class="modal fade" id="flowConfigModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">配置流量控制 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="flow-config-form">
					<div class="box box-default">
						<div class="box-body">
							<div class="form-group">
								<label>流量控制</label>
								<select class="form-control" id="port-flow-enabled">
									<option value="false">关闭</option>
									<option value="true">启用</option>
								</select>
							</div>
							<div class="form-group">
								<label>协商模式</label>
								<select class="form-control" id="port-flow-negotiation">
									<option value="auto">自动协商</option>
									<option value="force">强制开启</option>
									<option value="disable">强制关闭</option>
								</select>
							</div>
							<div class="form-group">
								<label>Pause 帧方向</label>
								<select class="form-control" id="port-pause-direction">
									<option value="both">RX & TX (双向)</option>
									<option value="rx">仅 RX (接收)</option>
									<option value="tx">仅 TX (发送)</option>
									<option value="none">None (无)</option>
								</select>
							</div>
							<div class="alert alert-info">
								<i class="fa fa-info-circle"></i> 说明：
								<ul style="margin: 5px 0 0 20px;">
									<li>RX: 本端可以接收对端的 Pause 帧并停止发送</li>
									<li>TX: 本端可以发送 Pause 帧告知对端暂停发送</li>
									<li>自动协商模式下，双方协商确定最终的流控状态</li>
								</ul>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortFlowConfig()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.port-checkbox').prop('checked', this.checked);
});

function refreshFlowStatus() {
	location.reload();
}

function editPortFlow(portId) {
	$('#edit-port-id').text(portId);
	$('#flowConfigModal').modal('show');
}

function batchEnableFlow() {
	var selected = [];
	$('.port-checkbox:checked').each(function() {
		selected.push($(this).val());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要配置的端口', {icon: 2});
		return;
	}
	if (confirm('确定要启用选中的 ' + selected.length + ' 个端口的流量控制吗？')) {
		layer.msg('批量启用功能待实现', {icon: 3});
	}
}

function saveGlobalFlowConfig() {
	var data = {
		enabled: $('#global-flow-enabled').val() === 'true',
		mode: $('#global-flow-mode').val(),
		backpressure: $('#backpressure-mode').val() === 'enable',
		pause_type: $('#pause-frame-type').val()
	};

	$.ajax({
		url: '/api/v1/config/flow-control/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('全局配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function savePortFlowConfig() {
	var portId = $('#edit-port-id').text();
	var data = {
		enabled: $('#port-flow-enabled').val() === 'true',
		negotiation: $('#port-flow-negotiation').val(),
		pause_direction: $('#port-pause-direction').val()
	};

	$.ajax({
		url: '/api/v1/config/flow-control/' + portId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('端口配置保存成功', {icon: 1});
				$('#flowConfigModal').modal('hide');
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
</script>
`

	boxContent := template.HTML(flowContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-tint"></i> 流量控制`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "流量控制",
		Description: "配置端口流量控制功能，通过 Pause 帧或背压机制防止拥塞导致的数据包丢失",
	}

	return panel, nil
}
