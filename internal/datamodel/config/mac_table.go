package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getMacTableContent MAC 地址表页面
func getMacTableContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	macContent := `
<style>
	.mac-section { margin-bottom: 30px; }
	.mac-table { font-size: 12px; }
	.mac-table .type-static { color: #28a745; font-weight: bold; }
	.mac-table .type-dynamic { color: #337ab7; }
	.mac-table .type-blackhole { color: #dc3545; }
	.search-box { display: inline-block; margin-left: 10px; }
	.mono-font { font-family: 'Consolas', 'Monaco', monospace; }
</style>

<div class="mac-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">MAC 地址表查询</h3>
			<div class="box-tools">
				<div class="search-box">
					<input type="text" class="form-control" id="mac-search" placeholder="MAC 地址搜索" style="display: inline-block; width: 180px;">
					<button type="button" class="btn btn-primary btn-sm" onclick="searchMAC()">
						<i class="fa fa-search"></i> 搜索
					</button>
				</div>
				<button type="button" class="btn btn-default btn-sm" onclick="refreshMacTable()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-success btn-sm" onclick="exportMacTable()">
					<i class="fa fa-download"></i> 导出
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">地址类型</label>
				<div class="col-sm-3">
					<select class="form-control" id="mac-type-filter">
						<option value="all">全部</option>
						<option value="dynamic">动态学习</option>
						<option value="static">静态配置</option>
						<option value="blackhole">黑洞</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">所属 VLAN</label>
				<div class="col-sm-3">
					<select class="form-control" id="vlan-filter">
						<option value="all">全部</option>
						<option value="1">VLAN 1</option>
						<option value="10">VLAN 10</option>
						<option value="20">VLAN 20</option>
						<option value="100">VLAN 100</option>
					</select>
				</div>
				<div class="col-sm-2">
					<button type="button" class="btn btn-info btn-sm" onclick="applyFilter()">
						<i class="fa fa-filter"></i> 应用筛选
					</button>
				</div>
			</div>
		</div>
	</div>
</div>

<div class="mac-section">
	<div class="box box-info">
		<div class="box-header with-border">
			<h3 class="box-title">MAC 地址表 (<span id="mac-count">0</span> 条)</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-warning btn-sm" onclick="clearDynamicMAC()">
					<i class="fa fa-eraser"></i> 清除动态表项
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover mac-table">
				<thead>
					<tr>
						<th width="50">序号</th>
						<th width="180">MAC 地址</th>
						<th width="80">VLAN</th>
						<th width="120">端口</th>
						<th width="100">类型</th>
						<th width="150">老化时间</th>
						<th width="120">操作</th>
					</tr>
				</thead>
				<tbody id="mac-table-body">
					<tr>
						<td>1</td>
						<td class="mono-font">00:1A:2B:3C:4D:5E</td>
						<td>10</td>
						<td>GE1/0/1</td>
						<td><span class="type-dynamic">动态学习</span></td>
						<td>280s</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMACDetail('00:1A:2B:3C:4D:5E')">详情</button>
						</td>
					</tr>
					<tr>
						<td>2</td>
						<td class="mono-font">00:1A:2B:3C:4D:5F</td>
						<td>10</td>
						<td>GE1/0/2</td>
						<td><span class="type-dynamic">动态学习</span></td>
						<td>150s</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMACDetail('00:1A:2B:3C:4D:5F')">详情</button>
						</td>
					</tr>
					<tr>
						<td>3</td>
						<td class="mono-font">AA:BB:CC:DD:EE:01</td>
						<td>20</td>
						<td>GE1/0/5</td>
						<td><span class="type-static">静态配置</span></td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteStaticMAC('AA:BB:CC:DD:EE:01')">删除</button>
						</td>
					</tr>
					<tr>
						<td>4</td>
						<td class="mono-font">00:00:00:00:00:00</td>
						<td>1</td>
						<td>-</td>
						<td><span class="type-blackhole">黑洞</span></td>
						<td>-</td>
						<td>
							<button class="btn btn-sm btn-danger" onclick="deleteStaticMAC('00:00:00:00:00:00')">删除</button>
						</td>
					</tr>
					<tr>
						<td>5</td>
						<td class="mono-font">11:22:33:44:55:66</td>
						<td>100</td>
						<td>GE1/0/24</td>
						<td><span class="type-dynamic">动态学习</span></td>
						<td>300s</td>
						<td>
							<button class="btn btn-sm btn-default" onclick="viewMACDetail('11:22:33:44:55:66')">详情</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<div class="mac-section">
	<div class="box box-warning">
		<div class="box-header with-border">
			<h3 class="box-title">静态 MAC 地址配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary btn-sm" onclick="addStaticMAC()">
					<i class="fa fa-plus"></i> 添加
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group row">
				<label class="col-sm-2 control-label">MAC 地址</label>
				<div class="col-sm-3">
					<input type="text" class="form-control" id="static-mac-address" placeholder="XX:XX:XX:XX:XX:XX">
					<small class="text-muted">格式：00:1A:2B:3C:4D:5E</small>
				</div>
				<label class="col-sm-2 control-label">VLAN ID</label>
				<div class="col-sm-3">
					<input type="number" class="form-control" id="static-mac-vlan" min="1" max="4094" value="1">
				</div>
			</div>
			<div class="form-group row">
				<label class="col-sm-2 control-label">绑定端口</label>
				<div class="col-sm-3">
					<select class="form-control" id="static-mac-port">
						<option value="GE1/0/1">GE1/0/1</option>
						<option value="GE1/0/2">GE1/0/2</option>
						<option value="GE1/0/3">GE1/0/3</option>
						<option value="GE1/0/4">GE1/0/4</option>
						<option value="GE1/0/5">GE1/0/5</option>
						<option value="GE1/0/6">GE1/0/6</option>
					</select>
				</div>
				<label class="col-sm-2 control-label">类型</label>
				<div class="col-sm-3">
					<select class="form-control" id="static-mac-type">
						<option value="static">静态绑定</option>
						<option value="blackhole">黑洞</option>
					</select>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- MAC 详情弹窗 -->
<div class="modal fade" id="macDetailModal" tabindex="-1">
	<div class="modal-dialog">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title">MAC 地址详情 - <span id="detail-mac-address"></span></h4>
			</div>
			<div class="modal-body">
				<table class="table">
					<tr><td width="150">MAC 地址</td><td class="mono-font" id="detail-mac">-</td></tr>
					<tr><td>VLAN ID</td><td id="detail-vlan">-</td></tr>
					<tr><td>端口</td><td id="detail-port">-</td></tr>
					<tr><td>类型</td><td id="detail-type">-</td></tr>
					<tr><td>老化时间</td><td id="detail-age">-</td></tr>
					<tr><td>OUI 厂商</td><td id="detail-oui">-</td></tr>
				</table>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
			</div>
		</div>
	</div>
</div>

<script>
function refreshMacTable() {
	location.reload();
}

function searchMAC() {
	var mac = $('#mac-search').val().trim();
	if (!mac) {
		layer.msg('请输入 MAC 地址', {icon: 2});
		return;
	}
	// TODO: 调用 API 搜索
	layer.msg('搜索功能待实现：' + mac, {icon: 3});
}

function applyFilter() {
	var type = $('#mac-type-filter').val();
	var vlan = $('#vlan-filter').val();
	// TODO: 调用 API 筛选
	layer.msg('筛选功能待实现：type=' + type + ', vlan=' + vlan, {icon: 3});
}

function clearDynamicMAC() {
	if (confirm('确定要清除所有动态学习的 MAC 地址表项吗？')) {
		$.ajax({
			url: '/api/v1/config/mac-table/dynamic',
			type: 'DELETE',
			success: function(res) {
				if (res.code === 200) {
					layer.msg('动态表项清除成功', {icon: 1});
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

function addStaticMAC() {
	var mac = $('#static-mac-address').val().trim();
	var vlan = parseInt($('#static-mac-vlan').val());
	var port = $('#static-mac-port').val();
	var type = $('#static-mac-type').val();

	if (!mac) {
		layer.msg('请输入 MAC 地址', {icon: 2});
		return;
	}

	$.ajax({
		url: '/api/v1/config/mac-table/static',
		type: 'POST',
		contentType: 'application/json',
		data: JSON.stringify({
			mac_address: mac,
			vlan_id: vlan,
			port: port,
			type: type
		}),
		success: function(res) {
			if (res.code === 200) {
				layer.msg('静态 MAC 添加成功', {icon: 1});
				$('#static-mac-address').val('');
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

function deleteStaticMAC(mac) {
	if (confirm('确定要删除静态 MAC 地址 ' + mac + ' 吗？')) {
		$.ajax({
			url: '/api/v1/config/mac-table/static/' + encodeURIComponent(mac),
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

function viewMACDetail(mac) {
	$('#detail-mac-address').text(mac);
	$('#macDetailModal').modal('show');
	// TODO: 加载 MAC 详情
	$('#detail-mac').text(mac);
	$('#detail-vlan').text('10');
	$('#detail-port').text('GE1/0/1');
	$('#detail-type').text('动态学习');
	$('#detail-age').text('280s');
	$('#detail-oui').text('Unknown (需查询 OUI 数据库)');
}

function exportMacTable() {
	layer.msg('正在导出...', {icon: 1});
	window.location.href = '/api/v1/config/mac-table/export';
}
</script>
`

	boxContent := template.HTML(macContent)
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
		Description: "查看和管理交换机 MAC 地址转发表，支持动态学习、静态配置、黑洞等类型，可搜索、筛选和导出",
	}

	return panel, nil
}
