package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getVLANContent VLAN 配置页面
func getVLANContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	vlanContent := `
<style>
	.vlan-section { margin-bottom: 30px; }
	.vlan-table { font-size: 13px; }
	.vlan-table .status-active { color: #28a745; font-weight: bold; }
	.vlan-table .status-inactive { color: #dc3545; }
	.port-tag { display: inline-block; padding: 2px 6px; margin: 1px; background: #337ab7; color: white; border-radius: 3px; font-size: 11px; }
	.tag-untagged { background: #f0ad4e; }
	.tag-tagged { background: #5cb85c; }
</style>

<div class="vlan-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">创建 VLAN</h3>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">VLAN ID</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="vlan-id" min="1" max="4094" placeholder="1-4094">
				</div>
				<label class="col-sm-2 control-label">VLAN 名称</label>
				<div class="col-sm-3">
					<input type="text" class="form-control" id="vlan-name" maxlength="32" placeholder="最多 32 字符">
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-primary" onclick="createVLAN()">
						<i class="fa fa-plus"></i> 创建
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">VLAN 类型</label>
				<div class="col-sm-3">
					<select class="form-control" id="vlan-type">
						<option value="standard">标准 VLAN</option>
						<option value="voice">语音 VLAN</option>
						<option value="guest">访客 VLAN</option>
						<option value="management">管理 VLAN</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="vlan-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">VLAN 列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelectedVLANs()">
					<i class="fa fa-trash"></i> 批量删除
				</button>
				<button type="button" class="btn btn-default btn-sm" onclick="refreshVLANs()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover vlan-table">
				<thead>
					<tr>
						<th width="40"><input type="checkbox" id="select-all"></th>
						<th width="80">VLAN ID</th>
						<th width="150">名称</th>
						<th width="100">类型</th>
						<th>端口成员</th>
						<th width="100">状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="vlan-table-body">
					<tr>
						<td><input type="checkbox" class="vlan-checkbox" value="1"></td>
						<td>1</td>
						<td>default</td>
						<td>standard</td>
						<td>
							<span class="port-tag tag-untagged">GE1/0/1</span>
							<span class="port-tag tag-untagged">GE1/0/2</span>
							<span class="port-tag tag-untagged">GE1/0/3</span>
							<span class="port-tag tag-untagged">GE1/0/4</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVLAN(1)">编辑</button>
							<button class="btn btn-sm btn-default" onclick="viewVLANPorts(1)">端口</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="vlan-checkbox" value="10"></td>
						<td>10</td>
						<td>Office</td>
						<td>standard</td>
						<td>
							<span class="port-tag tag-untagged">GE1/0/5</span>
							<span class="port-tag tag-untagged">GE1/0/6</span>
							<span class="port-tag tag-tagged">GE1/0/24</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVLAN(10)">编辑</button>
							<button class="btn btn-sm btn-default" onclick="viewVLANPorts(10)">端口</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="vlan-checkbox" value="20"></td>
						<td>20</td>
						<td>Guest</td>
						<td>guest</td>
						<td>
							<span class="port-tag tag-untagged">GE1/0/7</span>
							<span class="port-tag tag-untagged">GE1/0/8</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVLAN(20)">编辑</button>
							<button class="btn btn-sm btn-default" onclick="viewVLANPorts(20)">端口</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="vlan-checkbox" value="100"></td>
						<td>100</td>
						<td>Voice</td>
						<td>voice</td>
						<td>
							<span class="port-tag tag-tagged">GE1/0/1</span>
							<span class="port-tag tag-tagged">GE1/0/2</span>
						</td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editVLAN(100)">编辑</button>
							<button class="btn btn-sm btn-default" onclick="viewVLANPorts(100)">端口</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- VLAN 编辑弹窗 -->
<div class="modal fade" id="vlanEditModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">编辑 VLAN - <span id="edit-vlan-id"></span></h4>
			</div>
			<div class="modal-body">
				<form id="vlan-edit-form">
					<div class="box box-default">
						<div class="box-body">
							<div class="form-group">
								<label>VLAN 名称</label>
								<input type="text" class="form-control" id="edit-vlan-name" maxlength="32">
							</div>
							<div class="form-group">
								<label>VLAN 类型</label>
								<select class="form-control" id="edit-vlan-type">
									<option value="standard">标准 VLAN</option>
									<option value="voice">语音 VLAN</option>
									<option value="guest">访客 VLAN</option>
									<option value="management">管理 VLAN</option>
								</select>
							</div>
							<div class="form-group">
								<label>Untagged 端口</label>
								<div id="untagged-ports">
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/1"> GE1/0/1</label>
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/2"> GE1/0/2</label>
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/3"> GE1/0/3</label>
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/4"> GE1/0/4</label>
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/5"> GE1/0/5</label>
									<label class="checkbox-inline"><input type="checkbox" class="untagged-checkbox" value="GE1/0/6"> GE1/0/6</label>
								</div>
								<small class="text-muted">Untagged 端口：发送数据时剥离 VLAN 标签</small>
							</div>
							<div class="form-group">
								<label>Tagged 端口</label>
								<div id="tagged-ports">
									<label class="checkbox-inline"><input type="checkbox" class="tagged-checkbox" value="GE1/0/23"> GE1/0/23</label>
									<label class="checkbox-inline"><input type="checkbox" class="tagged-checkbox" value="GE1/0/24"> GE1/0/24</label>
								</div>
								<small class="text-muted">Tagged 端口：发送数据时保留 VLAN 标签 (用于 trunk)</small>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveVLANConfig()">保存</button>
			</div>
		</div>
	</div>
</div>

<!-- VLAN 端口详情弹窗 -->
<div class="modal fade" id="vlanPortsModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">VLAN 端口成员 - VLAN <span id="view-vlan-id"></span></h4>
			</div>
			<div class="modal-body">
				<div class="box box-default">
					<div class="box-body">
						<h5>Untagged 端口:</h5>
						<div id="view-untagged-ports" style="margin-bottom: 15px;"></div>
						<h5>Tagged 端口:</h5>
						<div id="view-tagged-ports"></div>
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
$('#select-all').change(function() {
	$('.vlan-checkbox').prop('checked', this.checked);
});

function refreshVLANs() {
	location.reload();
}

function createVLAN() {
	var id = parseInt($('#vlan-id').val());
	var name = $('#vlan-name').val();
	var type = $('#vlan-type').val();

	if (!id || id < 1 || id > 4094) {
		layer.msg('请输入有效的 VLAN ID (1-4094)', {icon: 2});
		return;
	}
	if (!name) {
		layer.msg('请输入 VLAN 名称', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/vlan',
		type: 'POST',
		contentType: 'application/json',
		data: JSON.stringify({
			vlan_id: id,
			name: name,
			type: type
		}),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('VLAN 创建成功', {icon: 1});
				$('#vlan-id').val('');
				$('#vlan-name').val('');
				refreshVLANs();
			} else {
				layer.msg(res.message || '创建失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('创建失败，请重试', {icon: 2});
		}
	});
}

function editVLAN(vlanId) {
	$('#edit-vlan-id').text(vlanId);
	$('#vlanEditModal').modal('show');
	// TODO: 加载 VLAN 配置
}

function viewVLANPorts(vlanId) {
	$('#view-vlan-id').text(vlanId);
	$('#vlanPortsModal').modal('show');
	// TODO: 加载 VLAN 端口成员
	$('#view-untagged-ports').html('<span class="port-tag tag-untagged">GE1/0/1</span> <span class="port-tag tag-untagged">GE1/0/2</span>');
	$('#view-tagged-ports').html('<span class="port-tag tag-tagged">GE1/0/24</span>');
}

function saveVLANConfig() {
	var vlanId = $('#edit-vlan-id').text();
	var untaggedPorts = [];
	var taggedPorts = [];

	$('.untagged-checkbox:checked').each(function() {
		untaggedPorts.push($(this).val());
	});
	$('.tagged-checkbox:checked').each(function() {
		taggedPorts.push($(this).val());
	});

	var data = {
		name: $('#edit-vlan-name').val(),
		type: $('#edit-vlan-type').val(),
		untagged_ports: untaggedPorts,
		tagged_ports: taggedPorts
	};

	$.ajax({
		url: '/api/v1/config/vlan/' + vlanId,
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('VLAN 配置保存成功', {icon: 1});
				$('#vlanEditModal').modal('hide');
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

function deleteSelectedVLANs() {
	var selected = [];
	$('.vlan-checkbox:checked').each(function() {
		selected.push(parseInt($(this).val()));
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的 VLAN', {icon: 2});
		return;
	}
	if (selected.indexOf(1) !== -1) {
		layer.msg('默认 VLAN 1 不能删除', {icon: 2});
		return;
	}
	if (confirm('确定要删除选中的 ' + selected.length + ' 个 VLAN 吗？')) {
		$.ajax({
			url: '/api/v1/config/vlan',
			type: 'DELETE',
			contentType: 'application/json',
			data: JSON.stringify({ vlan_ids: selected }),
			success: function(res) {
				if (res.code === 200) {
					layer.msg('VLAN 批量删除成功', {icon: 1});
					refreshVLANs();
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
</script>
`

	boxContent := template.HTML(vlanContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-sitemap"></i> VLAN 配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "VLAN 配置",
		Description: "配置虚拟局域网 (VLAN)，实现网络分段和广播域隔离。支持标准 VLAN、语音 VLAN、访客 VLAN 等类型",
	}

	return panel, nil
}
