package datamodel

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetRebootSaveContent 重启/保存页面
func GetRebootSaveContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	rebootSaveContent := `
	<style>
		.maintenance-section { margin-bottom: 30px; }
		.maintenance-section h4 { margin-bottom: 15px; color: #333; }
		.btn-danger { background-color: #d9534f; border-color: #d43f3a; }
		.btn-danger:hover { background-color: #c9302c; border-color: #ac2925; }
		.warning-text { color: #f0ad4e; font-weight: bold; }
	</style>

	<div class="maintenance-section">
		<h4><i class="fa fa-floppy-o"></i> 保存配置</h4>
		<div class="box box-default">
			<div class="box-body">
				<p>将当前运行配置写入启动配置文件，重启后生效。</p>
				<button type="button" class="btn btn-primary" onclick="saveConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
	</div>

	<div class="maintenance-section">
		<h4><i class="fa fa-refresh"></i> 重启交换机</h4>
		<div class="box box-default">
			<div class="box-body">
				<p class="warning-text">警告：重启期间网络将中断约 1-3 分钟。</p>
				<div class="checkbox">
					<label>
						<input type="checkbox" id="save-before-reboot" checked>
						重启之前自动保存配置
					</label>
				</div>
				<button type="button" class="btn btn-warning" onclick="rebootSwitch()">
					<i class="fa fa-power-off"></i> 重启交换机
				</button>
			</div>
		</div>
	</div>

	<div class="maintenance-section">
		<h4><i class="fa fa-exclamation-triangle"></i> 恢复出厂配置</h4>
		<div class="box box-default">
			<div class="box-body">
				<p class="warning-text">⚠️ 警告：此操作将清除所有配置并恢复出厂默认状态！</p>
				<p>当前配置将永久丢失，且无法恢复。</p>
				<button type="button" class="btn btn-danger" onclick="factoryReset()">
					<i class="fa fa-trash"></i> 恢复出厂配置
				</button>
			</div>
		</div>
	</div>

	<script>
	function saveConfig() {
		if (confirm('确定要保存配置吗？')) {
			fetch('/api/v1/system/save-config', { method: 'POST' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('配置保存成功');
					} else {
						alert(data.message || '保存失败');
					}
				});
		}
	}

	function rebootSwitch() {
		var saveBefore = document.getElementById('save-before-reboot').checked;
		if (confirm('确定要重启交换机吗？重启期间网络将中断约 1-3 分钟。')) {
			fetch('/api/v1/system/reboot', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ save_before_reboot: saveBefore })
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('设备正在重启，请稍后刷新页面...');
				} else {
					alert(data.message || '重启失败');
				}
			});
		}
	}

	function factoryReset() {
		var confirmation = prompt('⚠️ 警告：此操作将清除所有配置并恢复出厂默认状态！\\n\\n当前配置将永久丢失，且无法恢复。\\n重启后，设备管理 IP 将恢复为默认地址：192.168.1.1\\n\\n请输入 CONFIRM 以确认此操作：');
		if (confirmation === 'CONFIRM') {
			fetch('/api/v1/system/factory-reset', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ confirmation: 'CONFIRM' })
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('恢复出厂配置成功，设备正在重启...');
				} else {
					alert(data.message || '恢复失败');
				}
			});
		} else if (confirmation !== null) {
			alert('确认失败，请输入 CONFIRM');
		}
	}
	</script>
	`

	boxContent := template.HTML(rebootSaveContent)

	rebootBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-power-off"></i> 重启/保存`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(rebootBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "重启/保存",
		Description: "维护 → 重启/保存",
	}, nil
}

