package datamodel

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetSystemConfigPage 获取系统配置页面
func GetSystemConfigPage(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	content := `
	<div class="row">
		<div class="col-md-6">
			<div class="box box-primary">
				<div class="box-header with-border">
					<h3 class="box-title"><i class="fa fa-cog"></i> 系统信息</h3>
				</div>
				<div class="box-body">
					<table class="table table-bordered">
						<tr>
							<th width="150">数据库</th>
							<td>SQLite3 (data/admin.db)</td>
						</tr>
						<tr>
							<th>GoAdmin 版本</th>
							<td>v1.2.26</td>
						</tr>
						<tr>
							<th>应用名称</th>
							<td>switch-admin</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
		<div class="col-md-6">
			<div class="box box-success">
				<div class="box-header with-border">
					<h3 class="box-title"><i class="fa fa-refresh"></i> 运行模式设置</h3>
				</div>
				<div class="box-body">
					<div class="form-group">
						<label>当前运行模式</label>
						<div id="current-mode-display" style="margin:10px 0;">
							<span class="label label-default" id="mode-badge">加载中...</span>
							<span class="text-muted" id="mode-desc"></span>
						</div>
					</div>
					<div class="form-group">
						<label>切换模式</label>
						<select class="form-control" id="mode-select" style="max-width:200px;">
							<option value="mock">离线测试模式 (mock)</option>
							<option value="switch">交换机模式 (switch)</option>
						</select>
					</div>
					<button type="button" class="btn btn-success" id="switch-mode-btn" onclick="switchMode()">
						<i class="fa fa-check"></i> 切换模式
					</button>
					<span id="switch-result" style="margin-left:10px;"></span>
				</div>
			</div>
		</div>
	</div>

	<script>
	// 加载当前模式
	function loadCurrentMode() {
		fetch('/api/mode')
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					var mode = data.data.mode;
					var desc = data.data.description;
					updateModeDisplay(mode, desc);
					document.getElementById('mode-select').value = mode;
				}
			})
			.catch(err => {
				console.error('加载模式失败:', err);
				document.getElementById('mode-badge').textContent = '加载失败';
			});
	}

	function updateModeDisplay(mode, desc) {
		var badge = document.getElementById('mode-badge');
		var modeDesc = document.getElementById('mode-desc');

		if (mode === 'mock') {
			badge.className = 'label label-warning';
			badge.textContent = '离线测试模式';
		} else if (mode === 'switch') {
			badge.className = 'label label-success';
			badge.textContent = '交换机模式';
		} else {
			badge.className = 'label label-default';
			badge.textContent = mode;
		}
		modeDesc.textContent = '(' + desc + ')';
	}

	function switchMode() {
		var newMode = document.getElementById('mode-select').value;
		var btn = document.getElementById('switch-mode-btn');
		var result = document.getElementById('switch-result');

		// 禁用按钮
		btn.disabled = true;
		btn.innerHTML = '<i class="fa fa-spinner fa-spin"></i> 切换中...';
		result.textContent = '';

		fetch('/api/mode', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({mode: newMode})
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				result.innerHTML = '<span style="color:green;"><i class="fa fa-check"></i> 切换成功!</span>';
				loadCurrentMode();
			} else {
				result.innerHTML = '<span style="color:red;"><i class="fa fa-times"></i> ' + (data.message || '切换失败') + '</span>';
			}
			btn.disabled = false;
			btn.innerHTML = '<i class="fa fa-check"></i> 切换模式';
		})
		.catch(err => {
			result.innerHTML = '<span style="color:red;"><i class="fa fa-times"></i> 请求失败：' + err + '</span>';
			btn.disabled = false;
			btn.innerHTML = '<i class="fa fa-check"></i> 切换模式';
		});
	}

	// 页面加载时获取当前模式
	document.addEventListener('DOMContentLoaded', loadCurrentMode);
	</script>
	`

	boxContent := template.HTML(content)

	modeBox := boxComp.
		WithHeadBorder().
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(modeBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "系统配置",
		Description: "查看和管理系统配置",
	}, nil
}
