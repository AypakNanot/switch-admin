package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getWormProtectionContent 蠕虫攻击防护页面
func getWormProtectionContent(ctx *context.Context) (types.Panel, error) {
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
