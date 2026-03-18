package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getLogsContent 日志管理页面
func getLogsContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	logsContent := `
	<style>
		.log-table { width: 100%; border-collapse: collapse; font-size: 12px; }
		.log-table th, .log-table td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
		.log-table th { background-color: #f5f5f5; }
		.log-level-info { color: #337ab7; }
		.log-level-warn { color: #f0ad4e; }
		.log-level-error { color: #d9534f; }
		.filter-group { margin-bottom: 15px; }
		.filter-group select { padding: 5px; border: 1px solid #ddd; border-radius: 4px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">系统日志</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="clearLogs()">
						<i class="fa fa-trash"></i> 清除日志
					</button>
					<button type="button" class="btn btn-default" onclick="loadLogs()">
						<i class="fa fa-refresh"></i> 刷新
					</button>
				</div>
			</div>
		</div>
		<div class="box-body">
			<div class="filter-group">
				<label>日志级别：</label>
				<select id="log-level-filter" onchange="loadLogs()">
					<option value="all">全部</option>
					<option value="info">信息</option>
					<option value="warn">警告</option>
					<option value="error">错误</option>
				</select>
			</div>
			<table class="log-table">
				<thead>
					<tr>
						<th>时间</th>
						<th>级别</th>
						<th>模块</th>
						<th>消息</th>
					</tr>
				</thead>
				<tbody id="log-list">
					<tr>
						<td colspan="4" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
			<div style="margin-top: 10px; color: #999;">总条数：<span id="log-count">0</span></div>
		</div>
	</div>

	<script>
	function loadLogs() {
		var level = document.getElementById('log-level-filter').value;
		var url = '/api/v1/logs';
		if (level !== 'all') {
			url += '?level=' + level;
		}
		fetch(url)
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('log-list');
				if (data.code === 200 && data.data && data.data.logs) {
					var html = '';
					var count = 0;
					data.data.logs.forEach(function(log) {
						var levelClass = 'log-level-' + log.level;
						html += '<tr>' +
							'<td>' + log.timestamp + '</td>' +
							'<td class="' + levelClass + '">' + log.level.toUpperCase() + '</td>' +
							'<td>' + log.module + '</td>' +
							'<td>' + log.message + '</td>' +
							'</tr>';
						count++;
					});
					tbody.innerHTML = html;
					document.getElementById('log-count').textContent = count;
				} else {
					tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; color:#999;">暂无日志</td></tr>';
					document.getElementById('log-count').textContent = '0';
				}
			});
	}

	function clearLogs() {
		if (confirm('确定要清除所有日志吗？\\n\\n此操作不可恢复！')) {
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

	// 页面加载时获取日志列表
	loadLogs();
	</script>
	`

	boxContent := template.HTML(logsContent)

	logsBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-file-text"></i> 日志管理`).
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
