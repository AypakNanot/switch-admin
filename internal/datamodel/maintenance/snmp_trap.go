package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSNMPTrapContent SNMP Trap 配置页面
func getSNMPTrapContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	snmpTrapContent := `
	<style>
		.trap-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.trap-table th, .trap-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.trap-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 300px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.checkbox-group label { display: inline-block; margin-right: 15px; margin-bottom: 10px; font-weight: normal; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">Trap 使能配置</h3>
		</div>
		<div class="box-body">
			<div class="checkbox-group">
				<label><input type="checkbox" name="trap_enable" value="Coldstart"> Coldstart Trap</label>
				<label><input type="checkbox" name="trap_enable" value="Warmstart"> Warmstart Trap</label>
				<label><input type="checkbox" name="trap_enable" value="Linkup"> Linkup Trap</label>
				<label><input type="checkbox" name="trap_enable" value="Linkdown"> Linkdown Trap</label>
				<label><input type="checkbox" name="trap_enable" value="System"> System Trap</label>
				<label><input type="checkbox" name="trap_enable" value="Loopback-detect"> Loopback-detect</label>
			</div>
			<button type="button" class="btn btn-primary" onclick="saveTrapConfig()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">Trap 目标主机配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>目标地址</label>
				<input type="text" class="form-control" id="trap-target" placeholder="IPv4 或 IPv6 地址，例如 1.1.1.1">
			</div>
			<div class="form-group">
				<label>团体名称</label>
				<input type="text" class="form-control" id="trap-community" placeholder="不包含 ? > \" \\ 字符，长度为 1-128">
			</div>
			<div class="form-group">
				<label>UDP 端口</label>
				<input type="number" class="form-control" id="trap-port" value="162" min="0" max="65535">
			</div>
			<div class="form-group">
				<label>VRF ID</label>
				<select class="form-control" id="trap-vrf">
					<option value="mgmt-if">mgmt-if</option>
					<option value="default">default</option>
				</select>
			</div>
			<button type="button" class="btn btn-primary" onclick="addTrapHost()">
				<i class="fa fa-plus"></i> 新建
			</button>

			<table class="trap-table" style="margin-top: 20px;">
				<thead>
					<tr>
						<th>目标地址</th>
						<th>UDP 端口</th>
						<th>Vrf 名称</th>
						<th>团体名称</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="trap-host-list">
					<tr>
						<td colspan="5" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	function saveTrapConfig() {
		alert('Trap 配置保存功能待实现');
	}

	function addTrapHost() {
		alert('添加 Trap 目标主机功能待实现');
	}

	function loadTrapHosts() {
		fetch('/api/v1/snmp/trap/hosts')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('trap-host-list');
				if (data.code === 200 && data.data && data.data.hosts) {
					var html = '';
					data.data.hosts.forEach(function(host) {
						html += '<tr>' +
							'<td>' + host.target + '</td>' +
							'<td>' + host.port + '</td>' +
							'<td>' + host.vrf + '</td>' +
							'<td>' + host.community + '</td>' +
							'<td>' +
								'<button class="btn btn-sm btn-default" onclick="testTrap(\\'' + host.id + '\\')">测试</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteTrapHost(\\'' + host.id + '\\')">删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无配置</td></tr>';
				}
			});
	}

	function testTrap(id) {
		if (confirm('确定要发送测试 Trap 吗？')) {
			fetch('/api/v1/snmp/trap/hosts/' + id + '/test', { method: 'POST' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('测试 Trap 已发送，请检查网管平台是否收到告警');
					} else {
						alert(data.message || '发送失败');
					}
				});
		}
	}

	function deleteTrapHost(id) {
		if (confirm('确定要删除此 Trap 目标主机吗？')) {
			alert('删除功能待实现');
		}
	}

	// 页面加载时获取 Trap 目标主机列表
	loadTrapHosts();
	</script>
	`

	boxContent := template.HTML(snmpTrapContent)

	trapBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-exclamation-triangle"></i> SNMP Trap 配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(trapBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "SNMP Trap 配置",
		Description: "维护 → SNMP Trap 配置",
	}, nil
}
