package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPortIsolationContent 端口隔离页面
func getPortIsolationContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	isolationContent := `
<style>
	.isolation-section { margin-bottom: 30px; }
	.isolation-table { font-size: 13px; }
	.isolation-table .status-enabled { color: #28a745; font-weight: bold; }
	.isolation-table .status-disabled { color: #dc3545; }
	.isolation-group { background: #f0f7ff; padding: 10px; border-radius: 4px; margin-bottom: 10px; }
	.port-tag { display: inline-block; padding: 3px 8px; margin: 2px; background: #337ab7; color: white; border-radius: 3px; font-size: 12px; }
	.btn-add-port { margin-top: 10px; }
</style>

<div class="isolation-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">全局隔离模式</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveGlobalIsolation()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">隔离模式</label>
				<div class="col-sm-4">
					<select class="form-control" id="isolation-mode">
						<option value="disabled">禁用</option>
						<option value="l2">二层隔离</option>
						<option value="all">二层三层都隔离</option>
					</select>
					<small class="text-muted">二层隔离：隔离组内端口二层不通，但可通过上行端口通信</small>
				</div>
				<label class="col-sm-2 control-label">上行端口</label>
				<div class="col-sm-4">
					<select class="form-control" id="uplink-port">
						<option value="">- 无 -</option>
						<option value="GE1/0/1">GE1/0/1</option>
						<option value="GE1/0/2">GE1/0/2</option>
						<option value="GE1/0/3">GE1/0/3</option>
						<option value="GE1/0/4">GE1/0/4</option>
					</select>
					<small class="text-muted">隔离组内端口可与上行端口通信</small>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="isolation-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">端口隔离组</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateGroupModal()">
					<i class="fa fa-plus"></i> 新建隔离组
				</button>
				<button type="button" class="btn btn-default btn-sm" onclick="refreshIsolation()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<div id="isolation-groups">
				<div class="isolation-group">
					<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
						<strong>隔离组 1</strong>
						<div>
							<button class="btn btn-sm btn-primary" onclick="editGroup(1)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteGroup(1)">删除</button>
						</div>
					</div>
					<div>
						<span class="port-tag">GE1/0/1</span>
						<span class="port-tag">GE1/0/2</span>
						<span class="port-tag">GE1/0/3</span>
						<span class="port-tag">GE1/0/4</span>
					</div>
					<div style="margin-top: 8px; font-size: 12px; color: #666;">
						<i class="fa fa-info-circle"></i> 组内端口互相隔离，可通过上行端口 GE1/0/24 通信
					</div>
				</div>
				<div class="isolation-group">
					<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
						<strong>隔离组 2</strong>
						<div>
							<button class="btn btn-sm btn-primary" onclick="editGroup(2)">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteGroup(2)">删除</button>
						</div>
					</div>
					<div>
						<span class="port-tag">GE1/0/5</span>
						<span class="port-tag">GE1/0/6</span>
					</div>
					<div style="margin-top: 8px; font-size: 12px; color: #666;">
						<i class="fa fa-info-circle"></i> 组内端口互相隔离，无上行端口
					</div>
				</div>
			</div>
			<div id="no-groups" style="display: none; text-align: center; color: #999; padding: 30px;">
				暂无隔离组
			</div>
		</div>
	</div>
</div>

<div class="isolation-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">端口隔离状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshPortStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover isolation-table">
				<thead>
					<tr>
						<th width="100">端口</th>
						<th width="100">隔离状态</th>
						<th width="150">所属隔离组</th>
						<th width="150">上行端口</th>
						<th>隔离模式</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="port-status-body">
					<tr>
						<td>GE1/0/1</td>
						<td><span class="status-enabled">已隔离</span></td>
						<td>组 1</td>
						<td>GE1/0/24</td>
						<td>L2</td>
						<td><button class="btn btn-sm btn-default" onclick="removeFromGroup('GE1/0/1')">移除</button></td>
					</tr>
					<tr>
						<td>GE1/0/2</td>
						<td><span class="status-enabled">已隔离</span></td>
						<td>组 1</td>
						<td>GE1/0/24</td>
						<td>L2</td>
						<td><button class="btn btn-sm btn-default" onclick="removeFromGroup('GE1/0/2')">移除</button></td>
					</tr>
					<tr>
						<td>GE1/0/7</td>
						<td><span class="status-disabled">未隔离</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-primary" onclick="addToGroup('GE1/0/7')">加入</button></td>
					</tr>
					<tr>
						<td>GE1/0/8</td>
						<td><span class="status-disabled">未隔离</span></td>
						<td>-</td>
						<td>-</td>
						<td>-</td>
						<td><button class="btn btn-sm btn-primary" onclick="addToGroup('GE1/0/8')">加入</button></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑隔离组弹窗 -->
<div class="modal fade" id="groupModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="group-modal-title">新建隔离组</h4>
			</div>
			<div class="modal-body">
				<form id="group-form">
					<input type="hidden" id="group-id">
					<div class="form-group">
						<label>组 ID</label>
						<input type="number" class="form-control" id="group-number" min="1" max="64" placeholder="1-64">
					</div>
					<div class="form-group">
						<label>选择成员端口</label>
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
					<div class="form-group">
						<label>上行端口 (可选)</label>
						<select class="form-control" id="group-uplink">
							<option value="">- 无 -</option>
							<option value="GE1/0/23">GE1/0/23</option>
							<option value="GE1/0/24">GE1/0/24</option>
						</select>
					</div>
					<div class="alert alert-info">
						<i class="fa fa-info-circle"></i> 隔离组内端口之间二层流量互相隔离，但可以与上行端口通信。
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveGroup()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshIsolation() {
	location.reload();
}

function refreshPortStatus() {
	location.reload();
}

function showCreateGroupModal() {
	$('#group-modal-title').text('新建隔离组');
	$('#group-form')[0].reset();
	$('#group-id').val('');
	$('#groupModal').modal('show');
}

function editGroup(groupId) {
	$('#group-modal-title').text('编辑隔离组 ' + groupId);
	$('#group-id').val(groupId);
	$('#group-number').val(groupId);
	$('#groupModal').modal('show');
}

function deleteGroup(groupId) {
	if (confirm('确定要删除隔离组 ' + groupId + ' 吗？组内端口将恢复正常通信。')) {
		$.ajax({
			url: '/api/v1/config/port-isolation/group/' + groupId,
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

function saveGroup() {
	var selectedPorts = [];
	$('.port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	if (selectedPorts.length === 0) {
		layer.msg('请至少选择一个端口', {icon: 2});
		return;
	}

	var data = {
		group_id: parseInt($('#group-number').val()),
		ports: selectedPorts,
		uplink: $('#group-uplink').val()
	};

	$.ajax({
		url: '/api/v1/config/port-isolation/group',
		type: 'POST',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('保存成功', {icon: 1});
				$('#groupModal').modal('hide');
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

function removeFromGroup(portId) {
	if (confirm('确定要将端口 ' + portId + ' 从隔离组移除吗？')) {
		$.ajax({
			url: '/api/v1/config/port-isolation/port/' + portId,
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('移除成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 1000);
				} else {
					layer.msg(res.message || '移除失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('移除失败，请重试', {icon: 2});
			}
		});
	}
}

function addToGroup(portId) {
	layer.msg('请选择要加入的隔离组', {icon: 3});
	showCreateGroupModal();
}

function saveGlobalIsolation() {
	var data = {
		mode: $('#isolation-mode').val(),
		uplink_port: $('#uplink-port').val()
	};

	$.ajax({
		url: '/api/v1/config/port-isolation/global',
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
</script>
`

	boxContent := template.HTML(isolationContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-chain-broken"></i> 端口隔离`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "端口隔离",
		Description: "配置端口隔离组，实现组内端口二层流量隔离，增强网络安全性和隐私保护",
	}

	return panel, nil
}
