package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getLAGContent 链路聚合管理页面
func getLAGContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	lagContent := `
	<style>
		.lag-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.lag-table th, .lag-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.lag-table th { background-color: #f5f5f5; }
		.status-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; background-color: #5cb85c; color: white; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 200px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.btn-group { margin-top: 10px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">创建链路聚合组</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>LAG 名称</label>
				<input type="text" class="form-control" id="lag-name" placeholder="例如：lag1">
			</div>
			<div class="form-group">
				<label>聚合模式</label>
				<select class="form-control" id="lag-mode">
					<option value="lacp">LACP (802.3ad)</option>
					<option value="static">静态聚合</option>
				</select>
			</div>
			<div class="form-group">
				<label>成员端口</label>
				<input type="text" class="form-control" id="lag-ports" placeholder="例如：eth0/1,eth0/2">
			</div>
			<button type="button" class="btn btn-primary" onclick="createLAG()">
				<i class="fa fa-plus"></i> 创建聚合组
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">链路聚合组列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default" onclick="loadLAGs()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="lag-table">
				<thead>
					<tr>
						<th>ID</th>
						<th>名称</th>
						<th>模式</th>
						<th>成员端口</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="lag-list">
					<tr>
						<td colspan="6" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	function createLAG() {
		var name = document.getElementById('lag-name').value;
		var mode = document.getElementById('lag-mode').value;
		var portsStr = document.getElementById('lag-ports').value;

		if (!name || !portsStr) {
			alert('请输入 LAG 名称和成员端口');
			return;
		}

		var ports = portsStr.split(',').map(p => p.trim()).filter(p => p);

		fetch('/api/v1/network/lags', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({ name: name, mode: mode, ports: ports })
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('链路聚合组创建成功');
				document.getElementById('lag-name').value = '';
				document.getElementById('lag-ports').value = '';
				loadLAGs();
			} else {
				alert(data.message || '创建失败');
			}
		});
	}

	function loadLAGs() {
		fetch('/api/v1/network/lags')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('lag-list');
				if (data.code === 200 && data.data && data.data.lags) {
					var html = '';
					data.data.lags.forEach(function(lag) {
						var ports = lag.ports ? lag.ports.join(', ') : '无';
						html += '<tr>' +
							'<td>' + lag.id + '</td>' +
							'<td>' + lag.name + '</td>' +
							'<td>' + lag.mode + '</td>' +
							'<td>' + ports + '</td>' +
							'<td><span class="status-badge">' + lag.status + '</span></td>' +
							'<td>' +
								'<button class="btn btn-sm btn-primary" onclick="editLAG(' + lag.id + ')">编辑</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteLAG(' + lag.id + ')">删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="6" style="text-align:center; color:#999;">暂无链路聚合组</td></tr>';
				}
			});
	}

	function editLAG(id) {
		alert('编辑链路聚合组功能待实现，ID: ' + id);
	}

	function deleteLAG(id) {
		if (confirm('确定要删除链路聚合组 ' + id + ' 吗？')) {
			fetch('/api/v1/network/lags/' + id, {
				method: 'DELETE'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('链路聚合组删除成功');
					loadLAGs();
				} else {
					alert(data.message || '删除失败');
				}
			});
		}
	}

	// 页面加载时获取链路聚合组列表
	loadLAGs();
	</script>
	`

	boxContent := template.HTML(lagContent)

	lagBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-link"></i> 链路聚合管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(lagBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "链路聚合管理",
		Description: "网络 → 链路聚合管理",
	}, nil
}
