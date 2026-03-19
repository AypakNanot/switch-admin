package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getVLANContent VLAN 管理页面
func getVLANContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	vlanContent := `
	<style>
		.vlan-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.vlan-table th, .vlan-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.vlan-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 200px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.status-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; background-color: #5cb85c; color: white; }
		.btn-group { margin-top: 10px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">创建 VLAN</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>VLAN ID</label>
				<input type="number" class="form-control" id="vlan-id" placeholder="1-4094" min="1" max="4094">
			</div>
			<div class="form-group">
				<label>VLAN 名称</label>
				<input type="text" class="form-control" id="vlan-name" placeholder="3-32 位字母数字下划线">
			</div>
			<button type="button" class="btn btn-primary" onclick="createVLAN()">
				<i class="fa fa-plus"></i> 创建 VLAN
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">VLAN 列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="deleteVLANs()">
						<i class="fa fa-trash"></i> 批量删除
					</button>
					<button type="button" class="btn btn-default" onclick="loadVLANs()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="vlan-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>VLAN ID</th>
						<th>名称</th>
						<th>端口</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="vlan-list">
					<tr>
						<td colspan="6" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#vlan-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function createVLAN() {
		var id = parseInt(document.getElementById('vlan-id').value);
		var name = document.getElementById('vlan-name').value;

		if (!id || !name) {
			alert('请输入 VLAN ID 和名称');
			return;
		}

		if (id < 1 || id > 4094) {
			alert('VLAN ID 必须在 1-4094 范围内');
			return;
		}

		fetch('/api/v1/network/vlans', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({ id: id, name: name })
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('VLAN 创建成功');
				document.getElementById('vlan-id').value = '';
				document.getElementById('vlan-name').value = '';
				loadVLANs();
			} else {
				alert(data.message || '创建失败');
			}
		});
	}

	function loadVLANs() {
		fetch('/api/v1/network/vlans')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('vlan-list');
				if (data.code === 200 && data.data && data.data.vlans) {
					var html = '';
					data.data.vlans.forEach(function(vlan) {
						var ports = vlan.ports ? vlan.ports.join(', ') : '无';
						html += '<tr>' +
							'<td><input type="checkbox" value="' + vlan.id + '"></td>' +
							'<td>' + vlan.id + '</td>' +
							'<td>' + vlan.name + '</td>' +
							'<td>' + ports + '</td>' +
							'<td><span class="status-badge">' + vlan.status + '</span></td>' +
							'<td>' +
								'<button class="btn btn-sm btn-primary" onclick="editVLAN(' + vlan.id + ')">编辑</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteVLAN(' + vlan.id + ')">删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="6" style="text-align:center; color:#999;">暂无 VLAN</td></tr>';
				}
			});
	}

	function editVLAN(id) {
		alert('编辑 VLAN 功能待实现，ID: ' + id);
	}

	function deleteVLAN(id) {
		if (confirm('确定要删除 VLAN ' + id + ' 吗？')) {
			fetch('/api/v1/network/vlans/' + id, {
				method: 'DELETE'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('VLAN 删除成功');
					loadVLANs();
				} else {
					alert(data.message || '删除失败');
				}
			});
		}
	}

	function deleteVLANs() {
		var selected = document.querySelectorAll('#vlan-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的 VLAN');
			return;
		}

		var ids = [];
		selected.forEach(cb => ids.push(parseInt(cb.value)));

		if (confirm('确定要删除选中的 ' + ids.length + ' 个 VLAN 吗？')) {
			fetch('/api/v1/network/vlans', {
				method: 'DELETE',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ ids: ids })
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('VLAN 批量删除成功');
					loadVLANs();
				} else {
					alert(data.message || '批量删除失败');
				}
			});
		}
	}

	// 页面加载时获取 VLAN 列表
	loadVLANs();
	</script>
	`

	boxContent := template.HTML(vlanContent)

	vlanBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-network-wired"></i> VLAN 管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(vlanBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "VLAN 管理",
		Description: "网络 → VLAN 管理",
	}, nil
}
