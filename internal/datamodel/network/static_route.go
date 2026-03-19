package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetStaticRouteContent IPv4 静态路由信息页面
func GetStaticRouteContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	// 操作按钮
	actionButtons := `
	<div style="margin-bottom:15px;">
		<button class="btn btn-success" onclick="createStaticRoute()">
			<i class="fa fa-plus"></i> 新建
		</button>
		<button class="btn btn-danger" onclick="deleteStaticRoute()">
			<i class="fa fa-trash"></i> 删除
		</button>
	</div>
	`

	// 静态路由表格
	staticRouteTable := `
	<table class="table table-bordered table-hover" id="static-route-table">
		<thead>
			<tr>
				<th width="30"><input type="checkbox" id="select-all"></th>
				<th>目的 IP 地址</th>
				<th>目的 IP 掩码</th>
				<th>下一跳</th>
				<th>路由距离</th>
				<th>状态</th>
				<th>操作</th>
			</tr>
		</thead>
		<tbody id="static-route-body">
			<tr><td colspan="7" style="text-align:center;">暂无数据~</td></tr>
		</tbody>
	</table>
	`

	// 操作脚本
	actionScript := `
	<style>
	.status-active { color: green; }
	.status-warning { color: orange; }
	</style>
	<script>
	function loadStaticRoutes() {
		fetch('/api/v1/routes/static')
			.then(res => res.json())
			.then(data => {
				renderStaticRoutes(data.data);
			});
	}

	function renderStaticRoutes(data) {
		var tbody = document.getElementById('static-route-body');
		var html = '';
		if (!data.items || data.items.length === 0) {
			tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;">暂无数据~</td></tr>';
			return;
		}
		data.items.forEach(function(item, index) {
			var statusHtml = '<span class="status-active">● 生效</span>';
			if (item.status === 'warning') {
				statusHtml = '<span class="status-warning">⚠️ 下一跳不可达</span>';
			}

			html += '<tr>' +
				'<td><input type="checkbox" name="route-id" value="' + item.id + '"></td>' +
				'<td>' + item.dest_ip + '</td>' +
				'<td>' + item.dest_mask + '</td>' +
				'<td>' + item.next_hop + '</td>' +
				'<td>' + item.distance + '</td>' +
				'<td>' + statusHtml + '</td>' +
				'<td>' +
				'<button class="btn btn-sm btn-primary" onclick="editStaticRoute(\'' + item.id + '\')"><i class="fa fa-edit"></i> 编辑</button> ' +
				'<button class="btn btn-sm btn-danger" onclick="deleteStaticRouteById(\'' + item.id + '\')"><i class="fa fa-trash"></i> 删除</button>' +
				'</td>' +
				'</tr>';
		});
		tbody.innerHTML = html;
	}

	function createStaticRoute() {
		showModal('新建静态路由', '', true);
	}

	function editStaticRoute(id) {
		fetch('/api/v1/routes/static/' + id)
			.then(res => res.json())
			.then(data => {
				showModal('编辑静态路由', data.data, false);
			});
	}

	function showModal(title, routeData, isNew) {
		var destIp = routeData && routeData.dest_ip ? routeData.dest_ip : '';
		var destMask = routeData && routeData.dest_mask ? routeData.dest_mask : '';
		var nextHop = routeData && routeData.next_hop ? routeData.next_hop : '';
		var distance = routeData && routeData.distance ? routeData.distance : '1';

		var html = '<form id="route-form">' +
			'<div class="form-group">' +
			'<label>*目的 IP 地址</label>' +
			'<input type="text" class="form-control" id="form-dest-ip" value="' + destIp + '" placeholder="0.0.0.0">' +
			'</div>' +
			'<div class="form-group">' +
			'<label>*目的 IP 掩码</label>' +
			'<select class="form-control" id="form-dest-mask">' +
			'<option value="0.0.0.0"' + (destMask === '0.0.0.0' ? ' selected' : '') + '>0.0.0.0 (0)</option>' +
			'<option value="255.255.255.0"' + (destMask === '255.255.255.0' ? ' selected' : '') + '>255.255.255.0 (24)</option>' +
			'<option value="255.255.0.0"' + (destMask === '255.255.0.0' ? ' selected' : '') + '>255.255.0.0 (16)</option>' +
			'<option value="255.0.0.0"' + (destMask === '255.0.0.0' ? ' selected' : '') + '>255.0.0.0 (8)</option>' +
			'</select>' +
			'</div>' +
			'<div class="form-group">' +
			'<label>下一跳</label>' +
			'<input type="text" class="form-control" id="form-next-hop" value="' + nextHop + '">' +
			'</div>' +
			'<div class="form-group">' +
			'<label>路由距离 (1-255)</label>' +
			'<input type="number" class="form-control" id="form-distance" value="' + distance + '" min="1" max="255">' +
			'</div>' +
			'</form>';

		layer.open({
			type: 1,
			title: title,
			area: ['500px', '450px'],
			content: html,
			btn: ['应用', '返回'],
			yes: function() {
				saveRoute(routeData ? routeData.id : null);
			}
		});
	}

	function saveRoute(id) {
		var destIp = document.getElementById('form-dest-ip').value;
		var destMask = document.getElementById('form-dest-mask').value;
		var nextHop = document.getElementById('form-next-hop').value;
		var distance = document.getElementById('form-distance').value;

		if (!destIp) {
			layer.msg('请输入目的 IP 地址');
			return;
		}

		var url = id ? '/api/v1/routes/static/' + id : '/api/v1/routes/static';
		var method = id ? 'PUT' : 'POST';

		fetch(url, {
			method: method,
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({
				dest_ip: destIp,
				dest_mask: destMask,
				next_hop: nextHop,
				distance: parseInt(distance)
			})
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				layer.msg('配置成功');
				layer.closeAll();
				loadStaticRoutes();
			} else {
				layer.msg(data.message || '操作失败');
			}
		});
	}

	function deleteStaticRoute() {
		var checkboxes = document.querySelectorAll('input[name="route-id"]:checked');
		if (checkboxes.length === 0) {
			layer.msg('请选择要删除的路由');
			return;
		}
		layer.confirm('确认删除选中的路由吗？', {
			btn: ['确认', '取消']
		}, function() {
			var ids = [];
		checkboxes.forEach(function(cb) { ids.push(cb.value); });
			Promise.all(ids.map(function(id) {
				return fetch('/api/v1/routes/static/' + id, {method: 'DELETE'});
			})).then(function() {
				layer.msg('删除成功');
				loadStaticRoutes();
			});
		});
	}

	function deleteStaticRouteById(id) {
		layer.confirm('确认删除该路由吗？', {btn: ['确认', '取消']}, function() {
			fetch('/api/v1/routes/static/' + id, {method: 'DELETE'})
				.then(res => res.json())
				.then(data => {
					layer.msg('删除成功');
					loadStaticRoutes();
				});
		});
	}

	document.addEventListener('DOMContentLoaded', function() {
		loadStaticRoutes();

		// 全选
		document.getElementById('select-all').addEventListener('change', function() {
			var checked = this.checked;
			document.querySelectorAll('input[name="route-id"]').forEach(function(cb) {
				cb.checked = checked;
			});
		});
	});
	</script>
	`

	boxContent := template.HTML(actionButtons + staticRouteTable + actionScript)

	routeBox := boxComp.
		WithHeadBorder().
		SetHeader("IPv4 静态路由信息").
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(routeBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "IPv4 静态路由信息",
		Description: "网络 → IP 路由 → IPv4 静态路由信息",
	}, nil
}