// GetUsersContent 用户管理页面
func GetUsersContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	usersContent := `
	<style>
		.user-table { width: 100%; border-collapse: collapse; }
		.user-table th, .user-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.user-table th { background-color: #f5f5f5; }
		.btn-group { display: flex; gap: 5px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">用户列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-primary" onclick="showAddUserModal()">
						<i class="fa fa-plus"></i> 添加
					</button>
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
						<th>权限等级</th>
						<th>密码</th>
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
	var permissionMap = {
		1: 'admin',
		2: 'operator',
		3: 'viewer'
	};

	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#user-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function loadUsers() {
		fetch('/api/v1/users')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('user-list');
				if (data.code === 200 && data.data && data.data.list) {
					var html = '';
					data.data.list.forEach(function(user) {
						html += '<tr>' +
							'<td><input type="checkbox" value="' + user.username + '"></td>' +
							'<td>' + user.username + '</td>' +
							'<td>' + (permissionMap[user.role] || user.role) + '</td>' +
							'<td>******</td>' +
							'<td><button class="btn btn-sm btn-primary" onclick="editUser(\\'' + user.username + '\\')">编辑</button></td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无数据</td></tr>';
				}
			});
	}

	function showAddUserModal() {
		alert('添加用户功能待实现');
	}

	function editUser(username) {
		alert('编辑用户：' + username);
	}

	function deleteUsers() {
		var selected = document.querySelectorAll('#user-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的用户');
			return;
		}
		if (confirm('确定要删除选中的用户吗？')) {
			alert('删除功能待实现');
		}
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

// GetMaintenanceSystemConfigContent 系统配置页面（维护模块下）
func GetMaintenanceSystemConfigContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	systemConfigContent := `
	<style>
		.config-section { margin-bottom: 30px; }
		.config-section h4 { margin-bottom: 15px; color: #333; }
		.form-inline { display: inline-block; margin-right: 15px; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: inline-block; width: 80px; }
		.form-control { display: inline-block; width: 200px; }
	</style>

	<div class="config-section">
		<h4><i class="fa fa-network-wired"></i> 基础设置</h4>
		<div class="box box-default">
			<div class="box-body">
				<div class="form-group">
					<label>管理地址</label>
					<input type="text" class="form-control" placeholder="192.168.1.1" value="192.168.1.10">
				</div>
				<div class="form-group">
					<label>掩码</label>
					<select class="form-control">
						<option value="255.255.255.0">255.255.255.0 (24)</option>
						<option value="255.255.0.0">255.255.0.0 (16)</option>
						<option value="255.0.0.0">255.0.0.0 (8)</option>
					</select>
				</div>
				<div class="form-group">
					<label>默认网关</label>
					<input type="text" class="form-control" placeholder="192.168.1.254" value="192.168.1.254">
				</div>
				<button type="button" class="btn btn-primary" onclick="saveNetworkConfig()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<div class="config-section">
		<h4><i class="fa fa-thermometer-half"></i> 温度阈值</h4>
		<div class="box box-default">
			<div class="box-body">
				<div class="form-group">
					<label>温度阈值</label>
					<input type="text" class="form-control" value="5/65/80" placeholder="低温/预警/告警">
				</div>
				<p class="text-muted">格式：低温告警阈值/高温预警阈值/高温告警阈值 (默认值：5/65/80，单位：摄氏度)</p>
				<button type="button" class="btn btn-primary" onclick="saveTemperatureConfig()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<div class="config-section">
		<h4><i class="fa fa-info-circle"></i> 基本信息</h4>
		<div class="box box-default">
			<div class="box-body">
				<div class="form-group">
					<label>设备名</label>
					<input type="text" class="form-control" placeholder="Switch-01" value="Switch-01">
				</div>
				<div class="form-group">
					<label>联系</label>
					<input type="text" class="form-control" placeholder="管理员联系方式" value="admin@example.com">
				</div>
				<div class="form-group">
					<label>位置</label>
					<input type="text" class="form-control" placeholder="设备物理位置" value="机房 A-01">
				</div>
				<button type="button" class="btn btn-primary" onclick="saveDeviceInfo()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<div class="config-section">
		<h4><i class="fa fa-clock-o"></i> 时间日期</h4>
		<div class="box box-default">
			<div class="box-body">
				<div class="form-group">
					<label>时间日期</label>
					<input type="datetime-local" class="form-control" id="datetime-input">
				</div>
				<button type="button" class="btn btn-primary" onclick="saveDateTime()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<script>
	// 设置当前时间
	document.getElementById('datetime-input').value = new Date().toISOString().slice(0, 16);

	function saveNetworkConfig() {
		alert('网络配置保存功能待实现');
	}

	function saveTemperatureConfig() {
		alert('温度阈值配置保存功能待实现');
	}

	function saveDeviceInfo() {
		alert('设备信息保存功能待实现');
	}

	function saveDateTime() {
		alert('时间日期保存功能待实现');
	}
	</script>
	`

	boxContent := template.HTML(systemConfigContent)

	configBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-cogs"></i> 系统配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(configBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "系统配置",
		Description: "维护 → 系统配置",
	}, nil
}

// GetLoadConfigContent 加载配置页面
func GetLoadConfigContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	loadConfigContent := `
	<style>
		.config-table { width: 100%; border-collapse: collapse; }
		.config-table th, .config-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.config-table th { background-color: #f5f5f5; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">配置文件列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-primary" onclick="loadSelectedConfig()">
						<i class="fa fa-upload"></i> 加载
					</button>
					<button type="button" class="btn btn-default" onclick="loadConfigFiles()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="config-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>文件名</th>
						<th>时间</th>
						<th>大小</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="config-file-list">
					<tr>
						<td colspan="5" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#config-file-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function loadConfigFiles() {
		fetch('/api/v1/config/files')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('config-file-list');
				if (data.code === 200 && data.data && data.data.files) {
					var html = '';
					data.data.files.forEach(function(file) {
						html += '<tr>' +
							'<td><input type="checkbox" value="' + file.path + '"></td>' +
							'<td>' + file.name + '</td>' +
							'<td>' + file.created_at + '</td>' +
							'<td>' + file.size + '</td>' +
							'<td><button class="btn btn-sm btn-primary" onclick="loadConfig(\\'' + file.path + '\\')">加载</button></td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无配置文件</td></tr>';
				}
			});
	}

	function loadConfig(filePath) {
		if (confirm('确定要加载此配置文件吗？\\n当前未保存的配置将被覆盖。')) {
			fetch('/api/v1/config/load', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ file_path: filePath })
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert('配置加载成功，部分配置可能需要重启生效');
				} else {
					alert(data.message || '加载失败');
				}
			});
		}
	}

	function loadSelectedConfig() {
		var selected = document.querySelectorAll('#config-file-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择配置文件');
			return;
		}
		loadConfig(selected[0].value);
	}

	// 页面加载时获取配置文件列表
	loadConfigFiles();
	</script>
	`

	boxContent := template.HTML(loadConfigContent)

	configBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-upload"></i> 加载配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(configBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "加载配置",
		Description: "维护 → 加载配置",
	}, nil
}

// GetFilesContent 文件管理页面
func GetFilesContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	filesContent := `
	<style>
		.file-table { width: 100%; border-collapse: collapse; }
		.file-table th, .file-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.file-table th { background-color: #f5f5f5; }
		.memory-info { background: #f5f5f5; padding: 15px; margin-bottom: 20px; border-radius: 4px; }
	</style>

	<div class="memory-info">
		<strong>flash:</strong> 总大小 3.9G, 剩余空间：3.6G<br>
		<strong>flash:/boot</strong> 总大小 2.9G, 剩余空间：2.3G
	</div>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">文件列表</h3>
			<div class="box-tools">
				<div style="margin-bottom: 10px;">
					<input type="file" id="file-input" style="display: inline-block; width: 300px;">
					<button type="button" class="btn btn-primary" onclick="uploadFile()">
						<i class="fa fa-upload"></i> 上传文件
					</button>
					<button type="button" class="btn btn-warning" onclick="uploadFirmware()">
						<i class="fa fa-download"></i> 上传映像
					</button>
				</div>
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="deleteFiles()">
						<i class="fa fa-trash"></i> 删除选中的文件
					</button>
					<button type="button" class="btn btn-default" onclick="loadFiles()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="file-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>文件名</th>
						<th>目录</th>
						<th>大小</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="file-list">
					<tr>
						<td colspan="5" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#file-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function loadFiles() {
		fetch('/api/v1/files')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('file-list');
				if (data.code === 200 && data.data && data.data.files) {
					var html = '';
					data.data.files.forEach(function(file) {
						html += '<tr>' +
							'<td><input type="checkbox" value="' + file.path + '"></td>' +
							'<td>' + file.name + '</td>' +
							'<td>' + file.directory + '</td>' +
							'<td>' + file.size + '</td>' +
							'<td><button class="btn btn-sm btn-primary" onclick="downloadFile(\\'' + file.path + '\\')">下载</button></td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无文件</td></tr>';
				}
			});
	}

	function uploadFile() {
		alert('文件上传功能待实现');
	}

	function uploadFirmware() {
		alert('固件上传功能待实现');
	}

	function downloadFile(filePath) {
		window.location.href = '/api/v1/files/' + encodeURIComponent(filePath);
	}

	function deleteFiles() {
		var selected = document.querySelectorAll('#file-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的文件');
			return;
		}
		if (confirm('确定要删除选中的文件吗？')) {
			alert('删除功能待实现');
		}
	}

	// 页面加载时获取文件列表
	loadFiles();
	</script>
	`

	boxContent := template.HTML(filesContent)

	filesBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-file"></i> 文件管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(filesBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "文件管理",
		Description: "维护 → 文件管理",
	}, nil
}

// GetLogsContent 日志管理页面
func GetLogsContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	logsContent := `
	<style>
		.log-table { width: 100%; border-collapse: collapse; font-size: 12px; }
		.log-table th, .log-table td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
		.log-table th { background-color: #f5f5f5; }
		.filter-row { margin-bottom: 15px; }
		.filter-row .form-group { display: inline-block; margin-right: 15px; }
		.log-content { max-height: 500px; overflow-y: auto; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">系统日志</h3>
			<div class="box-tools">
				<div class="filter-row">
					<div class="form-group">
						<label>时间范围</label>
						<input type="datetime-local" class="form-control" id="start-time" style="width: 180px;">
					</div>
					<span style="margin: 0 10px;">~</span>
					<input type="datetime-local" class="form-control" id="end-time" style="width: 180px; display: inline-block;">
				</div>
				<div class="form-group">
					<label>级别</label>
					<select class="form-control" id="log-level" style="width: 120px; display: inline-block;">
						<option value="All">All</option>
						<option value="Emergency">Emergency</option>
						<option value="Alert">Alert</option>
						<option value="Critical">Critical</option>
						<option value="Error">Error</option>
						<option value="Warning">Warning</option>
						<option value="Notice">Notice</option>
						<option value="Info">Info</option>
						<option value="Debug">Debug</option>
					</select>
				</div>
				<div class="form-group">
					<label>模块</label>
					<select class="form-control" id="log-module" style="width: 120px; display: inline-block;">
						<option value="All">All</option>
						<option value="DHCLIENT">DHCLIENT</option>
						<option value="IMI">IMI</option>
						<option value="SNMP">SNMP</option>
						<option value="SYS">SYS</option>
					</select>
				</div>
				<button type="button" class="btn btn-primary" onclick="queryLogs()">
					<i class="fa fa-search"></i> 查询
				</button>
				<button type="button" class="btn btn-default" onclick="loadLogs()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
				<button type="button" class="btn btn-danger" onclick="clearLogs()">
					<i class="fa fa-trash"></i> 清除
				</button>
			</div>
		</div>
		<div class="box-body">
			<div class="log-content">
				<table class="log-table">
					<thead>
						<tr>
							<th>时间</th>
							<th>模块</th>
							<th>级别</th>
							<th>内容</th>
						</tr>
					</thead>
					<tbody id="log-list">
						<tr>
							<td colspan="4" style="text-align:center; color:#999;">加载中...</td>
						</tr>
					</tbody>
				</table>
			</div>
			<div style="margin-top: 10px; color: #999;">Total <span id="log-count">0</span> records.</div>
		</div>
	</div>

	<script>
	// 设置默认时间范围（最近 24 小时）
	var now = new Date();
	var yesterday = new Date(now.getTime() - 24*60*60*1000);
	document.getElementById('start-time').value = yesterday.toISOString().slice(0, 16);
	document.getElementById('end-time').value = now.toISOString().slice(0, 16);

	function loadLogs() {
		queryLogs();
	}

	function queryLogs() {
		var startTime = document.getElementById('start-time').value;
		var endTime = document.getElementById('end-time').value;
		var level = document.getElementById('log-level').value;
		var module = document.getElementById('log-module').value;

		var url = '/api/v1/logs?page=1&page_size=50';
		if (startTime) url += '&start_time=' + encodeURIComponent(startTime);
		if (endTime) url += '&end_time=' + encodeURIComponent(endTime);
		if (level && level !== 'All') url += '&level=' + encodeURIComponent(level);
		if (module && module !== 'All') url += '&module=' + encodeURIComponent(module);

		fetch(url)
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('log-list');
				if (data.code === 200 && data.data && data.data.logs) {
					var html = '';
					data.data.logs.forEach(function(log) {
						html += '<tr>' +
							'<td>' + log.timestamp + '</td>' +
							'<td>' + log.module + '</td>' +
							'<td>' + (log.level || '') + '</td>' +
							'<td>' + log.message + '</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
					document.getElementById('log-count').textContent = data.data.total || 0;
				} else {
					tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; color:#999;">暂无日志</td></tr>';
					document.getElementById('log-count').textContent = '0';
				}
			});
	}

	function clearLogs() {
		if (confirm('确定要清除所有日志吗？此操作不可恢复。')) {
			fetch('/api/v1/logs', { method: 'DELETE' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('日志已清除');
						loadLogs();
					} else {
						alert(data.message || '清除失败');
					}
				});
		}
	}

	// 页面加载时获取日志
	loadLogs();
	</script>
	`

	boxContent := template.HTML(logsContent)

	logsBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-history"></i> 日志管理`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(logsBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "日志管理",
		Description: "维护 → 日志管理",
	}, nil
}

