package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getLoadConfigContent 加载配置页面
func getLoadConfigContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	loadConfigContent := `
	<style>
		.config-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.config-table th, .config-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.config-table th { background-color: #f5f5f5; }
		.warning-text { color: #d9534f; font-weight: bold; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">配置文件列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-primary" onclick="loadConfig()">
						<i class="fa fa-download"></i> 加载选中配置
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
						<th><input type="radio" name="config-select" id="select-config"></th>
						<th>文件路径</th>
						<th>大小</th>
						<th>修改时间</th>
					</tr>
				</thead>
				<tbody id="config-list">
					<tr>
						<td colspan="4" style="text-align:center; color:#999;">加载中...</td>
					</tr>
				</tbody>
			</table>
			<p class="warning-text" style="margin-top: 15px;">⚠️ 注意：加载配置后，部分配置可能需要重启设备才能生效。</p>
		</div>
	</div>

	<script>
	function loadConfigFiles() {
		fetch('/api/v1/config/files')
			.then(res => res.json())
			.then(data => {
				var tbody = document.getElementById('config-list');
				if (data.code === 200 && data.data && data.data.files) {
					var html = '';
					data.data.files.forEach(function(file) {
						html += '<tr>' +
							'<td><input type="radio" name="config-select" value="' + file.file_path + '"></td>' +
							'<td>' + file.file_path + '</td>' +
							'<td>' + file.size + '</td>' +
							'<td>' + file.modified + '</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="4" style="text-align:center; color:#999;">暂无配置文件</td></tr>';
				}
			});
	}

	function loadConfig() {
		var selected = document.querySelector('input[name="config-select"]:checked');
		if (!selected) {
			alert('请选择要加载的配置文件');
			return;
		}

		if (confirm('确定要加载选中的配置文件吗？\\n\\n部分配置可能需要重启设备才能生效。')) {
			fetch('/api/v1/config/load', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({ file_path: selected.value })
			})
			.then(res => res.json())
			.then(data => {
				if (data.code === 200) {
					alert(data.message || '配置加载成功');
				} else {
					alert(data.message || '加载失败');
				}
			});
		}
	}

	// 页面加载时获取配置文件列表
	loadConfigFiles();
	</script>
	`

	boxContent := template.HTML(loadConfigContent)

	loadConfigBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-folder-open"></i> 加载配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(loadConfigBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "加载配置",
		Description: "维护 → 加载配置",
	}, nil
}
