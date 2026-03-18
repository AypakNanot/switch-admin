package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getSNMPContent SNMP 配置页面
func getSNMPContent(ctx *context.Context) (types.Panel, error) {
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