// GetSNMPContent SNMP 配置页面
func GetSNMPContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	snmpContent := `
	<style>
		.snmp-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.snmp-table th, .snmp-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.snmp-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 300px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		.radio-group label { display: inline-block; margin-right: 20px; font-weight: normal; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">基本配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>SNMP 状态</label>
				<div class="radio-group">
					<label><input type="radio" name="snmp_status" value="1"> 启用</label>
					<label><input type="radio" name="snmp_status" value="0" checked> 关闭</label>
				</div>
			</div>
			<div class="form-group">
				<label>SNMP 版本</label>
				<select class="form-control" id="snmp-version">
					<option value="all">All</option>
					<option value="v1">v1</option>
					<option value="v2c">v2c</option>
					<option value="v3">v3</option>
				</select>
			</div>
			<button type="button" class="btn btn-primary" onclick="saveSNMPConfig()">
				<i class="fa fa-save"></i> 应用
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">团体配置</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>团体名称</label>
				<input type="text" class="form-control" id="community-name" placeholder="不包含 ? > \" \\ 字符，长度为 1-128">
			</div>
			<div class="form-group">
				<label>访问方式</label>
				<select class="form-control" id="community-access">
					<option value="Read-Only">Read-Only</option>
					<option value="Read-Write">Read-Write</option>
				</select>
			</div>
			<button type="button" class="btn btn-primary" onclick="addCommunity()">
				<i class="fa fa-plus"></i> 新建
			</button>

			<table class="snmp-table" style="margin-top: 20px;">
				<thead>
					<tr>
						<th>团体名称</th>
						<th>访问模式</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="community-list">
					<tr>
						<td colspan="3" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	function saveSNMPConfig() {
		alert('SNMP 配置保存功能待实现');
	}

	function addCommunity() {
		alert('添加团体功能待实现');
	}

	function loadCommunities() {
		fetch('/api/v1/snmp/communities')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('community-list');
				if (data.code === 200 && data.data && data.data.communities) {
					var html = '';
					data.data.communities.forEach(function(comm) {
						html += '<tr>' +
							'<td>' + comm.name + '</td>' +
							'<td>' + comm.access + '</td>' +
							'<td><button class="btn btn-sm btn-danger" onclick="deleteCommunity(\\'' + comm.name + '\\')">删除</button></td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="3" style="text-align:center; color:#999;">暂无团体</td></tr>';
				}
			});
	}

	function deleteCommunity(name) {
		if (confirm('确定要删除团体 "' + name + '" 吗？')) {
			alert('删除功能待实现');
		}
	}

	// 页面加载时获取团体列表
	loadCommunities();
	</script>
	`

	boxContent := template.HTML(snmpContent)

	snmpBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-bell"></i> SNMP 配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(snmpBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "SNMP 配置",
		Description: "维护 → SNMP 配置",
	}, nil
}

