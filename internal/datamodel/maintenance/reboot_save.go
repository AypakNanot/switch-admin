package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getRebootSaveContent 重启/保存配置页面
func getRebootSaveContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	rebootSaveContent := `
	<style>
		.config-section { margin-bottom: 30px; }
		.btn-group { margin-top: 20px; }
		.warning-text { color: #d9534f; font-weight: bold; margin-top: 15px; }
	</style>

	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">保存配置</h3>
		</div>
		<div class="box-body">
			<p>当前配置仅保存在内存中，重启后会丢失。请点击"保存配置"按钮将当前配置保存到启动配置文件。</p>
			<div class="btn-group">
				<button type="button" class="btn btn-primary" onclick="saveConfig()">
					<i class="fa fa-save"></i> 保存配置
				</button>
			</div>
		</div>
	</div>

	<div class="box box-default config-section">
		<div class="box-header with-border">
			<h3 class="box-title">重启设备</h3>
		</div>
		<div class="box-body">
			<p>重启交换机设备。</p>
			<div class="btn-group">
				<button type="button" class="btn btn-warning" onclick="reboot()">
					<i class="fa fa-refresh"></i> 重启交换机
				</button>
			</div>
			<p class="warning-text">⚠️ 警告：重启会导致网络短暂中断！</p>
		</div>
	</div>

	<div class="box box-default config-section">
		<div class="box-header with-border">
			<h3 class="box-title">恢复出厂设置</h3>
		</div>
		<div class="box-body">
			<p>删除所有配置文件，恢复设备到出厂默认状态。</p>
			<div class="btn-group">
				<button type="button" class="btn btn-danger" onclick="factoryReset()">
					<i class="fa fa-trash"></i> 恢复出厂设置
				</button>
			</div>
			<p class="warning-text">⚠️ 严重警告：此操作将删除所有配置，且不可恢复！</p>
		</div>
	</div>

	<script>
	function saveConfig() {
		if (confirm('确定要保存当前配置到启动文件吗？')) {
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

	function reboot() {
		if (confirm('确定要重启交换机吗？\\n\\n警告：重启会导致网络短暂中断！')) {
			fetch('/api/v1/system/reboot', { method: 'POST' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('重启命令已发送，设备将在 30 秒后重启');
					} else {
						alert(data.message || '重启失败');
					}
				});
		}
	}

	function factoryReset() {
		var confirmation = prompt('⚠️ 严重警告：此操作将删除所有配置，且不可恢复！\\n\\n如果您确定要恢复出厂设置，请输入 CONFIRM 确认：');
		if (confirmation === 'CONFIRM') {
			fetch('/api/v1/system/factory-reset', { method: 'POST' })
				.then(res => res.json())
				.then(data => {
					if (data.code === 200) {
						alert('恢复出厂设置成功，设备将重启');
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

	rebootSaveBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-cog"></i> 重启/保存配置`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(rebootSaveBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "重启/保存配置",
		Description: "维护 → 重启/保存配置",
	}, nil
}
