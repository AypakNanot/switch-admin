package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getStackContent 堆叠配置页面
func getStackContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	stackContent := `
<style>
	.stack-section { margin-bottom: 30px; }
	.stack-table { font-size: 13px; }
	.stack-table .status-master { color: #28a745; font-weight: bold; }
	.stack-table .status-slave { color: #337ab7; }
	.stack-table .status-standby { color: #f0ad4e; }
	.stack-table .status-offline { color: #dc3545; }
	.stack-topo { padding: 20px; background: #f5f5f5; border-radius: 4px; text-align: center; }
	.stack-unit { display: inline-block; width: 150px; height: 120px; border: 3px solid #ddd; border-radius: 8px; background: white; margin: 10px; padding: 10px; position: relative; }
	.stack-unit.master { border-color: #28a745; background: #e8f5e9; }
	.stack-unit.slave { border-color: #337ab7; }
	.stack-unit.offline { border-color: #dc3545; background: #ffebee; opacity: 0.6; }
	.stack-unit .unit-id { font-size: 24px; font-weight: bold; color: #333; }
	.stack-unit .unit-role { font-size: 12px; color: #666; margin-top: 5px; }
	.stack-unit .unit-info { font-size: 11px; color: #999; margin-top: 5px; }
	.stack-link { display: inline-block; width: 50px; height: 4px; background: #28a745; vertical-align: middle; margin: 0 -5px; }
	.stack-link.down { background: #dc3545; }
	.priority-high { color: #28a745; }
	.priority-normal { color: #666; }
</style>

<div class="stack-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">堆叠拓扑</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshStack()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="showStackConfigModal()">
					<i class="fa fa-cog"></i> 堆叠配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="stack-topo">
				<div class="stack-unit master">
					<div class="unit-id">Unit 1</div>
					<div class="unit-role"><span class="status-master">Master</span></div>
					<div class="unit-info">SN: ABC123456</div>
					<div class="unit-info">Priority: 150</div>
					<div class="unit-info">Uptime: 15d 8h</div>
				</div>
				<div class="stack-link"></div>
				<div class="stack-unit slave">
					<div class="unit-id">Unit 2</div>
					<div class="unit-role"><span class="status-slave">Slave</span></div>
					<div class="unit-info">SN: DEF789012</div>
					<div class="unit-info">Priority: 100</div>
					<div class="unit-info">Uptime: 15d 8h</div>
				</div>
				<div class="stack-link"></div>
				<div class="stack-unit slave">
					<div class="unit-id">Unit 3</div>
					<div class="unit-role"><span class="status-slave">Slave</span></div>
					<div class="unit-info">SN: GHI345678</div>
					<div class="unit-info">Priority: 100</div>
					<div class="unit-info">Uptime: 15d 8h</div>
				</div>
				<div class="stack-link down"></div>
				<div class="stack-unit offline">
					<div class="unit-id">Unit 4</div>
					<div class="unit-role"><span class="status-offline">Offline</span></div>
					<div class="unit-info">SN: JKL901234</div>
					<div class="unit-info">Last seen: 2h ago</div>
				</div>
				<div style="margin-top: 30px; text-align: left;">
					<h5>堆叠信息:</h5>
					<table class="table table-condensed">
						<tr><td width="150">堆叠名称:</td><td>Core-Stack</td></tr>
						<tr><td>堆叠 ID:</td><td>1</td></tr>
						<tr><td>成员数量:</td><td>3 / 4 (在线/最大)</td></tr>
						<tr><td>堆叠带宽:</td><td>80 Gbps (双向)</td></tr>
						<tr><td>拓扑类型:</td><td>链型</td></tr>
					</table>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="stack-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">堆叠成员列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-warning btn-sm" onclick="syncConfig()">
					<i class="fa fa-exchange"></i> 同步配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover stack-table">
				<thead>
					<tr>
						<th width="80">成员 ID</th>
						<th width="100">角色</th>
						<th width="150">序列号</th>
						<th width="120">型号</th>
						<th width="100">优先级</th>
						<th width="150">运行时间</th>
						<th width="100">状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="stack-member-body">
					<tr>
						<td>1</td>
						<td><span class="status-master">Master</span></td>
						<td>ABC123456</td>
						<td>S5720-28X</td>
						<td><span class="priority-high">150</span></td>
						<td>15 天 8 小时 30 分</td>
						<td><span class="status-master">正常</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMemberDetail(1)">详情</button>
							<button class="btn btn-sm btn-primary" onclick="changePriority(1)">优先级</button>
						</td>
					</tr>
					<tr>
						<td>2</td>
						<td><span class="status-slave">Slave</span></td>
						<td>DEF789012</td>
						<td>S5720-28X</td>
						<td><span class="priority-normal">100</span></td>
						<td>15 天 8 小时 30 分</td>
						<td><span class="status-master">正常</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMemberDetail(2)">详情</button>
							<button class="btn btn-sm btn-primary" onclick="changePriority(2)">优先级</button>
						</td>
					</tr>
					<tr>
						<td>3</td>
						<td><span class="status-slave">Slave</span></td>
						<td>GHI345678</td>
						<td>S5720-28X</td>
						<td><span class="priority-normal">100</span></td>
						<td>15 天 8 小时 30 分</td>
						<td><span class="status-master">正常</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMemberDetail(3)">详情</button>
							<button class="btn btn-sm btn-primary" onclick="changePriority(3)">优先级</button>
						</td>
					</tr>
					<tr>
						<td>4</td>
						<td><span class="status-offline">-</span></td>
						<td>JKL901234</td>
						<td>S5720-28X</td>
						<td><span class="priority-normal">100</span></td>
						<td>-</td>
						<td><span class="status-offline">离线</span></td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMemberDetail(4)">详情</button>
							<button class="btn btn-sm btn-danger" onclick="removeMember(4)">移除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="stack-section">
	<div class="box box-success">
		<div class="box-header with-border">
			<h3 class="box-title">堆叠端口配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default btn-sm" onclick="refreshStackPort()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="150">成员交换机</th>
						<th>堆叠端口</th>
						<th width="150">对端端口</th>
						<th width="100">链路状态</th>
						<th width="100">带宽</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="stack-port-body">
					<tr>
						<td>Unit 1</td>
						<td>Stack-Port 1</td>
						<td>Unit 2: Stack-Port 2</td>
						<td><span class="status-master">Up</span></td>
						<td>40 Gbps</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>Unit 1</td>
						<td>Stack-Port 2</td>
						<td>-</td>
						<td><span class="status-offline">Down</span></td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 2')">详情</button>
						</td>
					</tr>
					<tr>
						<td>Unit 2</td>
						<td>Stack-Port 1</td>
						<td>Unit 1: Stack-Port 1</td>
						<td><span class="status-master">Up</span></td>
						<td>40 Gbps</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>Unit 2</td>
						<td>Stack-Port 2</td>
						<td>Unit 3: Stack-Port 1</td>
						<td><span class="status-master">Up</span></td>
						<td>40 Gbps</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 2')">详情</button>
						</td>
					</tr>
					<tr>
						<td>Unit 3</td>
						<td>Stack-Port 1</td>
						<td>Unit 2: Stack-Port 2</td>
						<td><span class="status-master">Up</span></td>
						<td>40 Gbps</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 1')">详情</button>
						</td>
					</tr>
					<tr>
						<td>Unit 3</td>
						<td>Stack-Port 2</td>
						<td>Unit 4: Stack-Port 1</td>
						<td><span class="status-offline">Down</span></td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewPortDetail('Stack-Port 2')">详情</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="stack-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">堆叠管理</h3>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">堆叠名称</label>
				<div class="col-sm-3">
					<input type="text" class="form-control" id="stack-name" value="Core-Stack" maxlength="32">
				</div>
				<label class="col-sm-2 control-label">堆叠 ID</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="stack-id" min="1" max="9" value="1">
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-primary btn-sm" onclick="saveStackConfig()">
						<i class="fa fa-save"></i> 保存
					</button>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">主备切换</label>
				<div class="col-sm-9">
					<button type="button" class="btn btn-warning btn-sm" onclick="forceMasterSwap()">
						<i class="fa fa-exchange"></i> 强制主备切换
					</button>
					<small class="text-muted" style="margin-left: 10px;">将当前备用交换机切换为主交换机，会导致网络短暂中断</small>
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">堆叠分裂检测</label>
				<div class="col-sm-3">
					<select class="form-control" id="split-detection">
						<option value="mad">MAD 检测</option>
						<option value="direct">直连检测</option>
						<option value="proxy">代理检测</option>
						<option value="disable">禁用</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">自动合并</label>
				<div class="col-sm-3">
					<select class="form-control" id="auto-merge">
						<option value="true">启用</option>
						<option value="false">禁用</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- 优先级配置弹窗 -->
<div class="modal fade" id="priorityModal" tabindex="-1">
	<div class="modal-dialog modal-sm">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">配置优先级 - Unit <span id="priority-unit-id"></span></h4>
			</div>
			<div class="modal-body">
				<div class="form-group">
					<label>优先级 (1-255)</label>
					<input type="number" class="form-control" id="member-priority" min="1" max="255" value="100">
					<small class="text-muted">值越大优先级越高，主交换机选举时优先选择优先级高的设备</small>
				</div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="savePriority()">保存</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshStack() {
	location.reload();
}

function refreshStackPort() {
	location.reload();
}

function showStackConfigModal() {
	layer.msg('堆叠配置功能待实现', {icon: 3});
}

function syncConfig() {
	if (confirm('确定要将主交换机的配置同步到所有备交换机吗？')) {
		$.ajax({
			url: '/api/v1/config/stack/sync',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('配置同步成功', {icon: 1});
				} else {
					layer.msg(res.message || '同步失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('同步失败，请重试', {icon: 2});
			}
		});
	}
}

function viewMemberDetail(unitId) {
	layer.msg('查看成员详情：Unit ' + unitId, {icon: 3});
}

function changePriority(unitId) {
	$('#priority-unit-id').text(unitId);
	$('#priorityModal').modal('show');
}

function savePriority() {
	var unitId = $('#priority-unit-id').text();
	var priority = parseInt($('#member-priority').val());

	if (priority < 1 || priority > 255) {
		layer.msg('优先级必须在 1-255 范围内', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/stack/member/' + unitId + '/priority',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify({ priority: priority }),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('优先级修改成功', {icon: 1});
				$('#priorityModal').modal('hide');
				setTimeout(function() { location.reload(); }, 1000);
			} else {
				layer.msg(res.message || '修改失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('修改失败，请重试', {icon: 2});
		}
	});
}

function removeMember(unitId) {
	if (confirm('确定要移除离线成员 Unit ' + unitId + ' 吗？')) {
		$.ajax({
			url: '/api/v1/config/stack/member/' + unitId,
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

function forceMasterSwap() {
	if (confirm('警告：强制主备切换会导致网络短暂中断！确定要继续吗？')) {
		$.ajax({
			url: '/api/v1/config/stack/master-swap',
			type: 'POST',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('主备切换成功', {icon: 1});
					setTimeout(function() { location.reload(); }, 2000);
				} else {
					layer.msg(res.message || '切换失败', {icon: 2});
				}
			},
			error: function() {
				layer.msg('切换失败，请重试', {icon: 2});
			}
		});
	}
}

function saveStackConfig() {
	var data = {
		name: $('#stack-name').val(),
		stack_id: parseInt($('#stack-id').val()),
		split_detection: $('#split-detection').val(),
		auto_merge: $('#auto-merge').val() === 'true'
	};

	$.ajax({
		url: '/api/v1/config/stack/global',
		type: 'PUT',
		contentType: 'application/json',
		data: JSON.stringify(data),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('堆叠配置保存成功', {icon: 1});
			} else {
				layer.msg(res.message || '保存失败', {icon: 2});
			}
		},
		error: function() {
			layer.msg('保存失败，请重试', {icon: 2});
		}
	});
}

function viewPortDetail(portName) {
	layer.msg('查看端口详情：' + portName, {icon: 3});
}
</script>
`

	boxContent := template.HTML(stackContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-server"></i> 堆叠配置`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "堆叠配置",
		Description: "配置交换机堆叠功能，实现多台交换机虚拟化为一台逻辑交换机。支持拓扑显示、成员管理、优先级配置、主备切换等",
	}

	return panel, nil
}