// GetSNMPTrapContent SNMP Trap 配置页面
func GetSNMPTrapContent(ctx *context.Context) (types.Panel, error) {
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

// GetWormProtectionContent 蠕虫攻击防护页面
func GetWormProtectionContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	wormContent := `
	<style>
		.worm-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.worm-table th, .worm-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.worm-table th { background-color: #f5f5f5; }
		.status-badge { padding: 3px 8px; border-radius: 3px; font-size: 12px; }
		.status-enabled { background-color: #5cb85c; color: white; }
		.status-disabled { background-color: #999; color: white; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">蠕虫规则列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-primary" onclick="showAddRuleModal()">
						<i class="fa fa-plus"></i> 新建
					</button>
					<button type="button" class="btn btn-danger" onclick="deleteRules()">
						<i class="fa fa-trash"></i> 删除
					</button>
					<button type="button" class="btn btn-default" onclick="clearStats()">
						<i class="fa fa-eraser"></i> 清除统计
					</button>
					<button type="button" class="btn btn-default" onclick="loadRules()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="worm-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>名称</th>
						<th>协议类型</th>
						<th>目的端口</th>
						<th>攻击统计</th>
						<th>状态</th>
						<th>操作</th>
					</tr>
				</thead>
				<tbody id="worm-rule-list">
					<tr>
						<td colspan="7" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
			<div style="margin-top: 10px; color: #999;">总条目：<span id="rule-count">0</span></div>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#worm-rule-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function loadRules() {
		fetch('/api/v1/security/worm/rules')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('worm-rule-list');
				if (data.code === 200 && data.data && data.data.rules) {
					var html = '';
					var count = 0;
					data.data.rules.forEach(function(rule) {
						var statusClass = rule.enabled ? 'status-enabled' : 'status-disabled';
						var statusText = rule.enabled ? 'enable' : 'disable';
						html += '<tr>' +
							'<td><input type="checkbox" value="' + rule.id + '"></td>' +
							'<td>' + rule.name + '</td>' +
							'<td>' + rule.protocol + '</td>' +
							'<td>' + rule.port + '</td>' +
							'<td>' + rule.stats + '</td>' +
							'<td><span class="status-badge ' + statusClass + '">' + statusText + '</span></td>' +
							'<td><button class="btn btn-sm btn-primary" onclick="editRule(\\'' + rule.id + '\\')">编辑</button></td>' +
							'</tr>';
						count++;
					});
					tbody.innerHTML = html;
					document.getElementById('rule-count').textContent = count;
				} else {
					tbody.innerHTML = '<tr><td colspan="7" style="text-align:center; color:#999;">暂无规则</td></tr>';
					document.getElementById('rule-count').textContent = '0';
				}
			});
	}

	function showAddRuleModal() {
		alert('新建规则功能待实现');
	}

	function editRule(id) {
		alert('编辑规则功能待实现');
	}

	function deleteRules() {
		var selected = document.querySelectorAll('#worm-rule-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的规则');
			return;
		}
		if (confirm('确定要删除选中的规则吗？')) {
			alert('删除功能待实现');
		}
	}

	function clearStats() {
		if (confirm('确定要清除所有统计吗？')) {
			fetch('/api/v1/security/worm/clear-stats', { method: 'POST' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('统计已清除');
						loadRules();
					}
				});
		}
	}

	// 页面加载时获取规则列表
	loadRules();
	</script>
	`

	boxContent := template.HTML(wormContent)

	wormBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-bug"></i> 蠕虫攻击防护`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(wormBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "蠕虫攻击防护",
		Description: "维护 → 蠕虫攻击防护",
	}, nil
}

