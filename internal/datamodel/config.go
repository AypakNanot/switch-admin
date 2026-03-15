package datamodel

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetPortsContent 端口状态页面
func GetPortsContent(ctx *context.Context) (types.Panel, error) {
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

// GetLinkAggregationContent 链路聚合页面
func GetLinkAggregationContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	laContent := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-normal { color: #28a745; }
	.status-degraded { color: #ffc107; }
	.status-down { color: #dc3545; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">链路聚合</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateModal()">
					<i class="fa fa-plus"></i> 新建
				</button>
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelected()">
					<i class="fa fa-trash"></i> 删除
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="50"><input type="checkbox" id="select-all"></th>
						<th>聚合组</th>
						<th>模式</th>
						<th>成员端口</th>
						<th>负载均衡</th>
						<th>最小活跃</th>
						<th>状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="la-table-body">
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag1</td>
						<td>LACP</td>
						<td>GE1/0/1-4</td>
						<td>src-dst-ip</td>
						<td>2</td>
						<td><span class="status-normal">正常</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag1')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag1')">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag2</td>
						<td>Static</td>
						<td>GE1/0/5-6</td>
						<td>src-dst-mac</td>
						<td>1</td>
						<td><span class="status-normal">正常</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag2')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag2')">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag3</td>
						<td>LACP</td>
						<td>GE1/0/9-12</td>
						<td>src-dst-mac</td>
						<td>2</td>
						<td><span class="status-degraded">降级</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag3')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag3')">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑聚合组弹窗 -->
<div class="modal fade" id="laModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="laModalTitle">新建链路聚合组</h4>
			</div>
			<div class="modal-body">
				<form id="la-form">
					<input type="hidden" id="la-id">

					<div class="box box-default">
						<div class="box-header with-border">
							<h4 class="box-title">基本信息</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">聚合组 ID</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-group-id" min="1" max="128" placeholder="范围：1-128">
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">模式</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-mode">
										<option value="LACP">LACP（动态）</option>
										<option value="Static">Static（静态）</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">描述</label>
								<div class="col-sm-9">
									<input type="text" class="form-control" id="la-description" placeholder="可选">
								</div>
							</div>
						</div>
					</div>

					<div class="box box-success">
						<div class="box-header with-border">
							<h4 class="box-title">负载均衡策略</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">负载均衡</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-load-balance">
										<option value="src-mac">src-mac</option>
										<option value="dst-mac">dst-mac</option>
										<option value="src-dst-mac">src-dst-mac</option>
										<option value="src-ip">src-ip</option>
										<option value="dst-ip">dst-ip</option>
										<option value="src-dst-ip" selected>src-dst-ip</option>
									</select>
								</div>
							</div>
						</div>
					</div>

					<div class="box box-info">
						<div class="box-header with-border">
							<h4 class="box-title">成员端口配置</h4>
						</div>
						<div class="box-body">
							<div class="form-group">
								<label>选择端口：</label>
								<div id="port-selection">
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/1"> GE1/0/1</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/2"> GE1/0/2</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/3"> GE1/0/3</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/4"> GE1/0/4</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/5"> GE1/0/5</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/6"> GE1/0/6</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/7"> GE1/0/7</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/8"> GE1/0/8</label>
								</div>
							</div>
							<div class="alert alert-warning">
								<i class="fa fa-exclamation-triangle"></i> 需确保：
								<ul style="margin: 5px 0 0 20px;">
									<li>所有端口速率/双工配置一致</li>
									<li>所有端口未加入其他聚合组</li>
									<li>所有端口未配置镜像/安全等功能</li>
								</ul>
							</div>
						</div>
					</div>

					<div class="box box-warning">
						<div class="box-header with-border">
							<h4 class="box-title">高级选项（LACP 模式）</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">最小活跃链路</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-min-active" min="1" value="1">
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">LACP 超时</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-lacp-timeout">
										<option value="long">长超时 (90s)</option>
										<option value="short">短超时 (3s)</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">LACP 优先级</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-priority" min="1" max="65535" value="32768">
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveLA()">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.row-checkbox').prop('checked', this.checked);
});

function refreshData() {
	location.reload();
}

function showCreateModal() {
	$('#laModalTitle').text('新建链路聚合组');
	$('#la-form')[0].reset();
	$('#la-id').val('');
	$('#laModal').modal('show');
}

function editLA(laId) {
	$('#laModalTitle').text('编辑链路聚合组 - ' + laId);
	// TODO: 加载聚合组配置
	$('#laModal').modal('show');
}

