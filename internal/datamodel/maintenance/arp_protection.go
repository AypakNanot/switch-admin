package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getARPProtectionContent ARP 攻击防护页面
func getARPProtectionContent(ctx *context.Context) (types.Panel, error) {
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