// GetDDoSProtectionContent DDoS 攻击防护页面
func GetDDoSProtectionContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	ddosContent := `
	<style>
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 150px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; display: inline-block; }
		.checkbox-group label { display: inline-block; margin-right: 20px; font-weight: normal; }
		.unit { margin-left: 5px; color: #999; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">DDoS 防护参数</h3>
		</div>
		<div class="box-body">
			<h4>参数信息</h4>
			<div class="form-group">
				<label>ICMP Flooding 防护</label>
				<input type="number" class="form-control" id="icmp-threshold" value="0" min="0" max="1000">
				<span class="unit">/ pps (0-1000)</span>
			</div>
			<div class="form-group">
				<label>UDP Flooding 防护</label>
				<input type="number" class="form-control" id="udp-threshold" value="0" min="0" max="1000">
				<span class="unit">/ pps (0-1000)</span>
			</div>
			<div class="form-group">
				<label>SYN Flooding 防护</label>
				<input type="number" class="form-control" id="syn-threshold" value="0" min="0" max="1000">
				<span class="unit">/ pps (0-1000)</span>
			</div>
			<div class="form-group">
				<label>小包攻击防护</label>
				<input type="number" class="form-control" id="small-packet-threshold" value="64" min="28" max="65535">
				<span class="unit">bytes (28-65535)</span>
			</div>
			<h4>特殊攻击防护</h4>
			<div class="checkbox-group">
				<label><input type="checkbox" id="smurf-protection"> Smurf 攻击防护</label>
				<label><input type="checkbox" id="fraggle-protection"> Fraggle 攻击防护</label>
				<label><input type="checkbox" id="mac-equal-protection"> MAC 地址相等防护</label>
				<label><input type="checkbox" id="ip-equal-protection"> IP 地址相等防护</label>
			</div>
			<div style="margin-top: 20px;">
				<button type="button" class="btn btn-primary" onclick="saveDDoSConfig()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<script>
	function saveDDoSConfig() {
		alert('DDoS 防护配置保存功能待实现');
	}

	function loadDDoSConfig() {
		fetch('/api/v1/security/ddos/config')
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data) {
					document.getElementById('icmp-threshold').value = data.data.icmp_threshold || 0;
					document.getElementById('udp-threshold').value = data.data.udp_threshold || 0;
					document.getElementById('syn-threshold').value = data.data.syn_threshold || 0;
					document.getElementById('small-packet-threshold').value = data.data.small_packet_threshold || 64;
					document.getElementById('smurf-protection').checked = data.data.smurf_protection || false;
					document.getElementById('fraggle-protection').checked = data.data.fraggle_protection || false;
					document.getElementById('mac-equal-protection').checked = data.data.mac_equal_protection || false;
					document.getElementById('ip-equal-protection').checked = data.data.ip_equal_protection || false;
				}
			});
	}

	// 页面加载时获取配置
	loadDDoSConfig();
	</script>
	`

	boxContent := template.HTML(ddosContent)

	ddosBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-shield"></i> DDoS 攻击防护`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(ddosBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "DDoS 攻击防护",
		Description: "维护 → DDoS 攻击防护",
	}, nil
}

