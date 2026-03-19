package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSTPContent STP 管理页面
func getSTPContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	stpContent := `
	<style>
		.stp-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.stp-table th, .stp-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.stp-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 200px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.status-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; }
		.status-forwarding { background-color: #5cb85c; color: white; }
		.status-blocking { background-color: #d9534f; color: white; }
		.status-listening { background-color: #f0ad4e; color: white; }
		.status-learning { background-color: #5bc0de; color: white; }
		.config-section { margin-bottom: 30px; }
		.btn-group { margin-top: 10px; }
	</style>

	<div class="box box-default config-section">
		<div class="box-header with-border">
			<h3 class="box-title">STP 基本配置</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-primary" onclick="saveSTPConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>STP 模式</label>
				<select class="form-control" id="stp-mode">
					<option value="disabled">Disabled</option>
					<option value="stp">STP (802.1D)</option>
					<option value="rstp">RSTP (802.1w)</option>
					<option value="mstp">MSTP (802.1s)</option>
				</select>
			</div>
			<div class="form-group">
				<label>交换机优先级</label>
				<input type="number" class="form-control" id="stp-priority" placeholder="4096 的倍数，默认 32768" value="32768">
				<small>值越小优先级越高，必须为 4096 的倍数</small>
			</div>
			<div class="form-group">
				<label>根桥 ID</label>
				<input type="text" class="form-control" id="stp-root-bridge" disabled value="正在计算...">
			</div>
		</div>
	</div>

	<div class="box box-default config-section">
		<div class="box-header with-border">
			<h3 class="box-title">STP 端口状态</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default" onclick="loadSTPStatus()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="stp-table">
				<thead>
					<tr>
						<th>端口</th>
						<th>角色</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="stp-port-list">
					<tr>
						<td colspan="4" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	function loadSTPConfig() {
		fetch('/api/v1/network/stp/config')
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data) {
					document.getElementById('stp-mode').value = data.data.mode || 'disabled';
					document.getElementById('stp-priority').value = data.data.priority || 32768;
					document.getElementById('stp-root-bridge').value = data.data.root_bridge || '正在计算...';
				}
			});
	}

	function saveSTPConfig() {
		var req = {
			enabled: document.getElementById('stp-mode').value !== 'disabled',
			mode: document.getElementById('stp-mode').value,
			priority: parseInt(document.getElementById('stp-priority').value)
		};

		if (req.priority % 4096 !== 0) {
			alert('交换机优先级必须是 4096 的倍数');
			return;
		}

		fetch('/api/v1/network/stp/config', {
			method: 'PUT',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(req)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('STP 配置保存成功');
				loadSTPConfig();
			} else {
				alert(data.message || '保存失败');
			}
		});
	}

	function loadSTPStatus() {
		fetch('/api/v1/network/stp/status')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('stp-port-list');
				if (data.code === 200 && data.data && data.data.port_states) {
					var html = '';
					data.data.port_states.forEach(function(port) {
						var statusClass = 'status-' + port.state.toLowerCase();
						html += '<tr>' +
							'<td>' + port.port + '</td>' +
							'<td>' + port.role + '</td>' +
							'<td><span class="status-badge ' + statusClass + '">' + port.state + '</span></td>' +
							'<td>' +
								'<button class="btn btn-sm btn-default" onclick="showPortDetail(\\'' + port.port + '\\')">详情</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; color:#999;">暂无 STP 端口数据</td></tr>';
				}
			});
	}

	function showPortDetail(port) {
		alert('端口 ' + port + ' 的 STP 详情功能待实现');
	}

	// 页面加载时加载 STP 配置和状态
	loadSTPConfig();
	loadSTPStatus();
	</script>
	`

	boxContent := template.HTML(stpContent)

	stpBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-random"></i> STP 管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(stpBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "STP 管理",
		Description: "网络 → STP 管理",
	}, nil
}
