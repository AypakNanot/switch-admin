package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getPortContent 端口管理页面
func getPortContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	portContent := `
	<style>
		.port-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.port-table th, .port-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.port-table th { background-color: #f5f5f5; }
		.status-up { background-color: #5cb85c; color: white; padding: 3px 8px; border-radius: 3px; font-size: 12px; }
		.status-down { background-color: #d9534f; color: white; padding: 3px 8px; border-radius: 3px; font-size: 12px; }
		.btn-group { margin-top: 10px; }
		.modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); }
		.modal-content { background-color: #fff; margin: 5% auto; padding: 20px; border-radius: 5px; width: 50%; }
		.close { float: right; font-size: 28px; font-weight: bold; cursor: pointer; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">端口列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default" onclick="loadPorts()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="port-table">
				<thead>
					<tr>
						<th>端口</th>
						<th>状态</th>
						<th>速度</th>
						<th>双工</th>
						<th>VLAN</th>
						<th>描述</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="port-list">
					<tr>
						<td colspan="7" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<!-- 编辑端口 modal -->
	<div id="editPortModal" class="modal">
		<div class="modal-content">
			<span class="close" onclick="closeModal()">&times;</span>
			<h3>编辑端口配置</h3>
			<div class="form-group">
				<label>端口</label>
				<input type="text" class="form-control" id="edit-port-name" disabled>
			</div>
			<div class="form-group">
				<label>描述</label>
				<input type="text" class="form-control" id="edit-description" placeholder="端口描述">
			</div>
			<div class="form-group">
				<label>速度</label>
				<select class="form-control" id="edit-speed">
					<option value="auto">自动</option>
					<option value="10">10 Mbps</option>
					<option value="100">100 Mbps</option>
					<option value="1000">1000 Mbps</option>
				</select>
			</div>
			<div class="form-group">
				<label>双工模式</label>
				<select class="form-control" id="edit-duplex">
					<option value="auto">自动</option>
					<option value="full">全双工</option>
					<option value="half">半双工</option>
				</select>
			</div>
			<div class="form-group">
				<label>VLAN</label>
				<input type="number" class="form-control" id="edit-vlan" placeholder="VLAN ID">
			</div>
			<div class="form-group">
				<label>
					<input type="checkbox" id="edit-enabled"> 启用端口
				</label>
			</div>
			<button type="button" class="btn btn-primary" onclick="savePortConfig()">保存配置</button>
		</div>
	</div>

	<script>
	var currentPort = null;

	function loadPorts() {
		fetch('/api/v1/network/ports')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('port-list');
				if (data.code === 200 && data.data && data.data.ports) {
					var html = '';
					data.data.ports.forEach(function(port) {
						var statusClass = port.status === 'up' ? 'status-up' : 'status-down';
						var description = port.description || '-';
						html += '<tr>' +
							'<td><strong>' + port.name + '</strong></td>' +
							'<td><span class="' + statusClass + '">' + port.status + '</span></td>' +
							'<td>' + port.speed + '</td>' +
							'<td>' + port.duplex + '</td>' +
							'<td>' + (port.vlan || 1) + '</td>' +
							'<td>' + description + '</td>' +
							'<td>' +
								'<button class="btn btn-sm btn-primary" onclick="editPort(\\'' + port.name + '\\', ' + JSON.stringify(port).replace(/"/g, '&quot;') + ')">编辑</button> ' +
								'<button class="btn btn-sm btn-warning" onclick="resetPort(\\'' + port.name + '\\')">重置</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="restartPort(\\'' + port.name + '\\')">重启</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; color:#999;">暂无端口数据</td></tr>';
				}
			});
	}

	function editPort(portName, port) {
		currentPort = portName;
		document.getElementById('edit-port-name').value = portName;
		document.getElementById('edit-description').value = port.description || '';
		document.getElementById('edit-speed').value = port.speed || 'auto';
		document.getElementById('edit-duplex').value = port.duplex || 'auto';
		document.getElementById('edit-vlan').value = port.vlan || 1;
		document.getElementById('edit-enabled').checked = port.status === 'up';
		document.getElementById('editPortModal').style.display = 'block';
	}

	function closeModal() {
		document.getElementById('editPortModal').style.display = 'none';
		currentPort = null;
	}

	function savePortConfig() {
		if (!currentPort) return;

		var req = {
			description: document.getElementById('edit-description').value,
			speed: document.getElementById('edit-speed').value,
			duplex: document.getElementById('edit-duplex').value,
			vlan: parseInt(document.getElementById('edit-vlan').value),
			enabled: document.getElementById('edit-enabled').checked
		};

		fetch('/api/v1/network/ports/' + encodeURIComponent(currentPort), {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(req)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('端口配置更新成功');
				closeModal();
				loadPorts();
			} else {
				alert(data.message || '更新失败');
			}
		});
	}

	function resetPort(portName) {
		if (confirm('确定要重置端口 "' + portName + '" 的配置吗？')) {
			fetch('/api/v1/network/ports/' + encodeURIComponent(portName) + '/reset', {
				method: 'POST'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('端口配置已重置');
					loadPorts();
				} else {
					alert(data.message || '重置失败');
				}
			});
		}
	}

	function restartPort(portName) {
		if (confirm('确定要重启端口 "' + portName + '" 吗？\\n\\n警告：重启会导致该端口连接短暂中断！')) {
			fetch('/api/v1/network/ports/' + encodeURIComponent(portName) + '/restart', {
				method: 'POST'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('端口已重启');
					loadPorts();
				} else {
					alert(data.message || '重启失败');
				}
			});
		}
	}

	// 点击 modal 外部关闭
	window.onclick = function(event) {
		var modal = document.getElementById('editPortModal');
		if (event.target == modal) {
			closeModal();
		}
	}

	// 页面加载时获取端口列表
	loadPorts();
	</script>
	`

	boxContent := template.HTML(portContent)

	portBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-ethernet"></i> 端口管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(portBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "端口管理",
		Description: "网络 → 端口管理",
	}, nil
}
