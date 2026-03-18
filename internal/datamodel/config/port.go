package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPortsContent 端口状态页面
func getPortsContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	portsContent := `
<style>
	.config-section { margin-bottom: 30px; }
	.config-section h4 { margin-bottom: 15px; color: #333; }
	.port-status-table { font-size: 13px; }
	.port-status-table .status-up { color: #28a745; font-weight: bold; }
	.port-status-table .status-down { color: #dc3545; }
	.port-status-table .status-disabled { color: #6c757d; }
	.btn-edit { padding: 2px 8px; font-size: 12px; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshPorts()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover port-status-table">
				<thead>
					<tr>
						<th width="50">端口</th>
						<th width="80">状态</th>
						<th width="80">链路</th>
						<th>速率/双工</th>
						<th>流控</th>
						<th>描述</th>
						<th>聚合组</th>
						<th width="100">操作</th>
					</tr>
				</thead>
				<tbody id="ports-table-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-up">启用</span></td>
						<td><span class="status-up">Up</span></td>
						<td>1000M/Full</td>
						<td>Off</td>
						<td>Server-A</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-edit btn-primary" onclick="editPort('GE1/0/1')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-up">启用</span></td>
						<td><span class="status-down">Down</span></td>
						<td>Auto/Auto</td>
						<td>Off</td>
						<td></td>
						<td>-</td>
						<td><button class="btn btn-sm btn-edit btn-primary" onclick="editPort('GE1/0/2')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td><span class="status-disabled">禁用</span></td>
						<td><span class="status-down">Down</span></td>
						<td>Auto/Auto</td>
						<td>Off</td>
						<td></td>
						<td>Ag1</td>
						<td><button class="btn btn-sm btn-edit btn-default" onclick="viewPort('GE1/0/3')" disabled>查看 ⚠️</button></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td><span class="status-up">启用</span></td>
						<td><span class="status-up">Up</span></td>
						<td>100M/Full</td>
						<td>On</td>
						<td>AP-Floor2</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-edit btn-primary" onclick="editPort('GE1/0/4')">编辑</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口编辑弹窗 -->
<div class="modal fade" id="portEditModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">编辑端口配置 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="port-edit-form">
					<div class="box box-default">
						<div class="box-header with-border">
							<h4 class="box-title">基本信息</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">管理状态</label>
								<div class="col-sm-9">
									<select class="form-control" id="port-admin-status">
										<option value="enable">启用</option>
										<option value="disable">禁用</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">描述</label>
								<div class="col-sm-9">
									<input type="text" class="form-control" id="port-description" maxlength="64" placeholder="最多 64 字符">
								</div>
							</div>
						</div>
					</div>

					<div class="box box-success">
						<div class="box-header with-border">
							<h4 class="box-title">物理参数</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">速率/双工</label>
								<div class="col-sm-9">
									<select class="form-control" id="port-speed-duplex">
										<option value="auto">自动协商</option>
										<option value="10H">10M 半双工</option>
										<option value="10F">10M 全双工</option>
										<option value="100H">100M 半双工</option>
										<option value="100F">100M 全双工</option>
										<option value="1000F">1000M 全双工</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">流量控制</label>
								<div class="col-sm-9">
									<select class="form-control" id="port-flow-control">
										<option value="off">关闭</option>
										<option value="on">开启</option>
									</select>
								</div>
							</div>
						</div>
					</div>

					<div class="alert alert-warning">
						<i class="fa fa-exclamation-triangle"></i> 注意事项：
						<ul style="margin: 5px 0 0 20px;">
							<li>自动协商失败时，默认使用最低速率（10M 半双工）</li>
							<li>配置固定速率时，需确保对端设备配置一致，否则可能无法通信</li>
						</ul>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortConfig()">应用</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshPorts() {
	location.reload();
}

function editPort(portId) {
	$('#edit-port-id').text(portId);
	// TODO: 加载端口配置
	$('#portEditModal').modal('show');
}

function viewPort(portId) {
	alert('端口 ' + portId + ' 已加入聚合组，不能单独修改配置。\\n\\n如需修改，请前往：配置 → 链路聚合 → 编辑聚合组');
}

function savePortConfig() {
	var portId = $('#edit-port-id').text();
	var data = {
		admin_status: $('#port-admin-status').val(),
		description: $('#port-description').val(),
		speed_duplex: $('#port-speed-duplex').val(),
		flow_control: $('#port-flow-control').val()
	};

	// TODO: 调用 API 保存配置
	console.log('保存端口配置:', portId, data);

	$.ajax({
		url: '/api/v1/ports/' + portId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('端口配置保存成功', {icon: 1});
				$('#portEditModal').modal('hide');
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

	boxContent := template.HTML(portsContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-ethernet"></i> 端口状态`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口状态",
		Description: "配置交换机端口的基础参数，包括管理状态、速率、双工模式、流控、描述等",
	}

	return panel, nil
}
