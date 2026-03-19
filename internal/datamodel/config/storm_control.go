package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getStormControlContent 风暴控制页面
func getStormControlContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	stormContent := `
<style>
	.storm-section { margin-bottom: 30px; }
	.storm-table { font-size: 13px; }
	.storm-table .status-enabled { color: #28a745; font-weight: bold; }
	.storm-table .status-disabled { color: #dc3545; }
	.form-control-sm { padding: 4px 8px; font-size: 12px; }
	.btn-save { margin-top: 10px; }
</style>

<div class="storm-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">全局风暴控制配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveGlobalStormConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">风暴控制类型</label>
				<div class="col-sm-4">
					<select class="form-control" id="storm-type">
						<option value="broadcast">广播风暴</option>
						<option value="multicast">组播风暴</option>
						<option value="unicast">单播风暴</option>
						<option value="all">全部类型</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">控制模式</label>
				<div class="col-sm-4">
					<select class="form-control" id="storm-mode">
						<option value="rate">速率模式 (pps)</option>
						<option value="percent">百分比模式 (%)</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">阈值</label>
				<div class="col-sm-4">
					<input type="number" class="form-control" id="storm-threshold" min="1" max="1000000" value="10000" placeholder="阈值">
				</div>
				<label class="col-sm-2 control-label">动作</label>
				<div class="col-sm-4">
					<select class="form-control" id="storm-action">
						<option value="block">阻断</option>
						<option value="shutdown">关闭端口</option>
						<option value="trap">发送告警</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="storm-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">端口风暴控制状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshStormStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="batchConfigure()">
					<i class="fa fa-cog"></i> 批量配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover storm-table">
				<thead>
					<tr>
						<th width="40"><input type="checkbox" id="select-all"></th>
						<th width="100">端口</th>
						<th width="100">广播</th>
						<th width="100">组播</th>
						<th width="100">单播</th>
						<th width="100">状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="storm-table-body">
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/1"></td>
						<td>GE1/0/1</td>
						<td>10000 pps</td>
						<td>10000 pps</td>
						<td>10000 pps</td>
						<td><span class="status-disabled">关闭</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortStorm('GE1/0/1')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/2"></td>
						<td>GE1/0/2</td>
						<td>5000 pps</td>
						<td>5000 pps</td>
						<td>5000 pps</td>
						<td><span class="status-enabled">启用</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortStorm('GE1/0/2')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/3"></td>
						<td>GE1/0/3</td>
						<td>10000 pps</td>
						<td>10000 pps</td>
						<td>10000 pps</td>
						<td><span class="status-disabled">关闭</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortStorm('GE1/0/3')">配置</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="port-checkbox" value="GE1/0/4"></td>
						<td>GE1/0/4</td>
						<td>8000 pps</td>
						<td>8000 pps</td>
						<td>8000 pps</td>
						<td><span class="status-enabled">启用</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editPortStorm('GE1/0/4')">配置</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口风暴配置弹窗 -->
<div class="modal fade" id="stormConfigModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">配置风暴控制 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="storm-config-form">
					<div class="box box-default">
						<div class="box-body">
							<div class="form-group">
								<label>启用风暴控制</label>
								<select class="form-control" id="port-storm-enabled">
									<option value="false">关闭</option>
									<option value="true">启用</option>
								</select>
							</div>
							<div class="form-group">
								<label>广播阈值 (pps)</label>
								<input type="number" class="form-control" id="port-broadcast-threshold" min="100" max="1000000" value="10000">
							</div>
							<div class="form-group">
								<label>组播阈值 (pps)</label>
								<input type="number" class="form-control" id="port-multicast-threshold" min="100" max="1000000" value="10000">
							</div>
							<div class="form-group">
								<label>单播阈值 (pps)</label>
								<input type="number" class="form-control" id="port-unicast-threshold" min="100" max="1000000" value="10000">
							</div>
							<div class="form-group">
								<label>超过阈值动作</label>
								<select class="form-control" id="port-storm-action">
									<option value="block">阻断风暴流量</option>
									<option value="shutdown">关闭端口</option>
									<option value="trap">仅发送告警</option>
								</select>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortStormConfig()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.port-checkbox').prop('checked', this.checked);
});

function refreshStormStatus() {
	location.reload();
}

function editPortStorm(portId) {
	$('#edit-port-id').text(portId);
	$('#stormConfigModal').modal('show');
}

function batchConfigure() {
	var selected = [];
	$('.port-checkbox:checked').each(function() {
		selected.push($(this).val());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要配置的端口', {icon: 2});
		return;
	}
	layer.msg('批量配置功能待实现，已选择 ' + selected.length + ' 个端口', {icon: 3});
}

function saveGlobalStormConfig() {
	var data = {
		storm_type: $('#storm-type').val(),
		storm_mode: $('#storm-mode').val(),
		threshold: parseInt($('#storm-threshold').val()),
		action: $('#storm-action').val()
	};

	$.ajax({
		url: '/api/v1/config/storm-control/global',
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

function savePortStormConfig() {
	var portId = $('#edit-port-id').text();
	var data = {
		enabled: $('#port-storm-enabled').val() === 'true',
		broadcast_threshold: parseInt($('#port-broadcast-threshold').val()),
		multicast_threshold: parseInt($('#port-multicast-threshold').val()),
		unicast_threshold: parseInt($('#port-unicast-threshold').val()),
		action: $('#port-storm-action').val()
	};

	$.ajax({
		url: '/api/v1/config/storm-control/' + portId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('端口配置保存成功', {icon: 1});
				$('#stormConfigModal').modal('hide');
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

	boxContent := template.HTML(stormContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-shield"></i> 风暴控制`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "风暴控制",
		Description: "配置端口广播、组播、单播风暴抑制阈值，防止网络风暴导致网络瘫痪",
	}

	return panel, nil
}