// GetARPProtectionContent ARP 攻击防护页面
func GetARPProtectionContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	arpContent := `
	<style>
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 150px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; display: inline-block; }
		.unit { margin-left: 5px; color: #999; }
		.warning-text { color: #d9534f; font-weight: bold; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">ARP 防护参数</h3>
		</div>
		<div class="box-body">
			<h4>参数信息</h4>
			<div class="form-group">
				<label>Arp 防护</label>
				<input type="number" class="form-control" id="arp-threshold" value="100" min="0" max="1000000">
				<span class="unit">/ pps (0-1000000)</span>
			</div>
			<p class="warning-text">配置为 0 会导致交换机不再学习动态 ARP!</p>
			<div style="margin-top: 20px;">
				<button type="button" class="btn btn-primary" onclick="saveARPConfig()">
					<i class="fa fa-save"></i> 应用
				</button>
			</div>
		</div>
	</div>

	<script>
	function saveARPConfig() {
		var threshold = parseInt(document.getElementById('arp-threshold').value);
		if (threshold === 0) {
			var confirmation = prompt('⚠️ 严重警告：此操作可能导致您永久失去设备管理权限！\\n\\n设置为 0 后，交换机将立即停止学习动态 ARP 表项。\\n\\n后果：\\n• 几分钟后，您电脑的 MAC 地址将从交换机 ARP 表中老化消失\\n• 交换机将丢弃所有来自您电脑的管理请求\\n• 您将无法通过 Web/SSH/Telnet 连接设备\\n• 只能前往机房，通过 Console 线或重启设备恢复\\n\\n如果您确定要设置为 0，请先完成以下准备工作：\\n☐ 在交换机上为当前管理 PC 配置静态 ARP 表项\\n☐ 或确保有其他方式可以访问设备（如 Console 线）\\n\\n请输入 I UNDERSTAND 确认您已知晓风险：');
			if (confirmation !== 'I UNDERSTAND') {
				if (confirmation !== null) {
					alert('确认失败，请输入 I UNDERSTAND');
				}
				document.getElementById('arp-threshold').value = 100; // 恢复默认值
				return;
			}
		}

		alert('ARP 防护配置保存功能待实现');
	}

	function loadARPConfig() {
		fetch('/api/v1/security/arp/config')
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data) {
					document.getElementById('arp-threshold').value = data.data.arp_threshold || 100;
				}
			});
	}

	// 页面加载时获取配置
	loadARPConfig();
	</script>
	`

	boxContent := template.HTML(arpContent)

	arpBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-lock"></i> ARP 攻击防护`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(arpBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "ARP 攻击防护",
		Description: "维护 → ARP 攻击防护",
	}, nil
}

