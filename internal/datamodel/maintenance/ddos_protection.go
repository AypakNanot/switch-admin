package maintenance

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getDDoSProtectionContent DDoS 攻击防护页面
func getDDoSProtectionContent(ctx *context.Context) (types.Panel, error) {
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
