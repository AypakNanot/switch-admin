package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getMulticastContent 组播配置页面
func getMulticastContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	multicastContent := `
<style>
	.mc-section { margin-bottom: 30px; }
	.mc-table { font-size: 13px; }
	.mc-table .status-active { color: #28a745; font-weight: bold; }
	.mc-table .status-inactive { color: #dc3545; }
	.group-tag { display: inline-block; padding: 3px 8px; background: #337ab7; color: white; border-radius: 3px; font-size: 12px; margin: 2px; font-family: monospace; }
	.mode-tag { display: inline-block; padding: 2px 6px; border-radius: 3px; font-size: 11px; }
	.mode-igmp { background: #5bc0de; color: white; }
	.mode-snooping { background: #5cb85c; color: white; }
	.mode-proxy { background: #f0ad4e; color: white; }
</style>

<div class="mc-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">IGMP Snooping 配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveIGMPSnooping()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">IGMP Snooping</label>
				<div class="col-sm-3">
					<select class="form-control" id="igmp-snooping-enabled">
						<option value="true">启用</option>
						<option value="false">关闭</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">版本</label>
				<div class="col-sm-3">
					<select class="form-control" id="igmp-version">
						<option value="1">IGMP v1</option>
						<option value="2">IGMP v2</option>
						<option value="3">IGMP v3</option>
					</select>
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-default btn-sm" onclick="resetIGMPDefault()">
						<i class="fa fa-undo"></i> 恢复默认
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">快速离开</label>
				<div class="col-sm-3">
					<select class="form-control" id="igmp-fast-leave">
						<option value="true">启用</option>
						<option value="false">禁用</option>
					</select>
					<small class="text-muted">启用后，收到 Leave 消息立即停止转发组播</small>
				</div>
				<label class="col-sm-2 control-label">路由器端口学习</label>
				<div class="col-sm-3">
					<select class="form-control" id="igmp-router-port">
						<option value="auto">自动学习</option>
						<option value="manual">手动配置</option>
					</select>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">查询间隔</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="igmp-query-interval" min="10" max="1000" value="125">
						<span class="input-group-addon">秒</span>
					</div>
				</div>
				<label class="col-sm-2 control-label">最大响应时间</label>
				<div class="col-sm-3">
					<div class="input-group">
						<input type="number" class="form-control" id="igmp-max-response" min="1" max="100" value="10">
						<span class="input-group-addon">秒</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="mc-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">组播组列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshMulticast()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-danger btn-sm" onclick="clearMulticastGroups()">
					<i class="fa fa-trash"></i> 清除动态组
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover mc-table">
				<thead>
					<tr>
						<th width="50">VLAN</th>
						<th width="180">组播组地址</th>
						<th width="150">成员端口</th>
						<th width="120">路由器端口</th>
						<th width="100">类型</th>
						<th width="100">状态</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="mc-group-body">
					<tr>
						<td>10</td>
						<td><span class="group-tag">239.1.1.1</span></td>
						<td>GE1/0/1, GE1/0/2</td>
						<td>GE1/0/24</td>
						<td><span class="mode-tag mode-snooping">动态</span></td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewGroupDetail('239.1.1.1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>10</td>
						<td><span class="group-tag">239.2.2.2</span></td>
						<td>GE1/0/3</td>
						<td>GE1/0/24</td>
						<td><span class="mode-tag mode-snooping">动态</span></td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewGroupDetail('239.2.2.2')">详情</button>
						</td>
					</tr>
					<tr>
						<td>20</td>
						<td><span class="group-tag">224.0.1.1</span></td>
						<td>GE1/0/5, GE1/0/6</td>
						<td>GE1/0/24</td>
						<td><span class="mode-tag mode-snooping">动态</span></td>
						<td><span class="status-active">活跃</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewGroupDetail('224.0.1.1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>100</td>
						<td><span class="group-tag">239.100.1.1</span></td>
						<td>GE1/0/10</td>
						<td>-</td>
						<td><span class="mode-tag mode-snooping">动态</span></td>
						<td><span class="status-inactive">无成员</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewGroupDetail('239.100.1.1')">详情</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="mc-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">静态组播组配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showAddStaticGroupModal()">
					<i class="fa fa-plus"></i> 添加
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="50">VLAN</th>
						<th width="180">组播组地址</th>
						<th>成员端口</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="static-mc-body">
					<tr>
						<td>10</td>
						<td><span class="group-tag">225.1.1.1</span></td>
						<td>GE1/0/1, GE1/0/2, GE1/0/3</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteStaticGroup('225.1.1.1')">删除</button>
						</td>
					</tr>
					<tr>
						<td>20</td>
						<td><span class="group-tag">225.2.2.2</span></td>
						<td>GE1/0/5, GE1/0/6</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteStaticGroup('225.2.2.2')">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="mc-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">组播 VLAN 配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="saveMulticastVLAN()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">组播 VLAN</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="multicast-vlan" min="1" max="4094" value="100">
					<small class="text-muted">用于承载组播流量的专用 VLAN</small>
				</div>
				<label class="col-sm-2 control-label">组播 VLAN 复制</label>
				<div class="col-sm-3">
					<select class="form-control" id="multicast-vlan-replicate">
						<option value="true">启用</option>
						<option value="false">禁用</option>
					</select>
					<small class="text-muted">启用后将组播流量复制到多个用户 VLAN</small>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">用户 VLAN 列表</label>
				<div class="col-sm-9">
					<input type="text" class="form-control" id="user-vlan-list" placeholder="多个 VLAN 用逗号分隔，如：10,20,30">
					<small class="text-muted">当启用 VLAN 复制时，组播流量会复制到这些 VLAN</small>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- 静态组播组添加弹窗 -->
<div class="modal fade" id="staticGroupModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">添加静态组播组</h4>
			</div>
			<div class="modal-body">
				<form id="static-group-form">
					<div class="form-group">
						<label>VLAN ID</label>
						<input type="number" class="form-control" id="static-vlan" min="1" max="4094" value="1">
					</div>
					<div class="form-group">
						<label>组播组地址</label>
						<input type="text" class="form-control" id="static-group-address" placeholder="224.0.0.0 - 239.255.255.255">
						<small class="text-muted">D 类 IP 地址范围：224.0.0.0 - 239.255.255.255</small>
					</div>
					<div class="form-group">
						<label>成员端口</label>
						<div id="static-member-ports">
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/1"> GE1/0/1</label>
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/2"> GE1/0/2</label>
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/3"> GE1/0/3</label>
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/4"> GE1/0/4</label>
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/5"> GE1/0/5</label>
							<label class="checkbox-inline"><input type="checkbox" class="member-port-checkbox" value="GE1/0/6"> GE1/0/6</label>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveStaticGroup()">添加</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshMulticast() {
	location.reload();
}

function resetIGMPDefault() {
	if (confirm('确定要恢复 IGMP Snooping 默认配置吗？')) {
		$.ajax({
			url: '/api/v1/config/multicast/igmp/reset',
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

function saveIGMPSnooping() {
	var data = {
		enabled: $('#igmp-snooping-enabled').val() === 'true',
		version: parseInt($('#igmp-version').val()),
		fast_leave: $('#igmp-fast-leave').val() === 'true',
		router_port_mode: $('#igmp-router-port').val(),
		query_interval: parseInt($('#igmp-query-interval').val()),
		max_response_time: parseInt($('#igmp-max-response').val())
	};

	$.ajax({
		url: '/api/v1/config/multicast/igmp',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('IGMP Snooping 配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function clearMulticastGroups() {
	if (confirm('确定要清除所有动态学习的组播组吗？静态组不受影响。')) {
		$.ajax({
			url: '/api/v1/config/multicast/groups/dynamic',
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('动态组清除成功', {icon: 1});
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

function viewGroupDetail(group) {
	layer.msg('查看组播组详情：' + group, {icon: 3});
}

function showAddStaticGroupModal() {
	$('#static-group-form')[0].reset();
	$('#staticGroupModal').modal('show');
}

function saveStaticGroup() {
	var selectedPorts = [];
	$('.member-port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	var data = {
		vlan_id: parseInt($('#static-vlan').val()),
		group_address: $('#static-group-address').val(),
		member_ports: selectedPorts
	};

	if (!data.group_address) {
		layer.msg('请输入组播组地址', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/multicast/group/static',
		type: 'POST',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('静态组播组添加成功', {icon: 1});
				$('#staticGroupModal').modal('hide');
				setTimeout(function() { location.reload(); }, 1000);
			} else {
				layer.msg(res.message || '添加失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('添加失败，请重试', {icon: 2});
		}
	});
}

function deleteStaticGroup(group) {
	if (confirm('确定要删除静态组播组 ' + group + ' 吗？')) {
		$.ajax({
			url: '/api/v1/config/multicast/group/static/' + encodeURIComponent(group),
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

function saveMulticastVLAN() {
	var data = {
		vlan_id: parseInt($('#multicast-vlan').val()),
		replicate: $('#multicast-vlan-replicate').val() === 'true',
		user_vlans: $('#user-vlan-list').val().split(',').map(function(v) { return parseInt(v.trim()); }).filter(function(v) { return !isNaN(v); })
	};

	$.ajax({
		url: '/api/v1/config/multicast/vlan',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('组播 VLAN 配置保存成功', {icon: 1});
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

	boxContent := template.HTML(multicastContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-sitemap"></i> 组播配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "组播配置",
		Description: "配置组播功能，包括 IGMP Snooping、组播组管理、静态组播组、组播 VLAN 等，优化组播流量转发",
	}

	return panel, nil
}
