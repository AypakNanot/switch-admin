package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSessionsContent 当前会话页面
func getSessionsContent(ctx *context.Context) (types.Panel, error) {
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