function deleteLA(laId) {
	layer.confirm('确定要删除聚合组 ' + laId + ' 吗？成员端口将恢复独立管理。', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function deleteSelected() {
	var selected = [];
	$('.row-checkbox:checked').each(function() {
		selected.push($(this).closest('tr').find('td:eq(1)').text());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的聚合组', {icon: 2});
		return;
	}
	layer.confirm('确定要删除选中的 ' + selected.length + ' 个聚合组吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 批量删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function saveLA() {
	var laId = $('#la-id').val();
	var selectedPorts = [];
	$('.port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	var data = {
		group_id: parseInt($('#la-group-id').val()),
		mode: $('#la-mode').val(),
		description: $('#la-description').val(),
		load_balance: $('#la-load-balance').val(),
		member_ports: selectedPorts,
		min_active: parseInt($('#la-min-active').val()),
		lacp_timeout: $('#la-lacp-timeout').val(),
		priority: parseInt($('#la-priority').val())
	};

	// TODO: 调用 API 保存配置
	console.log('保存聚合组配置:', data);

	layer.msg('保存成功', {icon: 1});
	$('#laModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}
</script>
`

	boxContent := template.HTML(laContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-link"></i> 链路聚合`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "链路聚合",
		Description: "将多个物理端口捆绑成一个逻辑端口，提升带宽和可靠性。支持 LACP 动态聚合和静态聚合。",
	}

	return panel, nil
}

// GetStormControlContent 风暴控制页面
func GetStormControlContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">风暴控制</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="applyAll()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered">
				<thead>
					<tr>
						<th>端口</th>
						<th>广播风暴控制 (pps)</th>
						<th>组播风暴控制 (pps)</th>
						<th>未知单播风暴控制 (pps)</th>
						<th>动作</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>GE1/0/1</td>
						<td><input type="number" class="form-control" value="1000" min="0" max="1000000"></td>
						<td><input type="number" class="form-control" value="1000" min="0" max="1000000"></td>
						<td><input type="number" class="form-control" value="1000" min="0" max="1000000"></td>
						<td>
							<select class="form-control">
								<option value="shutdown">Shutdown</option>
								<option value="trap">Trap</option>
								<option value="block">Block</option>
							</select>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<script>
function applyAll() {
	layer.msg('配置应用成功', {icon: 1});
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-wind"></i> 风暴控制`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "风暴控制",
		Description: "配置端口风暴控制，防止广播、组播、未知单播风暴。",
	}

	return panel, nil
}

// GetFlowControlContent 流量控制页面
func GetFlowControlContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-on { color: #28a745; font-weight: bold; }
	.status-off { color: #6c757d; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">流量控制配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>流控 (802.3x)</th>
						<th>入口限速</th>
						<th>出口限速</th>
						<th>QoS 策略</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="flow-control-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-off">关闭</span></td>
						<td>100 Mbps</td>
						<td>不限</td>
						<td>default</td>
						<td><button class="btn btn-sm btn-primary" onclick="editFlowControl('GE1/0/1')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-on">发送/接收</span></td>
						<td>不限</td>
						<td>50 Mbps</td>
						<td>VOICE</td>
						<td><button class="btn btn-sm btn-primary" onclick="editFlowControl('GE1/0/2')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td><span class="status-off">关闭</span></td>
						<td>不限</td>
						<td>不限</td>
						<td>default</td>
						<td><button class="btn btn-sm btn-primary" onclick="editFlowControl('GE1/0/3')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td><span class="status-on">仅接收</span></td>
						<td>500 Mbps</td>
						<td>500 Mbps</td>
						<td>VIDEO</td>
						<td><button class="btn btn-sm btn-primary" onclick="editFlowControl('GE1/0/4')">编辑</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 编辑流量控制弹窗 -->
<div class="modal fade" id="flowControlModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">编辑流量控制配置 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="flow-control-form">
					<div class="box box-default">
						<div class="box-header with-border">
							<h4 class="box-title">IEEE 802.3x 流量控制</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">流控模式</label>
								<div class="col-sm-9">
									<select class="form-control" id="flow-control-mode">
										<option value="off">关闭</option>
										<option value="tx-only">仅发送</option>
										<option value="rx-only">仅接收</option>
										<option value="both">发送和接收</option>
									</select>
								</div>
							</div>
							<div class="alert alert-info">
								<i class="fa fa-info-circle"></i>
								说明：当对端设备发送流量过快时，本端口发送暂停帧请求对端降低发送速率
							</div>
						</div>
					</div>

					<div class="box box-success">
						<div class="box-header with-border">
							<h4 class="box-title">端口限速</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">入口限速</label>
								<div class="col-sm-9">
									<div class="row">
										<div class="col-sm-6">
											<select class="form-control" id="ingress-rate-enable">
												<option value="off">关闭</option>
												<option value="on">启用</option>
											</select>
										</div>
										<div class="col-sm-6">
											<input type="number" class="form-control" id="ingress-rate" placeholder="Mbps" min="1" max="10000">
										</div>
									</div>
									<small class="text-muted">范围：1-10000 Mbps</small>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">入口突发大小</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="ingress-burst" value="1024" min="64" max="16384">
									<small class="text-muted">范围：64-16384 KB，默认：1024 KB</small>
								</div>
							</div>
							<hr>
							<div class="form-group row">
								<label class="col-sm-3 control-label">出口限速</label>
								<div class="col-sm-9">
									<div class="row">
										<div class="col-sm-6">
											<select class="form-control" id="egress-rate-enable">
												<option value="off">关闭</option>
												<option value="on">启用</option>
											</select>
										</div>
										<div class="col-sm-6">
											<input type="number" class="form-control" id="egress-rate" placeholder="Mbps" min="1" max="10000">
										</div>
									</div>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">出口突发大小</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="egress-burst" value="1024" min="64" max="16384">
								</div>
							</div>
						</div>
					</div>

					<div class="box box-info">
						<div class="box-header with-border">
							<h4 class="box-title">QoS 策略</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">应用策略</label>
								<div class="col-sm-9">
									<select class="form-control" id="qos-policy">
										<option value="default">default - 默认策略</option>
										<option value="VOICE">VOICE - 语音优先</option>
										<option value="VIDEO">VIDEO - 视频优先</option>
										<option value="DATA">DATA - 数据业务</option>
									</select>
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveFlowControl()">应用</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshData() {
	location.reload();
}

function editFlowControl(portId) {
	$('#edit-port-id').text(portId);
	// TODO: 加载端口流量控制配置
	$('#flow-control-form')[0].reset();
	$('#flowControlModal').modal('show');
}

function saveFlowControl() {
	var portId = $('#edit-port-id').text();
	var data = {
		port_id: portId,
		flow_control_mode: $('#flow-control-mode').val(),
		ingress_rate_enable: $('#ingress-rate-enable').val() === 'on',
		ingress_rate: parseInt($('#ingress-rate').val()) || 0,
		ingress_burst: parseInt($('#ingress-burst').val()) || 1024,
		egress_rate_enable: $('#egress-rate-enable').val() === 'on',
		egress_rate: parseInt($('#egress-rate').val()) || 0,
		egress_burst: parseInt($('#egress-burst').val()) || 1024,
		qos_policy: $('#qos-policy').val()
	};

	// TODO: 调用 API 保存
	console.log('保存流量控制配置:', data);
	layer.msg('配置应用成功', {icon: 1});
	$('#flowControlModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-tachometer-alt"></i> 流量控制`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "流量控制",
		Description: "配置 IEEE 802.3x 流控、端口限速、QoS 策略",
	}

	return panel, nil
}

// GetPortIsolationContent 端口隔离页面
func GetPortIsolationContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口隔离组</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateGroupModal()">
					<i class="fa fa-plus"></i> 新建
				</button>
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelectedGroups()">
					<i class="fa fa-trash"></i> 删除
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="50"><input type="checkbox" id="select-all"></th>
						<th>隔离组 ID</th>
						<th>组名称</th>
						<th>成员端口</th>
						<th>隔离模式</th>
						<th>状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="group-table-body">
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>1</td>
						<td>Isolation-Group-1</td>
						<td>GE1/0/1-8</td>
						<td>二层隔离</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editGroup(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteGroup(1)">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>2</td>
						<td>Isolation-Group-2</td>
						<td>GE1/0/9-16</td>
						<td>二层隔离</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editGroup(2)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteGroup(2)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">快速隔离配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="quickSetIsolation()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>隔离组</th>
						<th>隔离模式</th>
						<th>上行端口</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="port-isolation-body">
					<tr>
						<td>GE1/0/1</td>
						<td>
							<select class="form-control" style="width:150px;display:inline-block" id="port-group-1">
								<option value="">不属于隔离组</option>
								<option value="1" selected>Isolation-Group-1</option>
								<option value="2">Isolation-Group-2</option>
							</select>
						</td>
						<td>
							<select class="form-control" style="width:120px;display:inline-block" id="port-mode-1">
								<option value="l2">二层隔离</option>
								<option value="all">全部隔离</option>
							</select>
						</td>
						<td>
							<label class="checkbox-inline" style="margin:0">
								<input type="checkbox" class="uplink-check" data-port="1"> 设为上行端口
							</label>
						</td>
						<td><button class="btn btn-sm btn-primary" onclick="setPortIsolation('GE1/0/1')">设置</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td>
							<select class="form-control" style="width:150px;display:inline-block" id="port-group-2">
								<option value="" selected>不属于隔离组</option>
								<option value="1">Isolation-Group-1</option>
								<option value="2">Isolation-Group-2</option>
							</select>
						</td>
						<td>
							<select class="form-control" style="width:120px;display:inline-block" id="port-mode-2">
								<option value="l2">二层隔离</option>
								<option value="all">全部隔离</option>
							</select>
						</td>
						<td>
							<label class="checkbox-inline" style="margin:0">
								<input type="checkbox" class="uplink-check" data-port="2"> 设为上行端口
							</label>
						</td>
						<td><button class="btn btn-sm btn-primary" onclick="setPortIsolation('GE1/0/2')">设置</button></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td>
							<select class="form-control" style="width:150px;display:inline-block" id="port-group-3">
								<option value="" selected>不属于隔离组</option>
								<option value="1">Isolation-Group-1</option>
								<option value="2">Isolation-Group-2</option>
							</select>
						</td>
						<td>
							<select class="form-control" style="width:120px;display:inline-block" id="port-mode-3">
								<option value="l2">二层隔离</option>
								<option value="all">全部隔离</option>
							</select>
						</td>
						<td>
							<label class="checkbox-inline" style="margin:0">
								<input type="checkbox" class="uplink-check" data-port="3"> 设为上行端口
							</label>
						</td>
						<td><button class="btn btn-sm btn-primary" onclick="setPortIsolation('GE1/0/3')">设置</button></td>
					</tr>
				</tbody>
			</table>
			<div class="alert alert-info">
				<i class="fa fa-info-circle"></i>
				<strong>说明：</strong><br>
				• 端口隔离用于实现同一隔离组内端口之间不能互相通信，只能与上行端口通信<br>
				• 二层隔离：仅隔离二层流量，三层流量仍可通信<br>
				• 全部隔离：二层和三层流量都隔离<br>
				• 上行端口：隔离组内端口可以与上行端口通信
			</div>
		</div>
	</div>
</div>

<!-- 新建/编辑隔离组弹窗 -->
<div class="modal fade" id="groupModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="groupModalTitle">新建隔离组</h4>
			</div>
			<div class="modal-body">
				<form id="group-form">
					<input type="hidden" id="group-id-edit">
					<div class="form-group row">
						<label class="col-sm-3 control-label">隔离组 ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="group-id" min="1" max="64" placeholder="范围：1-64">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">组名称</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="group-name" placeholder="如：Isolation-Group-1">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">描述</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="group-description" placeholder="可选">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">成员端口</label>
						<div class="col-sm-9">
							<div id="port-selection">
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/1"> GE1/0/1</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/2"> GE1/0/2</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/3"> GE1/0/3</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/4"> GE1/0/4</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/5"> GE1/0/5</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/6"> GE1/0/6</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/7"> GE1/0/7</label>
								<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/8"> GE1/0/8</label>
							</div>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">上行端口</label>
						<div class="col-sm-9">
							<div id="uplink-selection">
								<label class="checkbox-inline"><input type="checkbox" class="uplink-checkbox" value="GE1/0/23"> GE1/0/23</label>
								<label class="checkbox-inline"><input type="checkbox" class="uplink-checkbox" value="GE1/0/24"> GE1/0/24</label>
							</div>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">隔离模式</label>
						<div class="col-sm-9">
							<select class="form-control" id="group-mode">
								<option value="l2">二层隔离</option>
								<option value="all">全部隔离</option>
							</select>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveGroup()">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.row-checkbox').prop('checked', this.checked);
});

function refreshData() {
	location.reload();
}

function showCreateGroupModal() {
	$('#groupModalTitle').text('新建隔离组');
	$('#group-form')[0].reset();
	$('#group-id-edit').val('');
	$('#group-id').prop('disabled', false);
	$('#groupModal').modal('show');
}

function editGroup(groupId) {
	$('#groupModalTitle').text('编辑隔离组 - ' + groupId);
	// TODO: 加载隔离组配置
	$('#group-id-edit').val(groupId);
	$('#group-id').prop('disabled', true);
	$('#groupModal').modal('show');
}

function saveGroup() {
	var selectedPorts = [];
	$('.port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});
	var uplinkPorts = [];
	$('.uplink-checkbox:checked').each(function() {
		uplinkPorts.push($(this).val());
	});

	var data = {
		group_id: parseInt($('#group-id').val()),
		name: $('#group-name').val(),
		description: $('#group-description').val(),
		member_ports: selectedPorts,
		uplink_ports: uplinkPorts,
		mode: $('#group-mode').val()
	};

	// TODO: 调用 API 保存
	console.log('保存隔离组:', data);
	layer.msg('保存成功', {icon: 1});
	$('#groupModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function deleteGroup(groupId) {
	layer.confirm('确定要删除隔离组 ' + groupId + ' 吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function deleteSelectedGroups() {
	var selected = [];
	$('.row-checkbox:checked').each(function() {
		selected.push($(this).closest('tr').find('td:eq(1)').text());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的隔离组', {icon: 2});
		return;
	}
	layer.confirm('确定要删除选中的 ' + selected.length + ' 个隔离组吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 批量删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function setPortIsolation(portId) {
	var groupId = $('#port-group-' + portId.split('/')[2]).val();
	var mode = $('#port-mode-' + portId.split('/')[2]).val();
	// TODO: 调用 API 设置
	console.log('设置端口隔离:', portId, groupId, mode);
	layer.msg('配置应用成功', {icon: 1});
}

function quickSetIsolation() {
	layer.msg('批量应用功能开发中...', {icon: 0});
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-columns"></i> 端口隔离`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口隔离",
		Description: "配置端口隔离组，实现端口间通信隔离",
	}

	return panel, nil
}

// GetPortMonitorContent 端口监测页面
func GetPortMonitorContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口镜像配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateMirrorModal()">
					<i class="fa fa-plus"></i> 新建镜像组
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>镜像组 ID</th>
						<th>镜像组名称</th>
						<th>镜像方向</th>
						<th>源端口</th>
						<th>目的端口</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="mirror-table-body">
					<tr>
						<td>1</td>
						<td>Mirror-Group-1</td>
						<td>双向</td>
						<td>GE1/0/1, GE1/0/2</td>
						<td>GE1/0/24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirror(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirror(1)">删除</button>
						</td>
					</tr>
					<tr>
						<td>2</td>
						<td>Mirror-Group-2</td>
						<td>入方向</td>
						<td>GE1/0/5</td>
						<td>GE1/0/23</td>
						<td><span class="status-inactive">未激活</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirror(2)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirror(2)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">端口流量统计</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshStats()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>接收包数</th>
						<th>发送包数</th>
						<th>接收字节数</th>
						<th>发送字节数</th>
						<th>错误包数</th>
						<th>广播包数</th>
					</tr>
				</thead>
				<tbody id="stats-body">
					<tr>
						<td>GE1/0/1</td>
						<td>1,234,567</td>
						<td>987,654</td>
						<td>1.2 GB</td>
						<td>856 MB</td>
						<td>0</td>
						<td>12,345</td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td>567,890</td>
						<td>432,109</td>
						<td>456 MB</td>
						<td>321 MB</td>
						<td>2</td>
						<td>5,678</td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
						<td>0</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑镜像组弹窗 -->
<div class="modal fade" id="mirrorModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="mirrorModalTitle">新建镜像组</h4>
			</div>
			<div class="modal-body">
				<form id="mirror-form">
					<input type="hidden" id="mirror-id-edit">
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像组 ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="mirror-id" min="1" max="4" placeholder="范围：1-4">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像组名称</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="mirror-name" placeholder="如：Mirror-Group-1">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像方向</label>
						<div class="col-sm-9">
							<select class="form-control" id="mirror-direction">
								<option value="both">双向（接收 + 发送）</option>
								<option value="ingress">入方向（仅接收）</option>
								<option value="egress">出方向（仅发送）</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">源端口（被镜像端口）</label>
						<div class="col-sm-9">
							<div id="source-port-selection">
								<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/1"> GE1/0/1</label>
								<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/2"> GE1/0/2</label>
								<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/3"> GE1/0/3</label>
								<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/4"> GE1/0/4</label>
							</div>
							<small class="text-muted">可选择多个源端口</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">目的端口（监控端口）</label>
						<div class="col-sm-9">
							<select class="form-control" id="dest-port">
								<option value="">请选择</option>
								<option value="GE1/0/23">GE1/0/23</option>
								<option value="GE1/0/24">GE1/0/24</option>
							</select>
							<small class="text-muted">目的端口不能与源端口重复</small>
						</div>
					</div>
					<div class="alert alert-warning">
						<i class="fa fa-exclamation-triangle"></i>
						<strong>注意事项：</strong><br>
						• 目的端口不能用于正常数据转发<br>
						• 目的端口速率应大于等于所有源端口速率之和<br>
						• 一个端口只能属于一个镜像组
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveMirror()">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshData() {
	location.reload();
}

function refreshStats() {
	location.reload();
}

function showCreateMirrorModal() {
	$('#mirrorModalTitle').text('新建镜像组');
	$('#mirror-form')[0].reset();
	$('#mirror-id-edit').val('');
	$('#mirror-id').prop('disabled', false);
	$('#mirrorModal').modal('show');
}

function editMirror(mirrorId) {
	$('#mirrorModalTitle').text('编辑镜像组 - ' + mirrorId);
	// TODO: 加载镜像组配置
	$('#mirror-id-edit').val(mirrorId);
	$('#mirror-id').prop('disabled', true);
	$('#mirrorModal').modal('show');
}

function saveMirror() {
	var sourcePorts = [];
	$('.source-port-checkbox:checked').each(function() {
		sourcePorts.push($(this).val());
	});

	var data = {
		mirror_id: parseInt($('#mirror-id').val()),
		name: $('#mirror-name').val(),
		direction: $('#mirror-direction').val(),
		source_ports: sourcePorts,
		dest_port: $('#dest-port').val()
	};

	// TODO: 调用 API 保存
	console.log('保存镜像组:', data);
	layer.msg('保存成功', {icon: 1});
	$('#mirrorModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function deleteMirror(mirrorId) {
	layer.confirm('确定要删除镜像组 ' + mirrorId + ' 吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-eye"></i> 端口监测`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口监测",
		Description: "配置端口镜像和流量统计",
	}

	return panel, nil
}

// GetVLANContent VLAN 配置页面
func GetVLANContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	vlanContent := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.port-mode-access { background-color: #d4edda; padding: 2px 6px; border-radius: 3px; }
	.port-mode-trunk { background-color: #fff3cd; padding: 2px 6px; border-radius: 3px; }
	.port-mode-hybrid { background-color: #cce5ff; padding: 2px 6px; border-radius: 3px; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">VLAN 列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateVlanModal()">
					<i class="fa fa-plus"></i> 新建
				</button>
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelectedVlans()">
					<i class="fa fa-trash"></i> 删除
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshVlanData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="50"><input type="checkbox" id="select-all-vlan"></th>
						<th width="80">VLAN ID</th>
						<th>VLAN 名称</th>
						<th>成员端口</th>
						<th>Tagged 端口</th>
						<th>状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="vlan-table-body">
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>1</td>
						<td>default</td>
						<td>GE1/0/1-24</td>
						<td>-</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlan(1)">编辑</button>
							<button class="btn btn-sm btn-default" onclick="viewVlanMembers(1)" disabled>默认</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>10</td>
						<td>Management</td>
						<td>GE1/0/1-4</td>
						<td>GE1/0/23-24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlan(10)">编辑</button>
							<button class="btn btn-sm btn-info" onclick="viewVlanMembers(10)">查看</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>20</td>
						<td>VOICE</td>
						<td>GE1/0/5-12</td>
						<td>GE1/0/23-24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlan(20)">编辑</button>
							<button class="btn btn-sm btn-info" onclick="viewVlanMembers(20)">查看</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>30</td>
						<td>VIDEO</td>
						<td>GE1/0/13-20</td>
						<td>GE1/0/23-24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlan(30)">编辑</button>
							<button class="btn btn-sm btn-info" onclick="viewVlanMembers(30)">查看</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>100</td>
						<td>Guest</td>
						<td>-</td>
						<td>-</td>
						<td><span class="status-inactive">未使用</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlan(100)">编辑</button>
							<button class="btn btn-sm btn-info" onclick="viewVlanMembers(100)">查看</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口 VLAN 配置 -->
<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">端口 VLAN 配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="batchSetPortVlan()">
					<i class="fa fa-save"></i> 批量应用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>VLAN 模式</th>
						<th>PVID</th>
						<th>Allowed VLAN</th>
						<th>Untagged VLAN</th>
						<th>Tagged VLAN</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="port-vlan-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="port-mode-access">Access</span></td>
						<td>10</td>
						<td>10</td>
						<td>10</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortVlan('GE1/0/1')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="port-mode-access">Access</span></td>
						<td>10</td>
						<td>10</td>
						<td>10</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortVlan('GE1/0/2')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/5</td>
						<td><span class="port-mode-trunk">Trunk</span></td>
						<td>1</td>
						<td>1,10,20</td>
						<td>1</td>
						<td>10,20</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortVlan('GE1/0/5')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/23</td>
						<td><span class="port-mode-hybrid">Hybrid</span></td>
						<td>1</td>
						<td>1-4094</td>
						<td>1,10,20,30</td>
						<td>100</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortVlan('GE1/0/23')">编辑</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑 VLAN 弹窗 -->
<div class="modal fade" id="vlanModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="vlanModalTitle">新建 VLAN</h4>
			</div>
			<div class="modal-body">
				<form id="vlan-form">
					<input type="hidden" id="vlan-id-edit">
					<div class="form-group row">
						<label class="col-sm-3 control-label">VLAN ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="vlan-id" min="1" max="4094" placeholder="范围：1-4094">
							<small class="text-muted">VLAN 1 为默认 VLAN，不能删除</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">VLAN 名称</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="vlan-name" maxlength="32" placeholder="最多 32 字符，如：Management、VOICE">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">描述</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="vlan-description" placeholder="可选">
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveVlan()">确定</button>
			</div>
		</div>
	</div>
</div>

<!-- 端口 VLAN 编辑弹窗 -->
<div class="modal fade" id="portVlanModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">端口 VLAN 配置 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="port-vlan-form">
					<div class="alert alert-warning">
						<i class="fa fa-exclamation-triangle"></i>
						<strong>⚠️ 自杀式切断防护：</strong>修改以下配置可能导致管理连接中断，请谨慎操作！
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">VLAN 模式</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-vlan-mode" onchange="onVlanModeChange()">
								<option value="access">Access - 接入端口，仅属于一个 VLAN</option>
								<option value="trunk">Trunk - 干道端口，可属于多个 VLAN</option>
								<option value="hybrid">Hybrid - 混合端口，灵活配置</option>
							</select>
						</div>
					</div>
					<div class="form-group row" id="pvid-group">
						<label class="col-sm-3 control-label">PVID（默认 VLAN ID）</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="port-pvid" min="1" max="4094" value="1">
						</div>
					</div>
					<div class="form-group row" id="allowed-vlan-group">
						<label class="col-sm-3 control-label">Allowed VLAN</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="port-allowed-vlan" placeholder="如：1,10,20 或 1-100">
							<small class="text-muted">支持逗号分隔或范围，如：1,10,20 或 1-100</small>
						</div>
					</div>
					<div class="form-group row" id="untagged-vlan-group">
						<label class="col-sm-3 control-label">Untagged VLAN</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="port-untagged-vlan" placeholder="如：1,10,20">
						</div>
					</div>
					<div class="form-group row" id="tagged-vlan-group">
						<label class="col-sm-3 control-label">Tagged VLAN</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="port-tagged-vlan" placeholder="如：20,30,100">
						</div>
					</div>
					<div class="alert alert-info">
						<strong>模式说明：</strong><br>
						• <strong>Access：</strong> 仅能属于一个 VLAN，发送不带 Tag 的帧，PVID=Allowed VLAN<br>
						• <strong>Trunk：</strong> 可属于多个 VLAN，默认发送带 Tag 的帧，可配置一个 PVID 发送不带 Tag<br>
						• <strong>Hybrid：</strong> 灵活配置，可配置多个 Untagged VLAN 和 Tagged VLAN
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortVlan()">应用</button>
			</div>
		</div>
	</div>
</div>

<!-- 查看 VLAN 成员弹窗 -->
<div class="modal fade" id="vlanMembersModal" tabindex="-1">
	<div class="modal-dialog modal-md">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">VLAN <span id="view-vlan-id"></span> 成员</h4>
			</div>
			<div class="modal-body">
				<div id="vlan-members-content"></div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all-vlan').change(function() {
	$('.row-checkbox').prop('checked', this.checked);
});

function refreshVlanData() {
	location.reload();
}

function showCreateVlanModal() {
	$('#vlanModalTitle').text('新建 VLAN');
	$('#vlan-form')[0].reset();
	$('#vlan-id-edit').val('');
	$('#vlan-id').prop('disabled', false);
	$('#vlanModal').modal('show');
}

function editVlan(vlanId) {
	$('#vlanModalTitle').text('编辑 VLAN - ' + vlanId);
	// TODO: 加载 VLAN 配置
	$('#vlan-id-edit').val(vlanId);
	$('#vlan-id').val(vlanId).prop('disabled', true);
	$('#vlanModal').modal('show');
}

function saveVlan() {
	var vlanId = parseInt($('#vlan-id').val());
	var vlanName = $('#vlan-name').val();

	if (!vlanId || vlanId < 1 || vlanId > 4094) {
		layer.msg('VLAN ID 必须在 1-4094 范围内', {icon: 2});
		return;
	}

	var data = {
		vlan_id: vlanId,
		name: vlanName,
		description: $('#vlan-description').val()
	};

	// TODO: 调用 API 保存
	console.log('保存 VLAN:', data);
	layer.msg('保存成功', {icon: 1});
	$('#vlanModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function deleteSelectedVlans() {
	var selected = [];
	$('.row-checkbox:checked').each(function() {
		var vlanId = $(this).closest('tr').find('td:eq(1)').text();
		if (vlanId == '1') {
			layer.msg('不能删除默认 VLAN 1', {icon: 2});
			return;
		}
		selected.push(vlanId);
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的 VLAN', {icon: 2});
		return;
	}
	layer.confirm('确定要删除选中的 ' + selected.length + ' 个 VLAN 吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 批量删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function viewVlanMembers(vlanId) {
	$('#view-vlan-id').text(vlanId);
	// TODO: 加载 VLAN 成员信息
	$('#vlan-members-content').html('<p>加载中...</p>');
	$('#vlanMembersModal').modal('show');
}

// 端口 VLAN 配置
function editPortVlan(portId) {
	$('#edit-port-id').text(portId);
	// TODO: 加载端口 VLAN 配置
	$('#portVlanModal').modal('show');
}

function onVlanModeChange() {
	var mode = $('#port-vlan-mode').val();
	if (mode === 'access') {
		$('#allowed-vlan-group').hide();
		$('#untagged-vlan-group').hide();
		$('#tagged-vlan-group').hide();
	} else if (mode === 'trunk') {
		$('#allowed-vlan-group').show();
		$('#untagged-vlan-group').hide();
		$('#tagged-vlan-group').show();
	} else { // hybrid
		$('#allowed-vlan-group').show();
		$('#untagged-vlan-group').show();
		$('#tagged-vlan-group').show();
	}
}

function savePortVlan() {
	var portId = $('#edit-port-id').text();
	var mode = $('#port-vlan-mode').val();
	var data = {
		port_id: portId,
		mode: mode,
		pvid: parseInt($('#port-pvid').val()),
		allowed_vlan: $('#port-allowed-vlan').val(),
		untagged_vlan: $('#port-untagged-vlan').val(),
		tagged_vlan: $('#port-tagged-vlan').val()
	};

	// 校验
	if (mode === 'access' && !data.pvid) {
		layer.msg('Access 模式必须配置 PVID', {icon: 2});
		return;
	}
	if (mode === 'trunk' && !data.allowed_vlan) {
		layer.msg('Trunk 模式必须配置 Allowed VLAN', {icon: 2});
		return;
	}

	// 自杀式切断防护检查
	if (portId === 'GE1/0/23') { // 假设这是管理端口
		layer.confirm('⚠️ 警告：修改此端口 VLAN 配置可能导致管理连接中断！\\n\\n确定要继续吗？', {
			btn: ['确定（我已知晓风险）', '取消'],
			icon: 'warning'
		}, function() {
			doSavePortVlan(data);
		});
	} else {
		doSavePortVlan(data);
	}
}

function doSavePortVlan(data) {
	// TODO: 调用 API 保存
	console.log('保存端口 VLAN 配置:', data);
	layer.msg('配置应用成功', {icon: 1});
	$('#portVlanModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function batchSetPortVlan() {
	layer.msg('批量应用功能开发中...', {icon: 0});
}

// 初始化
onVlanModeChange();
</script>
`

	boxContent := template.HTML(vlanContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-network-wired"></i> VLAN 配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "VLAN",
		Description: "配置 VLAN、端口 VLAN 模式（Access/Trunk/Hybrid）和成员关系",
	}

	return panel, nil
}

// GetMacTableContent MAC 地址表页面
func GetMacTableContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
	.type-static { color: #28a745; font-weight: bold; }
	.type-dynamic { color: #17a2b8; }
	.type-blackhole { color: #dc3545; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">MAC 地址表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showAddStaticMacModal()">
					<i class="fa fa-plus"></i> 添加静态 MAC
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-address-book"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">动态 MAC 数量</span>
							<span class="info-box-number">128</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-lock"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">静态 MAC 数量</span>
							<span class="info-box-number">12</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-red"><i class="fa fa-ban"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">黑洞 MAC 数量</span>
							<span class="info-box-number">3</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-chart-bar"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">MAC 表使用率</span>
							<span class="info-box-number">143/8192 (1.7%)</span>
						</div>
					</div>
				</div>
			</div>
			<hr>
			<div class="form-inline">
				<div class="form-group">
					<label>VLAN:</label>
					<select class="form-control" id="filter-vlan" style="width:100px;display:inline-block">
						<option value="all">全部</option>
						<option value="1">1</option>
						<option value="10">10</option>
						<option value="20">20</option>
						<option value="30">30</option>
					</select>
				</div>
				<div class="form-group">
					<label>类型:</label>
					<select class="form-control" id="filter-type" style="width:100px;display:inline-block">
						<option value="all">全部</option>
						<option value="dynamic">动态</option>
						<option value="static">静态</option>
						<option value="blackhole">黑洞</option>
					</select>
				</div>
				<div class="form-group">
					<label>端口:</label>
					<select class="form-control" id="filter-port" style="width:120px;display:inline-block">
						<option value="all">全部</option>
						<option value="GE1/0/1">GE1/0/1</option>
						<option value="GE1/0/2">GE1/0/2</option>
						<option value="GE1/0/3">GE1/0/3</option>
					</select>
				</div>
				<button type="button" class="btn btn-primary" onclick="filterMacTable()">
					<i class="fa fa-search"></i> 筛选
				</button>
			</div>
			<hr>
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>VLAN</th>
						<th>MAC 地址</th>
						<th>类型</th>
						<th>端口/聚合口</th>
						<th>状态</th>
						<th>老化时间</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="mac-table-body">
					<tr>
						<td>10</td>
						<td>00:1A:2B:3C:4D:5E</td>
						<td><span class="type-static">静态</span></td>
						<td>GE1/0/1</td>
						<td>活跃</td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteMac('00:1A:2B:3C:4D:5E')">删除</button>
						</td>
					</tr>
					<tr>
						<td>10</td>
						<td>00:11:22:33:44:55</td>
						<td><span class="type-dynamic">动态</span></td>
						<td>GE1/0/2</td>
						<td>活跃</td>
						<td>180s</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="deleteMac('00:11:22:33:44:55')" disabled>只读</button>
						</td>
					</tr>
					<tr>
						<td>20</td>
						<td>AA:BB:CC:DD:EE:FF</td>
						<td><span class="type-dynamic">动态</span></td>
						<td>GE1/0/5</td>
						<td>活跃</td>
						<td>120s</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="deleteMac('AA:BB:CC:DD:EE:FF')" disabled>只读</button>
						</td>
					</tr>
					<tr>
						<td>30</td>
						<td>11:22:33:44:55:66</td>
						<td><span class="type-blackhole">黑洞</span></td>
						<td>-</td>
						<td>丢弃</td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteMac('11:22:33:44:55:66')">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h4 class="box-title">MAC 地址学习配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveMacConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-3 control-label">MAC 地址学习</label>
				<div class="col-sm-9">
					<select class="form-control" id="mac-learning-enable">
						<option value="on">启用</option>
						<option value="off">关闭</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-3 control-label">MAC 地址学习上限</label>
				<div class="col-sm-9">
					<input type="number" class="form-control" id="mac-max-count" value="0" min="0" max="4096">
					<small class="text-muted">范围：0-4096，0 表示不限制</small>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-3 control-label">MAC 地址老化时间</label>
				<div class="col-sm-9">
					<input type="number" class="form-control" id="mac-age-time" value="300" min="10" max="630">
					<small class="text-muted">范围：10-630 秒，默认：300 秒</small>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- 添加静态 MAC 弹窗 -->
<div class="modal fade" id="staticMacModal" tabindex="-1">
	<div class="modal-dialog modal-md">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">添加静态 MAC 地址</h4>
			</div>
			<div class="modal-body">
				<form id="static-mac-form">
					<div class="form-group">
						<label>VLAN ID</label>
						<input type="number" class="form-control" id="static-mac-vlan" min="1" max="4094" placeholder="1-4094">
					</div>
					<div class="form-group">
						<label>MAC 地址</label>
						<input type="text" class="form-control" id="static-mac-addr" placeholder="格式：00:1A:2B:3C:4D:5E">
						<small class="text-muted">格式：XX:XX:XX:XX:XX:XX (十六进制)</small>
					</div>
					<div class="form-group">
						<label>绑定端口</label>
						<select class="form-control" id="static-mac-port">
							<option value="GE1/0/1">GE1/0/1</option>
							<option value="GE1/0/2">GE1/0/2</option>
							<option value="GE1/0/3">GE1/0/3</option>
							<option value="GE1/0/4">GE1/0/4</option>
							<option value="GE1/0/5">GE1/0/5</option>
							<option value="GE1/0/6">GE1/0/6</option>
						</select>
					</div>
					<div class="alert alert-info">
						<i class="fa fa-info-circle"></i>
						<strong>静态 MAC 说明：</strong><br>
						• 静态 MAC 地址不会老化，设备重启后仍存在<br>
						• 用于绑定特定设备到指定端口，增强安全性
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveStaticMac()">确定</button>
			</div>
		</div>
	</div>
</div>

<!-- 添加黑洞 MAC 弹窗 -->
<div class="modal fade" id="blackholeMacModal" tabindex="-1">
	<div class="modal-dialog modal-md">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">添加黑洞 MAC 地址</h4>
			</div>
			<div class="modal-body">
				<form id="blackhole-mac-form">
					<div class="form-group">
						<label>VLAN ID</label>
						<input type="number" class="form-control" id="blackhole-mac-vlan" min="1" max="4094">
					</div>
					<div class="form-group">
						<label>MAC 地址</label>
						<input type="text" class="form-control" id="blackhole-mac-addr" placeholder="格式：00:1A:2B:3C:4D:5E">
					</div>
					<div class="alert alert-warning">
						<i class="fa fa-exclamation-triangle"></i>
						<strong>黑洞 MAC 说明：</strong><br>
						• 黑洞 MAC 地址会被永久丢弃，不会转发<br>
						• 用于过滤恶意设备或广播风暴
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-danger" onclick="saveBlackholeMac()">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshData() {
	location.reload();
}

function filterMacTable() {
	var vlan = $('#filter-vlan').val();
	var type = $('#filter-type').val();
	var port = $('#filter-port').val();
	console.log('筛选:', vlan, type, port);
	// TODO: 调用 API 筛选
	layer.msg('筛选功能开发中...', {icon: 0});
}

function showAddStaticMacModal() {
	$('#static-mac-form')[0].reset();
	$('#staticMacModal').modal('show');
}

function saveStaticMac() {
	var data = {
		vlan_id: parseInt($('#static-mac-vlan').val()),
		mac_address: $('#static-mac-addr').val(),
		port_id: $('#static-mac-port').val(),
		type: 'static'
	};

	// MAC 地址格式校验
	var macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
	if (!macRegex.test(data.mac_address)) {
		layer.msg('MAC 地址格式不正确，应为：00:1A:2B:3C:4D:5E', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存
	console.log('添加静态 MAC:', data);
	layer.msg('添加成功', {icon: 1});
	$('#staticMacModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function deleteMac(macAddress) {
	layer.confirm('确定要删除 MAC 地址 ' + macAddress + ' 吗？', {
		btn: ['确定', '取消']
	}, function() {
		// TODO: 调用 API 删除
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function saveMacConfig() {
	var data = {
		mac_learning_enable: $('#mac-learning-enable').val() === 'on',
		mac_max_count: parseInt($('#mac-max-count').val()),
		mac_age_time: parseInt($('#mac-age-time').val())
	};

	// TODO: 调用 API 保存
	console.log('保存 MAC 配置:', data);
	layer.msg('配置保存成功', {icon: 1});
	setTimeout(function() { location.reload(); }, 1000);
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-address-book"></i> MAC 地址表`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "MAC 地址表",
		Description: "查看和管理 MAC 地址表，配置静态 MAC、黑洞 MAC",
	}

	return panel, nil
}

// GetSTPContent 生成树配置页面
func GetSTPContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.status-disabled { color: #dc3545; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">生成树协议配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-sitemap"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">STP 状态</span>
							<span class="info-box-number" style="color:#28a745">启用</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-circle"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">协议模式</span>
							<span class="info-box-number">RSTP</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-crown"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">桥优先级</span>
							<span class="info-box-number">4096</span>
						</div>
					</div>
				</div>
			</div>
			<hr>
			<form class="form-inline">
				<div class="form-group">
					<label>STP 状态:</label>
					<select class="form-control" id="stp-enable">
						<option value="on" selected>启用</option>
						<option value="off">关闭</option>
					</select>
				</div>
				<div class="form-group">
					<label>协议模式:</label>
					<select class="form-control" id="stp-mode">
						<option value="stp">STP</option>
						<option value="rstp" selected>RSTP (推荐)</option>
						<option value="mstp">MSTP</option>
					</select>
				</div>
				<div class="form-group">
					<label>桥优先级:</label>
					<select class="form-control" id="bridge-priority">
						<option value="0">0</option>
						<option value="4096" selected>4096</option>
						<option value="8192">8192</option>
						<option value="12288">12288</option>
						<option value="16384">16384</option>
						<option value="20480">20480</option>
						<option value="24576">24576</option>
						<option value="28672">28672</option>
						<option value="32768">32768 (默认)</option>
						<option value="36864">36864</option>
						<option value="40960">40960</option>
						<option value="45056">45056</option>
						<option value="49152">49152</option>
						<option value="53248">53248</option>
						<option value="57344">57344</option>
						<option value="61440">61440</option>
					</select>
				</div>
				<button type="button" class="btn btn-primary" onclick="saveStpConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</form>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">端口 STP 配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="batchSetPortStp()">
					<i class="fa fa-save"></i> 批量应用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>STP 状态</th>
						<th>端口优先级</th>
						<th>路径开销</th>
						<th>边缘端口</th>
						<th>BPDU 保护</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="port-stp-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-active">转发</span></td>
						<td>128</td>
						<td>20000</td>
						<td>否</td>
						<td>否</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortStp('GE1/0/1')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-active">转发</span></td>
						<td>128</td>
						<td>20000</td>
						<td>是</td>
						<td>是</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortStp('GE1/0/2')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td><span class="status-inactive">阻塞</span></td>
						<td>128</td>
						<td>20000</td>
						<td>否</td>
						<td>否</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortStp('GE1/0/3')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td><span class="status-disabled">禁用</span></td>
						<td>128</td>
						<td>20000</td>
						<td>否</td>
						<td>否</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortStp('GE1/0/4')">编辑</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口 STP 编辑弹窗 -->
<div class="modal fade" id="portStpModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">端口 STP 配置 - <span id="edit-port-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="port-stp-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">STP 状态</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-stp-enable">
								<option value="on">启用</option>
								<option value="off">禁用</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">端口优先级</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-priority">
								<option value="0">0</option>
								<option value="16">16</option>
								<option value="32">32</option>
								<option value="48">48</option>
								<option value="64">64</option>
								<option value="80">80</option>
								<option value="96">96</option>
								<option value="112">112</option>
								<option value="128" selected>128 (默认)</option>
								<option value="144">144</option>
								<option value="160">160</option>
								<option value="176">176</option>
								<option value="192">192</option>
								<option value="208">208</option>
								<option value="224">224</option>
								<option value="240">240</option>
							</select>
							<small class="text-muted">值越小优先级越高，必须为 16 的倍数</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">路径开销</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="port-cost" value="20000" min="1" max="200000000">
							<small class="text-muted">范围：1-200000000</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">边缘端口</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-edge">
								<option value="no">否</option>
								<option value="yes">是</option>
							</select>
							<small class="text-muted">边缘端口不参与 STP 计算，直接转发</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">BPDU 保护</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-bpdu-protection">
								<option value="no">否</option>
								<option value="yes">是</option>
							</select>
							<small class="text-muted">防止边缘端口接收 BPDU 报文</small>
						</div>
					</div>
					<div class="alert alert-info">
						<i class="fa fa-info-circle"></i>
						<strong>说明：</strong><br>
						• STP 用于消除二层环路，提供冗余备份<br>
						• RSTP 是 STP 的快速收敛版本<br>
						• MSTP 支持多个生成树实例，实现负载分担
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortStp()">应用</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshData() {
	location.reload();
}

function saveStpConfig() {
	var data = {
		stp_enable: $('#stp-enable').val() === 'on',
		stp_mode: $('#stp-mode').val(),
		bridge_priority: parseInt($('#bridge-priority').val())
	};

	// TODO: 调用 API 保存
	console.log('保存 STP 配置:', data);
	layer.msg('配置保存成功', {icon: 1});
	setTimeout(function() { location.reload(); }, 1000);
}

function editPortStp(portId) {
	$('#edit-port-id').text(portId);
	// TODO: 加载端口 STP 配置
	$('#portStpModal').modal('show');
}

function savePortStp() {
	var portId = $('#edit-port-id').text();
	var data = {
		port_id: portId,
		stp_enable: $('#port-stp-enable').val() === 'on',
		priority: parseInt($('#port-priority').val()),
		cost: parseInt($('#port-cost').val()),
		edge: $('#port-edge').val() === 'yes',
		bpdu_protection: $('#port-bpdu-protection').val() === 'yes'
	};

	// TODO: 调用 API 保存
	console.log('保存端口 STP 配置:', data);
	layer.msg('配置应用成功', {icon: 1});
	$('#portStpModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}

function batchSetPortStp() {
	layer.msg('批量应用功能开发中...', {icon: 0});
}
</script>
`

	boxContent := template.HTML(content)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-sitemap"></i> 生成树`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "生成树",
		Description: "配置 STP/RSTP/MSTP 生成树协议",
	}

	return panel, nil
}

// GetERPSContent ERPS 配置页面
func GetERPSContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	erpsContent := `
<style>
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.status-protection { color: #dc3545; font-weight: bold; }
	.role-rpl { color: #007bff; }
	.role-non-rpl { color: #6c757d; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">ERPS 状态概览</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshErpsStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-circle"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">ERPS 状态</span>
							<span class="info-box-number" style="color:#28a745">Idle</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-share-alt"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">协议模式</span>
							<span class="info-box-number">ERPS (G.8032)</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-exclamation-triangle"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">保护切换次数</span>
							<span class="info-box-number">0</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">ERPS 环配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="createErpsRing()">
					<i class="fa fa-plus"></i> 创建环
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>环 ID</th>
						<th>环名称</th>
						<th>VLAN ID</th>
						<th>RPL 端口</th>
						<th>成员端口</th>
						<th>环状态</th>
						<th>保护状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="erps-ring-body">
					<tr>
						<td>1</td>
						<td>Office-Ring</td>
						<td>100</td>
						<td>GE1/0/24</td>
						<td>GE1/0/1, GE1/0/2, GE1/0/23, GE1/0/24</td>
						<td><span class="status-active">Idle</span></td>
						<td><span class="role-non-rpl">保护</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editErpsRing(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteErpsRing(1)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- ERPS 环创建/编辑弹窗 -->
<div class="modal fade" id="erpsRingModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">ERPS 环配置 - <span id="erps-ring-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="erps-ring-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">环 ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="ring-id" min="1" max="16">
							<small class="text-muted">范围：1-16</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">环名称</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="ring-name" placeholder="例如：Office-Ring">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">控制 VLAN ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="ring-vlan" min="1" max="4094">
							<small class="text-muted">用于 ERPS 控制报文传输的 VLAN</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">RPL 端口</label>
						<div class="col-sm-9">
							<select class="form-control" id="rpl-port">
								<option value="">-- 选择 RPL 端口 --</option>
								<option value="GE1/0/1">GE1/0/1</option>
								<option value="GE1/0/2">GE1/0/2</option>
								<option value="GE1/0/23">GE1/0/23</option>
								<option value="GE1/0/24">GE1/0/24</option>
							</select>
							<small class="text-muted">环保护链路端口，阻塞以防环</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">成员端口</label>
						<div class="col-sm-9">
							<div class="checkbox-list" style="max-height: 200px; overflow-y: auto; border: 1px solid #ddd; padding: 10px;">
								<label class="checkbox-inline"><input type="checkbox" name="member-port" value="GE1/0/1"> GE1/0/1</label>
								<label class="checkbox-inline"><input type="checkbox" name="member-port" value="GE1/0/2"> GE1/0/2</label>
								<label class="checkbox-inline"><input type="checkbox" name="member-port" value="GE1/0/3"> GE1/0/3</label>
								<label class="checkbox-inline"><input type="checkbox" name="member-port" value="GE1/0/23"> GE1/0/23</label>
								<label class="checkbox-inline"><input type="checkbox" name="member-port" value="GE1/0/24"> GE1/0/24</label>
							</div>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">WTR 时间（分钟）</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="wtr-time" value="5" min="1" max="12">
							<small class="text-muted">等待恢复时间，默认 5 分钟</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">保护类型</label>
						<div class="col-sm-9">
							<select class="form-control" id="protection-type">
								<option value="1">1:1 保护</option>
								<option value="1+1">1+1 保护</option>
							</select>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveErpsRing()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshErpsStatus() {
	// TODO: 调用 API 刷新 ERPS 状态
	layer.msg('刷新成功', {icon: 1});
}

function createErpsRing() {
	$('#erps-ring-form')[0].reset();
	$('#erps-ring-id').text('新建');
	$('#erpsRingModal').modal('show');
}

function editErpsRing(ringId) {
	// TODO: 调用 API 获取 ERPS 环配置
	$('#erps-ring-id').text('Ring ' + ringId);
	$('#ring-id').val(ringId);
	$('#ring-name').val('Office-Ring');
	$('#ring-vlan').val(100);
	$('#rpl-port').val('GE1/0/24');
	$('#wtr-time').val(5);
	$('#erpsRingModal').modal('show');
}

function saveErpsRing() {
	var data = {
		ring_id: parseInt($('#ring-id').val()),
		ring_name: $('#ring-name').val(),
		control_vlan: parseInt($('#ring-vlan').val()),
		rpl_port: $('#rpl-port').val(),
		wtr_time: parseInt($('#wtr-time').val()),
		protection_type: $('#protection-type').val(),
		member_ports: []
	};
	$('input[name="member-port"]:checked').each(function() {
		data.member_ports.push($(this).val());
	});

	if (!data.ring_id || !data.ring_name || !data.control_vlan || !data.rpl_port) {
		layer.msg('请填写完整配置信息', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存 ERPS 环配置
	console.log('保存 ERPS 环配置:', data);
	layer.msg('保存成功', {icon: 1});
	$('#erpsRingModal').modal('hide');
}

function deleteErpsRing(ringId) {
	layer.confirm('确定要删除 ERPS 环 ' + ringId + ' 吗？', {
		icon: 3,
		title: '警告'
	}, function(index) {
		// TODO: 调用 API 删除 ERPS 环
		layer.msg('删除成功', {icon: 1});
		layer.close(index);
	});
}
</script>
`

	boxContent := template.HTML(erpsContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-share-alt"></i> ERPS`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "ERPS",
		Description: "配置以太环保护切换协议（ERPS），提供小于 50ms 的环网保护切换。",
	}, nil
}

// GetPoEContent PoE 配置页面
func GetPoEContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	poeContent := `
<style>
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.status-warning { color: #ffc107; font-weight: bold; }
	.status-error { color: #dc3545; font-weight: bold; }
	.poetyp-af { color: #007bff; }
	.poetyp-at { color: #28a745; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">PoE 状态概览</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshPoeStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-bolt"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">PoE 总预算</span>
							<span class="info-box-number">370W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-plug"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">已用功率</span>
							<span class="info-box-number">125W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-exclamation-triangle"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">剩余功率</span>
							<span class="info-box-number">245W</span>
						</div>
					</div>
				</div>
				<div class="col-md-3">
					<div class="info-box">
						<span class="info-box-icon bg-purple"><i class="fa fa-desktop"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">在线设备</span>
							<span class="info-box-number">8/24</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">全局 PoE 配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveGlobalPoeConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
		<div class="box-body">
			<form class="form-inline">
				<div class="form-group">
					<label>PoE 使能:</label>
					<select class="form-control" id="global-poe-enable">
						<option value="on" selected>启用</option>
						<option value="off">关闭</option>
					</select>
				</div>
				<div class="form-group">
					<label>PoE 标准:</label>
					<select class="form-control" id="global-poe-standard">
						<option value="802.3af">802.3af (PoE)</option>
						<option value="802.3at" selected>802.3at (PoE+)</option>
						<option value="802.3bt">802.3bt (PoE++)</option>
					</select>
				</div>
				<div class="form-group">
					<label>最大总功率 (W):</label>
					<input type="number" class="form-control" id="global-poe-max-power" value="370" min="15" max="370">
				</div>
				<div class="form-group">
					<label>管理模式:</label>
					<select class="form-control" id="global-poe-mode">
						<option value="auto">自动分配</option>
						<option value="manual">手动配置</option>
					</select>
				</div>
				<button type="button" class="btn btn-primary" onclick="saveGlobalPoeConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</form>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">端口 PoE 配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="batchSetPoe()">
					<i class="fa fa-save"></i> 批量应用
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>使能</th>
						<th>供电状态</th>
						<th>设备类型</th>
						<th>实际功率 (W)</th>
						<th>最大功率 (W)</th>
						<th>优先级</th>
						<th>电压 (V)</th>
						<th>电流 (mA)</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="poe-port-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-active">启用</span></td>
						<td><span class="status-active">供电中</span></td>
						<td><span class="poetyp-at">802.3at</span></td>
						<td>15.4</td>
						<td>30</td>
						<td>高</td>
						<td>53.2</td>
						<td>290</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortPoe('GE1/0/1')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-active">启用</span></td>
						<td><span class="status-active">供电中</span></td>
						<td><span class="poetyp-af">802.3af</span></td>
						<td>6.8</td>
						<td>15.4</td>
						<td>中</td>
						<td>48.0</td>
						<td>142</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortPoe('GE1/0/2')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td><span class="status-active">启用</span></td>
						<td><span class="status-inactive">未连接</span></td>
						<td>-</td>
						<td>0</td>
						<td>30</td>
						<td>中</td>
						<td>0</td>
						<td>0</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortPoe('GE1/0/3')">编辑</button></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td><span class="status-inactive">禁用</span></td>
						<td><span class="status-inactive">未供电</span></td>
						<td>-</td>
						<td>0</td>
						<td>0</td>
						<td>-</td>
						<td>0</td>
						<td>0</td>
						<td><button class="btn btn-sm btn-primary" onclick="editPortPoe('GE1/0/4')">编辑</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 端口 PoE 编辑弹窗 -->
<div class="modal fade" id="portPoeModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">端口 PoE 配置 - <span id="edit-port-poe-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="port-poe-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">PoE 使能</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-poe-enable">
								<option value="on">启用</option>
								<option value="off">禁用</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">最大供电功率 (W)</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="port-poe-max-power" value="30" min="0" max="90" step="0.1">
							<small class="text-muted">802.3af: 15.4W, 802.3at: 30W, 802.3bt: 90W</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">优先级</label>
						<div class="col-sm-9">
							<select class="form-control" id="port-poe-priority">
								<option value="low">低</option>
								<option value="medium" selected>中</option>
								<option value="high">高</option>
								<option value="critical">关键</option>
							</select>
							<small class="text-muted">功率不足时，高优先级端口优先供电</small>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePortPoe()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
	</div>
</div>

<script>
var currentEditPort = '';

function refreshPoeStatus() {
	// TODO: 调用 API 刷新 PoE 状态
	layer.msg('刷新成功', {icon: 1});
}

function saveGlobalPoeConfig() {
	var data = {
		poe_enable: $('#global-poe-enable').val() === 'on',
		poe_standard: $('#global-poe-standard').val(),
		max_total_power: parseFloat($('#global-poe-max-power').val()),
		management_mode: $('#global-poe-mode').val()
	};

	// TODO: 调用 API 保存全局 PoE 配置
	console.log('保存全局 PoE 配置:', data);
	layer.msg('保存成功', {icon: 1});
}

function editPortPoe(portId) {
	currentEditPort = portId;
	$('#edit-port-poe-id').text(portId);
	// TODO: 调用 API 获取端口 PoE 配置
	$('#port-poe-enable').val('on');
	$('#port-poe-max-power').val(30);
	$('#port-poe-priority').val('medium');
	$('#portPoeModal').modal('show');
}

function savePortPoe() {
	var data = {
		port_id: currentEditPort,
		poe_enable: $('#port-poe-enable').val() === 'on',
		max_power: parseFloat($('#port-poe-max-power').val()),
		priority: $('#port-poe-priority').val()
	};

	if (!data.port_id) {
		layer.msg('端口 ID 不能为空', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存端口 PoE 配置
	console.log('保存端口 PoE 配置:', data);
	layer.msg('保存成功', {icon: 1});
	$('#portPoeModal').modal('hide');
}

function batchSetPoe() {
	var selectedPorts = [];
	$('input[name="port-select"]:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	if (selectedPorts.length === 0) {
		layer.msg('请至少选择一个端口', {icon: 2});
		return;
	}

	var data = {
		ports: selectedPorts,
		poe_enable: $('#batch-poe-enable').val() === 'on',
		max_power: parseFloat($('#batch-poe-max-power').val()),
		priority: $('#batch-poe-priority').val()
	};

	// TODO: 调用 API 批量应用 PoE 配置
	console.log('批量应用 PoE 配置:', data);
	layer.msg('批量应用成功', {icon: 1});
}
</script>
`

	boxContent := template.HTML(poeContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-bolt"></i> PoE`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "PoE",
		Description: "配置 PoE (802.3at/af) 供电功能，管理受电设备（AP、IP 电话、摄像头等）。",
	}, nil
}

// GetPortMirrorContent 端口镜像页面
func GetPortMirrorContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	mirrorContent := `
<style>
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.direction-both { background-color: #e3f2fd; }
	.direction-ingress { background-color: #fff3e0; }
	.direction-egress { background-color: #f3e5f5; }
</style>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">端口镜像配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="createMirrorGroup()">
					<i class="fa fa-plus"></i> 创建镜像组
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>镜像组 ID</th>
						<th>镜像组名称</th>
						<th>镜像方向</th>
						<th>源端口</th>
						<th>目的端口</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="mirror-group-body">
					<tr>
						<td>1</td>
						<td>Monitor-Group-1</td>
						<td><span class="direction-both">双向 (both)</span></td>
						<td>GE1/0/1, GE1/0/2</td>
						<td>GE1/0/24</td>
						<td><span class="status-active">活动</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirrorGroup(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirrorGroup(1)">删除</button>
						</td>
					</tr>
					<tr>
						<td>2</td>
						<td>Ingress-Monitor</td>
						<td><span class="direction-ingress">入方向 (ingress)</span></td>
						<td>GE1/0/3</td>
						<td>GE1/0/23</td>
						<td><span class="status-active">活动</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirrorGroup(2)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirrorGroup(2)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h4 class="box-title"><i class="fa fa-exclamation-triangle"></i> 配置说明</h4>
		</div>
		<div class="box-body">
			<ul>
				<li><strong>源端口：</strong>被监控的端口，流量将被复制到目的端口</li>
				<li><strong>目的端口：</strong>监控端口，连接分析设备（如 Wireshark）</li>
				<li><strong>镜像方向：</strong>
					<ul>
						<li><span class="direction-both">双向 (both)</span> - 监控入站和出站流量</li>
						<li><span class="direction-ingress">入方向 (ingress)</span> - 仅监控入站流量</li>
						<li><span class="direction-egress">出方向 (egress)</span> - 仅监控出站流量</li>
					</ul>
				</li>
				<li><strong>注意事项：</strong>
					<ul>
						<li>目的端口不能用作其他用途，建议专用于流量分析</li>
						<li>每个镜像组只能有一个目的端口</li>
						<li>源端口可以有多个</li>
						<li>镜像会占用交换机 CPU 资源，请谨慎使用</li>
					</ul>
				</li>
			</ul>
		</div>
	</div>
</div>

<!-- 镜像组创建/编辑弹窗 -->
<div class="modal fade" id="mirrorGroupModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">镜像组配置 - <span id="mirror-group-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="mirror-group-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像组 ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="mirror-id" min="1" max="4">
							<small class="text-muted">范围：1-4</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像组名称</label>
						<div class="col-sm-9">
							<input type="text" class="form-control" id="mirror-name" placeholder="例如：Monitor-Group-1">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">镜像方向</label>
						<div class="col-sm-9">
							<select class="form-control" id="mirror-direction">
								<option value="both">双向 (both)</option>
								<option value="ingress">入方向 (ingress)</option>
								<option value="egress">出方向 (egress)</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">源端口（多选）</label>
						<div class="col-sm-9">
							<div class="checkbox-list" style="max-height: 200px; overflow-y: auto; border: 1px solid #ddd; padding: 10px;">
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/1"> GE1/0/1</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/2"> GE1/0/2</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/3"> GE1/0/3</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/4"> GE1/0/4</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/5"> GE1/0/5</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/6"> GE1/0/6</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/23"> GE1/0/23</label>
								<label class="checkbox-inline"><input type="checkbox" name="source-port" value="GE1/0/24"> GE1/0/24</label>
							</div>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">目的端口</label>
						<div class="col-sm-9">
							<select class="form-control" id="dest-port">
								<option value="">-- 选择目的端口 --</option>
								<option value="GE1/0/1">GE1/0/1</option>
								<option value="GE1/0/2">GE1/0/2</option>
								<option value="GE1/0/23">GE1/0/23</option>
								<option value="GE1/0/24">GE1/0/24</option>
							</select>
							<small class="text-muted">连接流量分析设备的端口</small>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveMirrorGroup()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
	</div>
</div>

<script>
function createMirrorGroup() {
	$('#mirror-group-form')[0].reset();
	$('#mirror-group-id').text('新建');
	$('#mirror-id').prop('disabled', false);
	$('#mirrorGroupModal').modal('show');
}

function editMirrorGroup(groupId) {
	// TODO: 调用 API 获取镜像组配置
	$('#mirror-group-id').text('Group ' + groupId);
	$('#mirror-id').val(groupId);
	$('#mirror-id').prop('disabled', true);
	$('#mirror-name').val('Monitor-Group-' + groupId);
	$('#mirror-direction').val('both');
	$('#dest-port').val('GE1/0/24');
	$('#mirrorGroupModal').modal('show');
}

function saveMirrorGroup() {
	var data = {
		mirror_id: parseInt($('#mirror-id').val()),
		mirror_name: $('#mirror-name').val(),
		direction: $('#mirror-direction').val(),
		dest_port: $('#dest-port').val(),
		source_ports: []
	};
	$('input[name="source-port"]:checked').each(function() {
		data.source_ports.push($(this).val());
	});

	if (!data.mirror_id || !data.mirror_name || !data.dest_port) {
		layer.msg('请填写完整配置信息', {icon: 2});
		return;
	}

	if (data.source_ports.length === 0) {
		layer.msg('请至少选择一个源端口', {icon: 2});
		return;
	}

	// 检查源端口和目的端口是否重复
	if (data.source_ports.indexOf(data.dest_port) !== -1) {
		layer.msg('目的端口不能同时作为源端口', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存镜像组配置
	console.log('保存镜像组配置:', data);
	layer.msg('保存成功', {icon: 1});
	$('#mirrorGroupModal').modal('hide');
}

function deleteMirrorGroup(groupId) {
	layer.confirm('确定要删除镜像组 ' + groupId + ' 吗？', {
		icon: 3,
		title: '警告'
	}, function(index) {
		// TODO: 调用 API 删除镜像组
		layer.msg('删除成功', {icon: 1});
		layer.close(index);
	});
}
</script>
`

	boxContent := template.HTML(mirrorContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-eye"></i> 端口镜像`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "端口镜像",
		Description: "配置端口镜像（Port Mirroring），将指定端口的流量复制到监控端口进行分析。",
	}, nil
}

// GetMulticastContent 组播配置页面
func GetMulticastContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	multicastContent := `
<style>
	.status-active { color: #28a745; font-weight: bold; }
	.status-inactive { color: #6c757d; }
	.status-joined { color: #007bff; }
	.mode-snooping { background-color: #e8f5e9; }
	.mode-flooding { background-color: #fff3e0; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">IGMP Snooping 状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshMulticastStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-sitemap"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">IGMP Snooping</span>
							<span class="info-box-number" style="color:#28a745">启用</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-group"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">组播组数量</span>
							<span class="info-box-number">5</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-desktop"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">组成员端口</span>
							<span class="info-box-number">12</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">全局组播配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveGlobalMulticastConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
		<div class="box-body">
			<form class="form-inline">
				<div class="form-group">
					<label>IGMP Snooping:</label>
					<select class="form-control" id="global-igmp-snooping">
						<option value="on" selected>启用</option>
						<option value="off">关闭</option>
					</select>
				</div>
				<div class="form-group">
					<label>IGMP 版本:</label>
					<select class="form-control" id="global-igmp-version">
						<option value="1">IGMP v1</option>
						<option value="2" selected>IGMP v2 (推荐)</option>
						<option value="3">IGMP v3</option>
					</select>
				</div>
				<div class="form-group">
					<label>查询间隔 (秒):</label>
					<input type="number" class="form-control" id="global-query-interval" value="125" min="10" max="600">
				</div>
				<div class="form-group">
					<label>组成员超时 (秒):</label>
					<input type="number" class="form-control" id="global-group-timeout" value="260" min="60" max="600">
				</div>
				<button type="button" class="btn btn-primary" onclick="saveGlobalMulticastConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</form>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">VLAN 组播配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="addVlanMulticast()">
					<i class="fa fa-plus"></i> 添加 VLAN
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>VLAN ID</th>
						<th>IGMP Snooping</th>
						<th>查询器</th>
						<th>组播组</th>
						<th>未知组播处理</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="vlan-multicast-body">
					<tr>
						<td>100</td>
						<td><span class="status-active">启用</span></td>
						<td><span class="mode-snooping">使能</span></td>
						<td>224.0.0.1, 239.1.1.1</td>
						<td><span class="mode-snooping">丢弃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlanMulticast(100)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteVlanMulticast(100)">删除</button>
						</td>
					</tr>
					<tr>
						<td>200</td>
						<td><span class="status-active">启用</span></td>
						<td><span class="mode-flooding">禁用</span></td>
						<td>224.0.0.1</td>
						<td><span class="mode-flooding">泛洪</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVlanMulticast(200)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteVlanMulticast(200)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">组播组列表</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshGroupTable()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>组播组地址</th>
						<th>VLAN ID</th>
						<th>成员端口</th>
						<th>类型</th>
						<th>超时时间</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="multicast-group-body">
					<tr>
						<td>224.0.0.1</td>
						<td>100</td>
						<td>GE1/0/1, GE1/0/2, GE1/0/3</td>
						<td>永久</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-warning" onclick="clearMulticastGroup('224.0.0.1', 100)">清除</button></td>
					</tr>
					<tr>
						<td>239.1.1.1</td>
						<td>100</td>
						<td>GE1/0/1, GE1/0/5</td>
						<td>动态</td>
						<td>260s</td>
						<td><button class="btn btn-sm btn-warning" onclick="clearMulticastGroup('239.1.1.1', 100)">清除</button></td>
					</tr>
					<tr>
						<td>239.255.255.250</td>
						<td>200</td>
						<td>GE1/0/2, GE1/0/4, GE1/0/6</td>
						<td>动态</td>
						<td>180s</td>
						<td><button class="btn btn-sm btn-warning" onclick="clearMulticastGroup('239.255.255.250', 200)">清除</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- VLAN 组播配置弹窗 -->
<div class="modal fade" id="vlanMulticastModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">VLAN 组播配置 - VLAN <span id="vlan-multicast-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="vlan-multicast-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">VLAN ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="vlan-id" min="1" max="4094">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">IGMP Snooping</label>
						<div class="col-sm-9">
							<select class="form-control" id="vlan-igmp-snooping">
								<option value="on">启用</option>
								<option value="off">关闭</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">IGMP 查询器</label>
						<div class="col-sm-9">
							<select class="form-control" id="vlan-querier">
								<option value="on">启用</option>
								<option value="off">关闭</option>
							</select>
							<small class="text-muted">启用后，交换机将发送 IGMP 查询报文</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">查询间隔 (秒)</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="vlan-query-interval" value="125" min="10" max="600">
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">未知组播处理</label>
						<div class="col-sm-9">
							<select class="form-control" id="vlan-unknown-multicast">
								<option value="drop">丢弃</option>
								<option value="flood">泛洪</option>
							</select>
							<small class="text-muted">未学习到的组播流量处理方式</small>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveVlanMulticast()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
	</div>
</div>

<script>
var currentEditVlan = null;

function refreshMulticastStatus() {
	// TODO: 调用 API 刷新组播状态
	layer.msg('刷新成功', {icon: 1});
}

function refreshGroupTable() {
	// TODO: 调用 API 刷新组播组表
	layer.msg('刷新成功', {icon: 1});
}

function saveGlobalMulticastConfig() {
	var data = {
		igmp_snooping_enable: $('#global-igmp-snooping').val() === 'on',
		igmp_version: parseInt($('#global-igmp-version').val()),
		query_interval: parseInt($('#global-query-interval').val()),
		group_timeout: parseInt($('#global-group-timeout').val())
	};

	// TODO: 调用 API 保存全局组播配置
	console.log('保存全局组播配置:', data);
	layer.msg('保存成功', {icon: 1});
}

function addVlanMulticast() {
	$('#vlan-multicast-form')[0].reset();
	$('#vlan-multicast-id').text('新建');
	$('#vlan-id').prop('disabled', false);
	$('#vlanMulticastModal').modal('show');
}

function editVlanMulticast(vlanId) {
	currentEditVlan = vlanId;
	$('#vlan-multicast-id').text(vlanId);
	$('#vlan-id').val(vlanId);
	$('#vlan-id').prop('disabled', true);
	// TODO: 调用 API 获取 VLAN 组播配置
	$('#vlan-igmp-snooping').val('on');
	$('#vlan-querier').val('on');
	$('#vlan-query-interval').val(125);
	$('#vlan-unknown-multicast').val('drop');
	$('#vlanMulticastModal').modal('show');
}

function saveVlanMulticast() {
	var data = {
		vlan_id: parseInt($('#vlan-id').val()),
		igmp_snooping_enable: $('#vlan-igmp-snooping').val() === 'on',
		querier_enable: $('#vlan-querier').val() === 'on',
		query_interval: parseInt($('#vlan-query-interval').val()),
		unknown_multicast_action: $('#vlan-unknown-multicast').val()
	};

	if (!data.vlan_id) {
		layer.msg('VLAN ID 不能为空', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存 VLAN 组播配置
	console.log('保存 VLAN 组播配置:', data);
	layer.msg('保存成功', {icon: 1});
	$('#vlanMulticastModal').modal('hide');
}

function deleteVlanMulticast(vlanId) {
	layer.confirm('确定要删除 VLAN ' + vlanId + ' 的组播配置吗？', {
		icon: 3,
		title: '警告'
	}, function(index) {
		// TODO: 调用 API 删除 VLAN 组播配置
		layer.msg('删除成功', {icon: 1});
		layer.close(index);
	});
}

function clearMulticastGroup(groupAddress, vlanId) {
	layer.confirm('确定要清除组播组 ' + groupAddress + ' (VLAN ' + vlanId + ') 吗？', {
		icon: 3,
		title: '警告'
	}, function(index) {
		// TODO: 调用 API 清除组播组
		layer.msg('清除成功', {icon: 1});
		layer.close(index);
	});
}
</script>
`

	boxContent := template.HTML(multicastContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-group"></i> 组播`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "组播",
		Description: "配置 IGMP Snooping 等组播功能，优化组播流量转发，减少带宽浪费。",
	}, nil
}

// GetResourceContent 资源配置页面
func GetResourceContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	resourceContent := `
<style>
	.usage-low { color: #28a745; }
	usage-medium { color: #ffc107; font-weight: bold; }
	.usage-high { color: #dc3545; font-weight: bold; }
	.progress-bar-low { background-color: #28a745; }
	.progress-bar-medium { background-color: #ffc107; }
	.progress-bar-high { background-color: #dc3545; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">资源使用概览</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshResourceStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-6">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-microchip"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">CPU 使用率</span>
							<span class="info-box-number">23%</span>
							<div class="progress">
								<div class="progress-bar progress-bar-low" style="width: 23%"></div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-6">
					<div class="info-box">
						<span class="info-box-icon bg-purple"><i class="fa fa-memory"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">内存使用率</span>
							<span class="info-box-number">45%</span>
							<div class="progress">
								<div class="progress-bar progress-bar-medium" style="width: 45%"></div>
							</div>
						</div>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-6">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-hdd"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">Flash 使用率</span>
							<span class="info-box-number">128MB / 512MB (25%)</span>
							<div class="progress">
								<div class="progress-bar progress-bar-low" style="width: 25%"></div>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-6">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-tachometer"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">带宽利用率</span>
							<span class="info-box-number">1.2Gbps / 10Gbps (12%)</span>
							<div class="progress">
								<div class="progress-bar progress-bar-low" style="width: 12%"></div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">系统资源详情</h4>
		</div>
		<div class="box-body">
			<table class="table table-bordered">
				<tbody>
					<tr>
						<th width="200">设备型号</th>
						<td>S5750-24X4S</td>
						<th width="200">序列号</th>
						<td>SW20240101001</td>
					</tr>
					<tr>
						<th>硬件版本</th>
						<td>V2.0</td>
						<th>软件版本</th>
						<td>V3.2.1 Build 20240101</td>
					</tr>
					<tr>
						<th>启动时间</th>
						<td>2024-01-15 08:30:00</td>
						<th>运行时间</th>
						<td>15 天 12 小时 30 分钟</td>
					</tr>
					<tr>
						<th>CPU 型号</th>
						<td>Broadcom BCM56340</td>
						<th>CPU 核心数</th>
						<td>4 核 @ 1.2GHz</td>
					</tr>
					<tr>
						<th>总内存</th>
						<td>2GB DDR4</td>
						<th>可用内存</th>
						<td>1.1GB</td>
					</tr>
					<tr>
						<th>总 Flash</th>
						<td>512MB</td>
						<th>可用 Flash</th>
						<td>384MB</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">端口带宽使用</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshPortBandwidth()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>端口</th>
						<th>入站带宽</th>
						<th>出站带宽</th>
						<th>总带宽</th>
						<th>端口速率</th>
						<th>利用率</th>
						<th>状态</th>
					</tr>
				</thead>
				<tbody id="port-bandwidth-body">
					<tr>
						<td>GE1/0/1</td>
						<td>45.2 Mbps</td>
						<td>123.8 Mbps</td>
						<td>169.0 Mbps</td>
						<td>1000 Mbps</td>
						<td>
							<div class="progress" style="margin-bottom: 0;">
								<div class="progress-bar progress-bar-low" style="width: 17%">17%</div>
							</div>
						</td>
						<td><span class="status-active">正常</span></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td>234.5 Mbps</td>
						<td>456.7 Mbps</td>
						<td>691.2 Mbps</td>
						<td>1000 Mbps</td>
						<td>
							<div class="progress" style="margin-bottom: 0;">
								<div class="progress-bar progress-bar-medium" style="width: 69%">69%</div>
							</div>
						</td>
						<td><span class="status-active">正常</span></td>
					</tr>
					<tr>
						<td>GE1/0/3</td>
						<td>12.3 Mbps</td>
						<td>8.9 Mbps</td>
						<td>21.2 Mbps</td>
						<td>1000 Mbps</td>
						<td>
							<div class="progress" style="margin-bottom: 0;">
								<div class="progress-bar progress-bar-low" style="width: 2%">2%</div>
							</div>
						</td>
						<td><span class="status-active">正常</span></td>
					</tr>
					<tr>
						<td>GE1/0/4</td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td>1000 Mbps</td>
						<td>
							<div class="progress" style="margin-bottom: 0;">
								<div class="progress-bar" style="width: 0%"></div>
							</div>
						</td>
						<td><span class="status-inactive">Down</span></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h4 class="box-title"><i class="fa fa-bell"></i> 资源告警阈值</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveAlarmThreshold()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
		<div class="box-body">
			<form class="form-inline">
				<div class="form-group">
					<label>CPU 告警阈值 (%):</label>
					<input type="number" class="form-control" id="cpu-alarm-threshold" value="80" min="50" max="95">
				</div>
				<div class="form-group">
					<label>内存告警阈值 (%):</label>
					<input type="number" class="form-control" id="memory-alarm-threshold" value="85" min="50" max="95">
				</div>
				<div class="form-group">
					<label>带宽告警阈值 (%):</label>
					<input type="number" class="form-control" id="bandwidth-alarm-threshold" value="90" min="50" max="100">
				</div>
				<button type="button" class="btn btn-primary" onclick="saveAlarmThreshold()">
					<i class="fa fa-save"></i> 保存
				</button>
			</form>
		</div>
	</div>
</div>

<script>
function refreshResourceStatus() {
	// TODO: 调用 API 刷新资源状态
	layer.msg('刷新成功', {icon: 1});
}

function refreshPortBandwidth() {
	// TODO: 调用 API 刷新端口带宽
	layer.msg('刷新成功', {icon: 1});
}

function saveAlarmThreshold() {
	var data = {
		cpu_threshold: parseInt($('#cpu-alarm-threshold').val()),
		memory_threshold: parseInt($('#memory-alarm-threshold').val()),
		bandwidth_threshold: parseInt($('#bandwidth-alarm-threshold').val())
	};

	// TODO: 调用 API 保存告警阈值
	console.log('保存告警阈值:', data);
	layer.msg('保存成功', {icon: 1});
}
</script>
`

	boxContent := template.HTML(resourceContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-microchip"></i> 资源`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "资源",
		Description: "查看和管理交换机资源使用情况（CPU、内存、Flash、带宽）。",
	}, nil
}

// GetStackContent 堆叠配置页面
func GetStackContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	stackContent := `
<style>
	.status-master { color: #28a745; font-weight: bold; }
	.status-slave { color: #007bff; }
	.status-standby { color: #ffc107; font-weight: bold; }
	.status-offline { color: #dc3545; }
	.role-master { background-color: #d4edda; }
	.role-slave { background-color: #e3f2fd; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">堆叠状态概览</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshStackStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="row">
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-green"><i class="fa fa-server"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">堆叠状态</span>
							<span class="info-box-number" style="color:#28a745">正常</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-blue"><i class="fa fa-sitemap"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">堆叠成员</span>
							<span class="info-box-number">2/4</span>
						</div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="info-box">
						<span class="info-box-icon bg-yellow"><i class="fa fa-link"></i></span>
						<div class="info-box-content">
							<span class="info-box-text">堆叠带宽</span>
							<span class="info-box-number">20Gbps</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">堆叠成员列表</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="addStackMember()">
					<i class="fa fa-plus"></i> 添加成员
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th>成员 ID</th>
						<th>设备型号</th>
						<th>序列号</th>
						<th>角色</th>
						<th>优先级</th>
						<th>堆叠端口</th>
						<th>软件版本</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="stack-member-body">
					<tr class="role-master">
						<td>1</td>
						<td>S5750-24X4S</td>
						<td>SW20240101001</td>
						<td><span class="status-master">Master</span></td>
						<td>255</td>
						<td>Stack1/1, Stack1/2</td>
						<td>V3.2.1</td>
						<td><span class="status-active">在线</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editStackMember(1)">编辑</button>
							<button class="btn btn-sm btn-warning" onclick="changeStackMaster(1)">设为主设备</button>
						</td>
					</tr>
					<tr class="role-slave">
						<td>2</td>
						<td>S5750-24X4S</td>
						<td>SW20240101002</td>
						<td><span class="status-slave">Slave</span></td>
						<td>100</td>
						<td>Stack1/1, Stack1/2</td>
						<td>V3.2.1</td>
						<td><span class="status-active">在线</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editStackMember(2)">编辑</button>
							<button class="btn btn-sm btn-warning" onclick="changeStackMaster(2)">设为主设备</button>
						</td>
					</tr>
					<tr>
						<td>3</td>
						<td>-</td>
						<td>-</td>
						<td><span class="status-offline">空闲</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td><span class="status-offline">离线</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editStackMember(3)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="removeStackMember(3)">移除</button>
						</td>
					</tr>
					<tr>
						<td>4</td>
						<td>-</td>
						<td>-</td>
						<td><span class="status-offline">空闲</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td><span class="status-offline">离线</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editStackMember(4)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="removeStackMember(4)">移除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h4 class="box-title">全局堆叠配置</h4>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveStackConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
		<div class="box-body">
			<form class="form-inline">
				<div class="form-group">
					<label>堆叠使能:</label>
					<select class="form-control" id="stack-enable">
						<option value="on" selected>启用</option>
						<option value="off">关闭</option>
					</select>
				</div>
				<div class="form-group">
					<label>堆叠模式:</label>
					<select class="form-control" id="stack-mode">
						<option value="auto">自动堆叠</option>
						<option value="manual" selected>手动配置</option>
					</select>
				</div>
				<div class="form-group">
					<label>主设备选举模式:</label>
					<select class="form-control" id="master-election">
						<option value="auto" selected>自动选举</option>
						<option value="priority">优先级优先</option>
						<option value="mac">MAC 地址优先</option>
					</select>
				</div>
				<button type="button" class="btn btn-primary" onclick="saveStackConfig()">
					<i class="fa fa-save"></i> 保存
				</button>
			</form>
		</div>
	</div>
</div>

<div class="config-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h4 class="box-title"><i class="fa fa-exclamation-triangle"></i> 堆叠说明</h4>
		</div>
		<div class="box-body">
			<ul>
				<li><strong>Master 设备：</strong>堆叠系统的主控设备，负责管理和控制整个堆叠</li>
				<li><strong>Slave 设备：</strong>堆叠系统的从设备，受 Master 设备管理</li>
				<li><strong>优先级：</strong>范围 1-255，值越大越容易成为 Master，默认 100</li>
				<li><strong>注意事项：</strong>
					<ul>
						<li>堆叠成员必须是相同型号的设备</li>
						<li>建议使用相同的软件版本</li>
						<li>修改优先级后需要重启设备才能生效</li>
						<li>主设备切换会导致短暂的业务中断</li>
					</ul>
				</li>
			</ul>
		</div>
	</div>
</div>

<!-- 堆叠成员编辑弹窗 -->
<div class="modal fade" id="stackMemberModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">堆叠成员配置 - 成员 <span id="stack-member-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="stack-member-form">
					<div class="form-group row">
						<label class="col-sm-3 control-label">成员 ID</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="member-id" min="1" max="4" disabled>
							<small class="text-muted">成员 ID 不可修改</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">优先级</label>
						<div class="col-sm-9">
							<input type="number" class="form-control" id="member-priority" value="100" min="1" max="255">
							<small class="text-muted">范围：1-255，值越大越容易成为 Master</small>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">堆叠端口 1</label>
						<div class="col-sm-9">
							<select class="form-control" id="stack-port-1">
								<option value="">-- 禁用 --</option>
								<option value="Stack1/1">Stack1/1</option>
								<option value="Stack1/2">Stack1/2</option>
							</select>
						</div>
					</div>
					<div class="form-group row">
						<label class="col-sm-3 control-label">堆叠端口 2</label>
						<div class="col-sm-9">
							<select class="form-control" id="stack-port-2">
								<option value="">-- 禁用 --</option>
								<option value="Stack1/1">Stack1/1</option>
								<option value="Stack1/2">Stack1/2</option>
							</select>
							<small class="text-muted">建议使用两个端口以实现链路冗余</small>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveStackMember()">
					<i class="fa fa-save"></i> 保存
				</button>
			</div>
		</div>
	</div>
</div>

<script>
var currentEditMember = null;

function refreshStackStatus() {
	// TODO: 调用 API 刷新堆叠状态
	layer.msg('刷新成功', {icon: 1});
}

function addStackMember() {
	layer.msg('请通过物理连接添加堆叠成员', {icon: 0});
}

function editStackMember(memberId) {
	currentEditMember = memberId;
	$('#stack-member-id').text(memberId);
	$('#member-id').val(memberId);
	// TODO: 调用 API 获取堆叠成员配置
	$('#member-priority').val(memberId === 1 ? 255 : 100);
	$('#stack-port-1').val('Stack1/1');
	$('#stack-port-2').val('Stack1/2');
	$('#stackMemberModal').modal('show');
}

function saveStackMember() {
	var data = {
		member_id: parseInt($('#member-id').val()),
		priority: parseInt($('#member-priority').val()),
		stack_port_1: $('#stack-port-1').val(),
		stack_port_2: $('#stack-port-2').val()
	};

	if (!data.member_id) {
		layer.msg('成员 ID 不能为空', {icon: 2});
		return;
	}

	// TODO: 调用 API 保存堆叠成员配置
	console.log('保存堆叠成员配置:', data);
	layer.msg('保存成功，重启后生效', {icon: 1});
	$('#stackMemberModal').modal('hide');
}

function changeStackMaster(memberId) {
	layer.confirm('确定要将成员 ' + memberId + ' 设为主设备吗？<br><strong style="color: #dc3545;">警告：此操作会导致短暂的业务中断！</strong>', {
		icon: 3,
		title: '警告',
		showCancelButton: true,
		confirmButtonText: '确认切换',
		cancelButtonText: '取消'
	}, function(index) {
		// TODO: 调用 API 切换主设备
		layer.msg('主设备切换成功', {icon: 1});
		layer.close(index);
	});
}

function removeStackMember(memberId) {
	layer.confirm('确定要移除成员 ' + memberId + ' 吗？', {
		icon: 3,
		title: '警告'
	}, function(index) {
		// TODO: 调用 API 移除堆叠成员
		layer.msg('移除成功', {icon: 1});
		layer.close(index);
	});
}

function saveStackConfig() {
	var data = {
		stack_enable: $('#stack-enable').val() === 'on',
		stack_mode: $('#stack-mode').val(),
		master_election: $('#master-election').val()
	};

	// TODO: 调用 API 保存全局堆叠配置
	console.log('保存全局堆叠配置:', data);
	layer.msg('保存成功', {icon: 1});
}
</script>
`

	boxContent := template.HTML(stackContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-server"></i> 堆叠`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "堆叠",
		Description: "配置交换机堆叠功能，将多台交换机虚拟化为单台设备，简化管理。",
	}, nil
}
