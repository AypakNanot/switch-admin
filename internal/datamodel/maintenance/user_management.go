package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getUsersContent 用户管理页面
func getUsersContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	usersContent := `
	<style>
		.user-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.user-table th, .user-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.user-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 200px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.role-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; background-color: #337ab7; color: white; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">创建用户</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>用户名</label>
				<input type="text" class="form-control" id="username" placeholder="3-32 位字母数字下划线">
			</div>
			<div class="form-group">
				<label>密码</label>
				<input type="password" class="form-control" id="password" placeholder="至少 6 位">
			</div>
			<div class="form-group">
				<label>角色</label>
				<select class="form-control" id="role">
					<option value="0">超级管理员 (super-admin)</option>
					<option value="1">管理员 (admin)</option>
					<option value="2">操作员 (operator)</option>
					<option value="3">只读用户 (readonly)</option>
				</select>
			</div>
			<button type="button" class="btn btn-primary" onclick="createUser()">
				<i class="fa fa-plus"></i> 创建用户
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">用户列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="deleteUsers()">
						<i class="fa fa-trash"></i> 删除
					</button>
					<button type="button" class="btn btn-default" onclick="loadUsers()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="user-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>用户名</th>
						<th>角色</th>
						<th>创建时间</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="user-list">
					<tr>
						<td colspan="5" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#user-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function createUser() {
		var username = document.getElementById('username').value;
		var password = document.getElementById('password').value;
		var role = parseInt(document.getElementById('role').value);

		if (!username || !password) {
			alert('请输入用户名和密码');
			return;
		}

		fetch('/api/v1/users', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({ username: username, password: password, role: role })
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('用户创建成功');
				document.getElementById('username').value = '';
				document.getElementById('password').value = '';
				loadUsers();
			} else {
				alert(data.message || '创建失败');
			}
		});
	}

	function loadUsers() {
		fetch('/api/v1/users')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('user-list');
				if (data.code === 200 && data.data && data.data.users) {
					var html = '';
					data.data.users.forEach(function(user) {
						html += '<tr>' +
							'<td><input type="checkbox" value="' + user.username + '"></td>' +
							'<td>' + user.username + '</td>' +
							'<td><span class="role-badge">' + user.role_name + '</span></td>' +
							'<td>' + user.created_at + '</td>' +
							'<td>' +
								'<button class="btn btn-sm btn-primary" onclick="editUser(\\'' + user.username + '\\')">编辑</button> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteUser(\\'' + user.username + '\\')">删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无用户</td></tr>';
				}
			});
	}

	function deleteUser(username) {
		if (confirm('确定要删除用户 "' + username + '" 吗？')) {
			alert('删除功能待实现');
		}
	}

	function deleteUsers() {
		var selected = document.querySelectorAll('#user-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的用户');
			return;
		}
		if (confirm('确定要删除选中的用户吗？')) {
			alert('批量删除功能待实现');
		}
	}

	function editUser(username) {
		alert('编辑用户功能待实现');
	}

	// 页面加载时获取用户列表
	loadUsers();
	</script>
	`

	boxContent := template.HTML(usersContent)

	usersBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-users"></i> 用户管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(usersBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "用户管理",
		Description: "维护 → 用户管理",
	}, nil
}
