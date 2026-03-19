package network

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getACLContent ACL 管理页面
func getACLContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	aclContent := `
	<style>
		.acl-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.acl-table th, .acl-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.acl-table th { background-color: #f5f5f5; }
		.status-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; background-color: #5cb85c; color: white; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 200px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.btn-group { margin-top: 10px; }
		.modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); }
		.modal-content { background-color: #fff; margin: 5% auto; padding: 20px; border-radius: 5px; width: 60%; }
		.close { float: right; font-size: 28px; font-weight: bold; cursor: pointer; }
		.rule-item { border: 1px solid #ddd; padding: 10px; margin-bottom: 10px; border-radius: 4px; background-color: #f9f9f9; }
		.rule-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
		.action-permit { color: #5cb85c; font-weight: bold; }
		.action-deny { color: #d9534f; font-weight: bold; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">创建 ACL</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>ACL 名称</label>
				<input type="text" class="form-control" id="acl-name" placeholder="3-32 位字母数字下划线">
			</div>
			<div class="form-group">
				<label>ACL 类型</label>
				<select class="form-control" id="acl-type">
					<option value="standard">标准 ACL (基于源地址)</option>
					<option value="extended">扩展 ACL (基于源/目的地址、端口、协议)</option>
				</select>
			</div>
			<button type="button" class="btn btn-primary" onclick="createACL()">
				<i class="fa fa-plus"></i> 创建 ACL
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">ACL 列表</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-default" onclick="loadACLs()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="acl-table">
				<thead>
					<tr>
						<th>ID</th>
						<th>名称</th>
						<th>类型</th>
						<th>规则数</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="acl-list">
					<tr>
						<td colspan="6" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<!-- ACL 规则管理 modal -->
	<div id="ruleModal" class="modal">
		<div class="modal-content">
			<span class="close" onclick="closeRuleModal()">&times;</span>
			<h3 id="rule-modal-title">ACL 规则管理</h3>

			<div style="margin-bottom: 20px;">
				<h4>添加规则</h4>
				<div class="form-group">
					<label>动作</label>
					<select class="form-control" id="rule-action" style="width: 150px;">
						<option value="permit">允许</option>
						<option value="deny">拒绝</option>
					</select>
				</div>
				<div class="form-group">
					<label>源地址</label>
					<input type="text" class="form-control" id="rule-source" placeholder="例如：192.168.1.0/24 或 any">
				</div>
				<div class="form-group">
					<label>目的地址</label>
					<input type="text" class="form-control" id="rule-dest" placeholder="例如：10.0.0.0/8 或 any">
				</div>
				<div class="form-group">
					<label>协议</label>
					<select class="form-control" id="rule-protocol" style="width: 150px;">
						<option value="ip">IP</option>
						<option value="tcp">TCP</option>
						<option value="udp">UDP</option>
						<option value="icmp">ICMP</option>
					</select>
				</div>
				<div class="form-group">
					<label>端口</label>
					<input type="text" class="form-control" id="rule-port" placeholder="例如：80,443 或 8000-9000">
				</div>
				<button type="button" class="btn btn-primary" onclick="addRule()">添加规则</button>
			</div>

			<div>
				<h4>现有规则</h4>
				<div id="rule-list">
					<!-- 规则列表将动态加载 -->
				</div>
			</div>
		</div>
	</div>

	<script>
	var currentACLID = null;

	function createACL() {
		var name = document.getElementById('acl-name').value;
		var type = document.getElementById('acl-type').value;

		if (!name) {
			alert('请输入 ACL 名称');
			return;
		}

		fetch('/api/v1/network/acls', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({ name: name, type: type })
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('ACL 创建成功');
				document.getElementById('acl-name').value = '';
				loadACLs();
			} else {
				alert(data.message || '创建失败');
			}
		});
	}

	function loadACLs() {
		fetch('/api/v1/network/acls')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('acl-list');
				if (data.code === 200 && data.data && data.data.acls) {
					var html = '';
					data.data.acls.forEach(function(acl) {
						html += '<tr>' +
							'<td>' + acl.id + '</td>' +
							'<td>' + acl.name + '</td>' +
							'<td>' + acl.type + '</td>' +
							'<td>' + acl.rules + '</td>' +
							'<td><span class="status-badge">' + acl.status + '</span></td>' +
							'<td>' +
								'<button class="btn btn-sm btn-primary" onclick="editACL(' + acl.id + ')">编辑</button> ' +
								'<button class="btn btn-sm btn-info" onclick="showRules(' + acl.id + ', \\'' + acl.name + '\\')">规则</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteACL(' + acl.id + ')">删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="6" style="text-align:center; color:#999;">暂无 ACL</td></tr>';
				}
			});
	}

	function editACL(id) {
		alert('编辑 ACL 功能待实现，ID: ' + id);
	}

	function deleteACL(id) {
		if (confirm('确定要删除 ACL ' + id + ' 吗？')) {
			fetch('/api/v1/network/acls/' + id, {
				method: 'DELETE'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('ACL 删除成功');
					loadACLs();
				} else {
					alert(data.message || '删除失败');
				}
			});
		}
	}

	function showRules(aclID, aclName) {
		currentACLID = aclID;
		document.getElementById('rule-modal-title').innerText = 'ACL 规则管理 - ' + aclName;
		document.getElementById('ruleModal').style.display = 'block';
		loadRules(aclID);
	}

	function closeRuleModal() {
		document.getElementById('ruleModal').style.display = 'none';
		currentACLID = null;
	}

	function loadRules(aclID) {
		fetch('/api/v1/network/acls/' + aclID + '/rules')
			.then(res => res.json())
			.then(data => {
				var ruleList = document.getElementById('rule-list');
				if (data.code === 200 && data.data && data.data.rules) {
					var html = '';
					data.data.rules.forEach(function(rule, index) {
						var actionClass = rule.action === 'permit' ? 'action-permit' : 'action-deny';
						html += '<div class="rule-item">' +
							'<div class="rule-header">' +
								'<strong>规则 ' + (index + 1) + '</strong>' +
								'<span class="' + actionClass + '">' + rule.action.toUpperCase() + '</span>' +
							'</div>' +
							'<div>源地址：' + (rule.source || 'any') + '</div>' +
							'<div>目的地址：' + (rule.destination || 'any') + '</div>' +
							'<div>协议：' + (rule.protocol || 'IP') + '</div>' +
							'<div>端口：' + (rule.port || 'any') + '</div>' +
							'<button class="btn btn-sm btn-danger" style="margin-top:10px;" onclick="deleteRule(' + aclID + ', ' + rule.id + ')">删除</button>' +
							'</div>';
					});
					ruleList.innerHTML = html || '<p style="color:#999;">暂无规则</p>';
				} else {
					ruleList.innerHTML = '<p style="color:#999;">暂无规则</p>';
				}
			});
	}

	function addRule() {
		if (!currentACLID) return;

		var req = {
			action: document.getElementById('rule-action').value,
			source: document.getElementById('rule-source').value || 'any',
			destination: document.getElementById('rule-dest').value || 'any',
			protocol: document.getElementById('rule-protocol').value,
			port: document.getElementById('rule-port').value || 'any'
		};

		fetch('/api/v1/network/acls/' + currentACLID + '/rules', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify(req)
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('规则添加成功');
				loadRules(currentACLID);
				// 清空表单
				document.getElementById('rule-source').value = '';
				document.getElementById('rule-dest').value = '';
				document.getElementById('rule-port').value = '';
			} else {
				alert(data.message || '添加失败');
			}
		});
	}

	function deleteRule(aclID, ruleID) {
		if (confirm('确定要删除此规则吗？')) {
			fetch('/api/v1/network/acls/' + aclID + '/rules/' + ruleID, {
				method: 'DELETE'
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('规则删除成功');
					loadRules(aclID);
				} else {
					alert(data.message || '删除失败');
				}
			});
		}
	}

	// 点击 modal 外部关闭
	window.onclick = function(event) {
		var modal = document.getElementById('ruleModal');
		if (event.target == modal) {
			closeRuleModal();
		}
	}

	// 页面加载时获取 ACL 列表
	loadACLs();
	</script>
	`

	boxContent := template.HTML(aclContent)

	aclBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-shield-alt"></i> ACL 管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(aclBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "ACL 管理",
		Description: "网络 → ACL 管理",
	}, nil
}
