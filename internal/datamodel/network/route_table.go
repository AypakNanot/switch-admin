package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetRouteTableContent IPv4 路由表页面
func GetRouteTableContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	// 搜索表单
	searchForm := `
	<form id="route-search-form" class="form-inline" role="form">
		<div class="form-group">
			<label>目的 IP</label>
			<input type="text" class="form-control" id="dest_ip" placeholder="支持模糊匹配，如 10.10.">
		</div>
		<div class="form-group">
			<label>协议类型</label>
			<select class="form-control" id="protocol">
				<option value="">全部</option>
				<option value="Static">Static(静态)</option>
				<option value="OSPF">OSPF</option>
				<option value="RIP">RIP</option>
				<option value="BGP">BGP</option>
				<option value="Connected">Connected(直连)</option>
				<option value="Local">Local(本地)</option>
			</select>
		</div>
		<button type="button" class="btn btn-primary" onclick="searchRoutes()">
			<i class="fa fa-search"></i> 查询
		</button>
	</form>

	<script>
	function searchRoutes() {
		var destIp = document.getElementById('dest_ip').value;
		var protocol = document.getElementById('protocol').value;
		var url = '/api/v1/routes/table?page=1&page_size=50';
		if (destIp) url += '&dest_ip=' + encodeURIComponent(destIp);
		if (protocol) url += '&protocol=' + encodeURIComponent(protocol);

		fetch(url)
			.then(res => res.json())
			.then(data => {
				renderRouteTable(data.data);
			});
	}

	function renderRouteTable(data) {
		var tbody = document.getElementById('route-table-body');
		var html = '';
		if (!data.items || data.items.length === 0) {
			tbody.innerHTML = '<tr><td colspan="6" style="text-align:center;">暂无数据~</td></tr>';
			return;
		}
		data.items.forEach(function(item) {
			var protocolClass = 'label-default';
			if (item.protocol === 'Static') protocolClass = 'label-primary';
			else if (item.protocol === 'OSPF') protocolClass = 'label-info';
			else if (item.protocol === 'BGP') protocolClass = 'label-warning';
			else if (item.protocol === 'Connected') protocolClass = 'label-success';

			html += '<tr>' +
				'<td>' + item.dest_ip + '</td>' +
				'<td>' + item.dest_mask + '</td>' +
				'<td><span class="label ' + protocolClass + '">' + item.protocol + '</span></td>' +
				'<td>' + item.next_hop + '</td>' +
				'<td>' + item.out_port + '</td>' +
				'<td>' + item.metric + '</td>' +
				'</tr>';
		});
		tbody.innerHTML = html;

		// 更新分页信息
		document.getElementById('total').textContent = data.total;
		document.getElementById('page').textContent = data.page;
		document.getElementById('total_pages').textContent = data.total_pages;
	}
	</script>
	`

	// 路由表格
	routeTable := `
	<table class="table table-bordered table-hover">
		<thead>
			<tr>
				<th>目的 IP 地址</th>
				<th>目的 IP 掩码</th>
				<th>协议类型</th>
				<th>下一跳</th>
				<th>出端口</th>
				<th>度量值</th>
			</tr>
		</thead>
		<tbody id="route-table-body">
			<tr><td colspan="6" style="text-align:center;">暂无数据~</td></tr>
		</tbody>
	</table>
	<div class="pull-right" style="margin-top:10px;">
		<span class="text-muted">
			总计：<span id="total">0</span> 条 |
			第 <span id="page">1</span> 页 / 共 <span id="total_pages">0</span> 页
		</span>
	</div>
	`

	// 页面加载时自动查询
	autoLoadScript := `
	<script>
	document.addEventListener('DOMContentLoaded', function() {
		searchRoutes();
	});
	</script>
	`

	boxContent := template.HTML(searchForm + routeTable + autoLoadScript)

	routeBox := boxComp.
		WithHeadBorder().
		SetHeader("IPv4 路由表").
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(routeBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "IPv4 路由表",
		Description: "网络 → IP 路由 → IPv4 路由表",
	}, nil
}
