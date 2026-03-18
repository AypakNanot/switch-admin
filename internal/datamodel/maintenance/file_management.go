package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getFilesContent 文件管理页面
func getFilesContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	filesContent := `
	<style>
		.file-table { width: 100%; border-collapse: collapse; font-size: 13px; }
		.file-table th, .file-table td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
		.file-table th { background-color: #f5f5f5; }
		.form-group { margin-bottom: 15px; }
		.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
		.form-control { width: 300px; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">上传文件</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>选择文件</label>
				<input type="file" class="form-control" id="file-input">
			</div>
			<button type="button" class="btn btn-primary" onclick="uploadFile()">
				<i class="fa fa-upload"></i> 上传
			</button>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">上传固件</h3>
		</div>
		<div class="box-body">
			<div class="form-group">
				<label>选择固件文件</label>
				<input type="file" class="form-control" id="firmware-input" accept=".bin,.img">
			</div>
			<button type="button" class="btn btn-warning" onclick="uploadFirmware()">
				<i class="fa fa-upload"></i> 上传固件
			</button>
			<p style="margin-top: 10px; color: #999;">支持的固件格式：.bin, .img</p>
		</div>
	</div>

	<div class="box box-default" style="margin-top: 20px;">
		<div class="box-header with-border">
			<h3 class="box-title">文件列表</h3>
			<div class="box-tools">
				<div class="btn-group">
					<button type="button" class="btn btn-danger" onclick="deleteFiles()">
						<i class="fa fa-trash"></i> 删除
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
						<th>大小</th>
						<th>修改时间</th>
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

	function uploadFile() {
		var fileInput = document.getElementById('file-input');
		if (!fileInput.files.length) {
			alert('请选择文件');
			return;
		}
		var formData = new FormData();
		formData.append('file', fileInput.files[0]);
		fetch('/api/v1/files/upload', {
			method: 'POST',
			body: formData
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('文件上传成功');
				loadFiles();
			} else {
				alert(data.message || '上传失败');
			}
		});
	}

	function uploadFirmware() {
		var fileInput = document.getElementById('firmware-input');
		if (!fileInput.files.length) {
			alert('请选择固件文件');
			return;
		}
		var formData = new FormData();
		formData.append('file', fileInput.files[0]);
		fetch('/api/v1/files/firmware', {
			method: 'POST',
			body: formData
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				alert('固件上传成功');
				loadFiles();
			} else {
				alert(data.message || '上传失败');
			}
		});
	}

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
							'<td>' + file.size + '</td>' +
							'<td>' + file.modified + '</td>' +
							'<td>' +
								'<a href="/api/v1/files/download?path=' + encodeURIComponent(file.path) + '" class="btn btn-sm btn-default"><i class="fa fa-download"></i> 下载</a> ' +
								'<button class="btn btn-sm btn-danger" onclick="deleteFile(\\'' + file.path + '\\')"><i class="fa fa-trash"></i> 删除</button>' +
							'</td>' +
							'</tr>';
					});
					tbody.innerHTML = html;
				} else {
					tbody.innerHTML = '<tr><td colspan="5" style="text-align:center; color:#999;">暂无文件</td></tr>';
				}
			});
	}

	function deleteFile(path) {
		if (confirm('确定要删除此文件吗？')) {
			fetch('/api/v1/files?path=' + encodeURIComponent(path), { method: 'DELETE' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('文件删除成功');
						loadFiles();
					} else {
						alert(data.message || '删除失败');
					}
				});
		}
	}

	function deleteFiles() {
		var selected = document.querySelectorAll('#file-list input[type="checkbox"]:checked');
		if (selected.length === 0) {
			alert('请选择要删除的文件');
			return;
		}
		if (confirm('确定要删除选中的文件吗？')) {
			alert('批量删除功能待实现');
		}
	}

	// 页面加载时获取文件列表
	loadFiles();
	</script>
	`

	boxContent := template.HTML(filesContent)

	filesBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-folder-open"></i> 文件管理`).
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