// GetSessionsContent 当前会话页面
func GetSessionsContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	sessionsContent := `
	<style>
		.session-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.session-table th, .session-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.session-table th { background-color: #f5f5f5; }
		.current-session { color: #5cb85c; font-weight: bold; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">当前会话列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="deleteSessions()">
						<i class="fa fa-trash"></i> 删除
					</button>
					<button type="button" class="btn btn-default" onclick="loadSessions()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<table class="session-table">
				<thead>
					<tr>
						<th><input type="checkbox" id="select-all"></th>
						<th>用户名</th>
						<th>会话 ID</th>
						<th>超时时间</th>
						<th>客户端 IP</th>
					</tr>
				</thead>
				<tbody id="session-list">
					<tr>
						<td colspan="5" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<script>
	document.getElementById('select-all').addEventListener('change', function(e) {
		var checkboxes = document.querySelectorAll('#session-list input[type="checkbox"]');
		checkboxes.forEach(cb => cb.checked = e.target.checked);
	});

	function loadSessions() {
		fetch('/api/v1/sessions')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('session-list');
				if (data.code === 200 && data.data && data.data.sessions) {
					var html = '';
					data.data.sessions.forEach(function(session) {
						var currentMarker = session.is_current ? ' <span class="current-session">(*)</span>' : '';
						html += '<tr>' +
							'<td><input type="checkbox" value="' + session.session_id + '"></td>' +
							'<td>' + session.username + currentMarker + '</td>' +
							'<td>' + session.session_id + '</td>' +
							'<td>' + session.timeout_at + '</td>' +
							'<td>' + session.client_ip + '</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无会话</td></tr>';
				}
			});
	}

	function deleteSessions() {
		var selected = document.querySelectorAll('#session-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要终止的会话');
			return;
		}
		if (confirm('确定要终止选中的会话吗？')) {
			alert('删除功能待实现');
		}
	}

	// 页面加载时获取会话列表
	loadSessions();
	</script>
	`

	boxContent := template.HTML(sessionsContent)

	sessionsBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-clock-o"></i> 当前会话`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(sessionsBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "当前会话",
		Description: "维护 → 当前会话",
	}, nil
}
