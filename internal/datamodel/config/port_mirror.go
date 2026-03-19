package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPortMirrorContent 端口镜像页面
func getPortMirrorContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	mirrorContent := `
<style>
	.mirror-section { margin-bottom: 30px; }
	.mirror-table { font-size: 13px; }
	.mirror-table .status-active { color: #28a745; font-weight: bold; }
	.mirror-table .status-inactive { color: #dc3545; }
	.direction-tag { display: inline-block; padding: 2px 8px; border-radius: 3px; font-size: 11px; margin-right: 5px; }
	.direction-ingress { background: #5bc0de; color: white; }
	.direction-egress { background: #5cb85c; color: white; }
	.direction-both { background: #f0ad4e; color: white; }
</style>

<div class="mirror-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口镜像配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateMirrorModal()">
					<i class="fa fa-plus"></i> 新建镜像
				</button>
				<button type="button" class="btn btn-default btn-sm" onclick="refreshMirror()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">观察端口 (监控口)</label>
				<div class="col-sm-3">
					<select class="form-control" id="observer-port">
						<option value="">- 请选择 -</option>
						<option value="GE1/0/23">GE1/0/23</option>
						<option value="GE1/0/24">GE1/0/24</option>
					</select>
					<small class="text-muted">用于连接监控设备的端口</small>
				</div>
				<label class="col-sm-2 control-label">观察 VLAN</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="observer-vlan" min="1" max="4094" placeholder="可选">
					<small class="text-muted">用于远程镜像的专用 VLAN</small>
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-primary btn-sm" onclick="saveObserverConfig()">
						<i class="fa fa-save"></i> 保存
					</button>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="mirror-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">镜像会话列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelectedMirrors()">
					<i class="fa fa-trash"></i> 批量删除
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover mirror-table">
				<thead>
					<tr>
						<th width="40"><input type="checkbox" id="select-all"></th>
						<th width="80">会话 ID</th>
						<th width="150">名称</th>
						<th>源端口</th>
						<th>目的端口</th>
						<th width="120">镜像方向</th>
						<th width="80">状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="mirror-table-body">
					<tr>
						<td><input type="checkbox" class="mirror-checkbox" value="1"></td>
						<td>1</td>
						<td>Monitor-AP</td>
						<td>
							<span class="direction-tag direction-both">双向</span>
							GE1/0/1, GE1/0/2
						</td>
						<td>GE1/0/24</td>
						<td>
							<span class="direction-tag direction-both">双向</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirror(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirror(1)">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="mirror-checkbox" value="2"></td>
						<td>2</td>
						<td>Monitor-Server</td>
						<td>
							<span class="direction-tag direction-ingress">入向</span>
							GE1/0/5
						</td>
						<td>GE1/0/24</td>
						<td>
							<span class="direction-tag direction-egress">出向</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirror(2)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirror(2)">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="mirror-checkbox" value="3"></td>
						<td>3</td>
						<td>Test-Mirror</td>
						<td>
							<span class="direction-tag direction-both">双向</span>
							GE1/0/10
						</td>
						<td>GE1/0/23</td>
						<td>
							<span class="direction-tag direction-both">双向</span>
						</td>
						<td><span class="status-inactive">禁用</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editMirror(3)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteMirror(3)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="mirror-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">流镜像配置 (基于 ACL)</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showFlowMirrorModal()">
					<i class="fa fa-plus"></i> 新建流镜像
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="80">会话 ID</th>
						<th width="150">名称</th>
						<th>匹配规则</th>
						<th>源端口</th>
						<th>目的端口</th>
						<th width="100">状态</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="flow-mirror-body">
					<tr>
						<td>10</td>
						<td>HTTP-Monitor</td>
						<td>ACL 100 (TCP 80)</td>
						<td>GE1/0/1</td>
						<td>GE1/0/24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editFlowMirror(10)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteFlowMirror(10)">删除</button>
						</td>
					</tr>
					<tr>
						<td>11</td>
						<td>Voice-Monitor</td>
						<td>ACL 101 (UDP 5060)</td>
						<td>GE1/0/2</td>
						<td>GE1/0/24</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editFlowMirror(11)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteFlowMirror(11)">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑镜像会话弹窗 -->
<div class="modal fade" id="mirrorModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="mirror-modal-title">新建镜像会话</h4>
			</div>
			<div class="modal-body">
				<form id="mirror-form">
					<input type="hidden" id="mirror-id">
					<div class="box box-default">
						<div class="box-body">
							<div class="form-group">
								<label>会话名称</label>
								<input type="text" class="form-control" id="mirror-name" maxlength="32" placeholder="最多 32 字符">
							</div>
							<div class="form-group">
								<label>源端口 (被监控端口)</label>
								<div id="source-port-selection">
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/1"> GE1/0/1</label>
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/2"> GE1/0/2</label>
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/3"> GE1/0/3</label>
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/4"> GE1/0/4</label>
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/5"> GE1/0/5</label>
									<label class="checkbox-inline"><input type="checkbox" class="source-port-checkbox" value="GE1/0/6"> GE1/0/6</label>
								</div>
							</div>
							<div class="form-group">
								<label>目的端口 (观察端口)</label>
								<select class="form-control" id="dest-port">
									<option value="GE1/0/23">GE1/0/23</option>
									<option value="GE1/0/24">GE1/0/24</option>
								</select>
								<small class="text-muted">需确保目的端口速率不低于源端口总和</small>
							</div>
							<div class="form-group">
								<label>镜像方向</label>
								<select class="form-control" id="mirror-direction">
									<option value="both">双向 (ingress + egress)</option>
									<option value="ingress">入方向 (ingress)</option>
									<option value="egress">出方向 (egress)</option>
								</select>
							</div>
							<div class="form-group">
								<label>状态</label>
								<select class="form-control" id="mirror-status">
									<option value="active">启用</option>
									<option value="inactive">禁用</option>
								</select>
							</div>
							<div class="alert alert-warning">
								<i class="fa fa-exclamation-triangle"></i> 注意事项:
								<ul style="margin: 5px 0 0 20px;">
									<li>目的端口不能同时作为源端口</li>
									<li>一个端口可以被多个镜像会话作为源端口</li>
									<li>目的端口将不能用于正常数据转发</li>
								</ul>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveMirror()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.mirror-checkbox').prop('checked', this.checked);
});

function refreshMirror() {
	location.reload();
}

function showCreateMirrorModal() {
	$('#mirror-modal-title').text('新建镜像会话');
	$('#mirror-form')[0].reset();
	$('#mirror-id').val('');
	$('#mirrorModal').modal('show');
}

function editMirror(id) {
	$('#mirror-modal-title').text('编辑镜像会话 ' + id);
	$('#mirror-id').val(id);
	$('#mirrorModal').modal('show');
	// TODO: 加载镜像配置
}

function saveObserverConfig() {
	var data = {
		observer_port: $('#observer-port').val(),
		observer_vlan: $('#observer-vlan').val() ? parseInt($('#observer-vlan').val()) : null
	};

	if (!data.observer_port) {
		layer.msg('请选择观察端口', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/port-mirror/observer',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('观察端口配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function saveMirror() {
	var selectedPorts = [];
	$('.source-port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	if (selectedPorts.length === 0) {
		layer.msg('请至少选择一个源端口', {icon: 2});
		return;
	}

	var data = {
		name: $('#mirror-name').val(),
		source_ports: selectedPorts,
		dest_port: $('#dest-port').val(),
		direction: $('#mirror-direction').val(),
		status: $('#mirror-status').val()
	};

	var url = '/api/v1/config/port-mirror';
	var method = 'POST';
	var id = $('#mirror-id').val();
	if (id) {
		url = url + '/' + id;
	}

	$.ajax({
		url: url,
		type: method,
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('保存成功', {icon: 1});
				$('#mirrorModal').modal('hide');
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

function deleteMirror(id) {
	if (confirm('确定要删除镜像会话 ' + id + ' 吗？')) {
		$.ajax({
			url: '/api/v1/config/port-mirror/' + id,
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('删除成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
				} else {
					layer.msg(res.message || '删除失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('删除失败，请重试', {icon: 2});
			}
		});
	}
}

function deleteSelectedMirrors() {
	var selected = [];
	$('.mirror-checkbox:checked').each(function() {
		selected.push(parseInt($(this).val()));
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的镜像会话', {icon: 2});
		return;
	}
	if (confirm('确定要删除选中的 ' + selected.length + ' 个镜像会话吗？')) {
		$.ajax({
			url: '/api/v1/config/port-mirror',
			type: 'DELETE',
			contentType: 'application/json',
			data: JSON.stringify({ ids: selected }),
			success: function(res) {
				if (res.code === 200) {
					layer.msg('批量删除成功', {icon: 1});
					refreshMirror();
				} else {
					layer.msg(res.message || '删除失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('删除失败，请重试', {icon: 2});
			}
		});
	}
}

function showFlowMirrorModal() {
	layer.msg('流镜像配置功能待实现', {icon: 3});
}

function editFlowMirror(id) {
	layer.msg('流镜像编辑功能待实现，ID: ' + id, {icon: 3});
}

function deleteFlowMirror(id) {
	if (confirm('确定要删除流镜像 ' + id + ' 吗？')) {
		layer.msg('删除功能待实现', {icon: 3});
	}
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

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口镜像",
		Description: "配置端口镜像 (SPAN) 功能，将源端口的流量复制到监控端口进行分析。支持本地镜像、远程镜像、流镜像",
	}

	return panel, nil
}
