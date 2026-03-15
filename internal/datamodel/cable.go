package datamodel

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetCableTestContent 虚拟电缆检测页面
func GetCableTestContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	cableTestForm := `
	<form id="cable-test-form" class="form-horizontal" role="form">
		<div class="form-group">
			<label class="col-sm-2 control-label">选择端口</label>
			<div class="col-sm-4">
				<select class="form-control" id="port-select">
					<option value="">-- 请选择电口端口 --</option>
				</select>
				<p class="help-block" id="port-hint"></p>
			</div>
		</div>
		<div class="form-group">
			<div class="col-sm-offset-2 col-sm-4">
				<button type="button" class="btn btn-primary" id="detect-btn" onclick="startCableTest()" disabled>
					<i class="fa fa-wrench"></i> 检测
				</button>
			</div>
		</div>
	</form>

	<div id="cable-test-result" style="margin-top:20px; display:none;">
		<div class="box box-info">
			<div class="box-header with-border">
				<h3 class="box-title">检测结果</h3>
			</div>
			<div class="box-body">
				<table class="table table-bordered">
					<tr>
						<th width="150">端口状态</th>
						<td id="result-port-status">-</td>
					</tr>
					<tr>
						<th>电缆状态</th>
						<td id="result-cable-status">-</td>
					</tr>
					<tr>
						<th>电缆长度</th>
						<td id="result-cable-length">-</td>
					</tr>
					<tr>
						<th>故障描述</th>
						<td id="result-fault-desc">-</td>
					</tr>
				</table>

				<h4>线对状态详情</h4>
				<table class="table table-bordered table-condensed">
					<thead>
						<tr>
							<th>线对</th>
							<th>针脚</th>
							<th>状态</th>
							<th>故障距离</th>
						</tr>
					</thead>
					<tbody id="pairs-body">
						<tr><td>Pair A</td><td>1,2</td><td id="pair-a-status">-</td><td id="pair-a-distance">-</td></tr>
						<tr><td>Pair B</td><td>3,6</td><td id="pair-b-status">-</td><td id="pair-b-distance">-</td></tr>
						<tr><td>Pair C</td><td>4,5</td><td id="pair-c-status">-</td><td id="pair-c-distance">-</td></tr>
						<tr><td>Pair D</td><td>7,8</td><td id="pair-d-status">-</td><td id="pair-d-distance">-</td></tr>
					</tbody>
				</table>
			</div>
		</div>
	</div>

	<script>
	var currentTaskId = null;

	// 加载端口列表
	function loadPorts() {
		fetch('/api/v1/diagnostic/cable/ports')
			.then(res => res.json())
			.then(data => {
				var select = document.getElementById('port-select');
				var html = '<option value="">-- 请选择电口端口 --</option>';
				data.data.ports.forEach(function(port) {
					if (port.detectable) {
						html += '<option value="' + port.port_id + '" data-hint="' + (port.hint || '') + '">' + port.label + '</option>';
					}
				});
				// 添加不可选中的端口（仅显示）
				data.data.ports.forEach(function(port) {
					if (!port.detectable) {
						html += '<option value="" disabled class="text-muted">[不可用] ' + port.label + ' - ' + port.hint + '</option>';
					}
				});
				select.innerHTML = html;
			});
	}

	// 端口选择变化
	document.getElementById('port-select') && document.getElementById('port-select').addEventListener('change', function() {
		var selected = this.options[this.selectedIndex];
		var hint = selected.getAttribute('data-hint') || '';
		document.getElementById('port-hint').textContent = hint;
		document.getElementById('detect-btn').disabled = !this.value;
	});

	function startCableTest() {
		var portId = document.getElementById('port-select').value;

		if (!portId) {
			alert('请选择端口');
			return;
		}

		document.getElementById('detect-btn').disabled = true;
		document.getElementById('detect-btn').innerHTML = '<i class="fa fa-spinner fa-spin"></i> 检测中...';
		document.getElementById('cable-test-result').style.display = 'none';

		fetch('/api/v1/diagnostic/cable', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({port_id: portId})
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				displayResult(data.data);
			} else {
				alert(data.message || '检测失败');
			}
			document.getElementById('detect-btn').disabled = false;
			document.getElementById('detect-btn').innerHTML = '<i class="fa fa-wrench"></i> 检测';
		})
		.catch(err => {
			alert('请求失败：' + err);
			document.getElementById('detect-btn').disabled = false;
			document.getElementById('detect-btn').innerHTML = '<i class="fa fa-wrench"></i> 检测';
		});
	}

	function displayResult(result) {
		document.getElementById('cable-test-result').style.display = 'block';

		// 端口状态
		var adminStatus = result.admin_status === 'up' ? 'UP (已使能)' : 'DOWN (已关闭)';
		var linkStatus = result.link_status === 'up' ? 'UP (链路连通)' : 'DOWN (链路断开)';
		document.getElementById('result-port-status').textContent = 'Admin: ' + adminStatus + ', Link: ' + linkStatus;

		// 电缆状态
		var cableStatusMap = {
			'normal': '正常',
			'open': '断路',
			'short': '短路',
			'cross': '线序错误',
			'impedance': '阻抗异常'
		};
		var statusText = cableStatusMap[result.cable_status] || result.cable_status;

		var statusClass = 'label label-success';
		if (result.cable_status === 'open') statusClass = 'label label-danger';
		else if (result.cable_status === 'short') statusClass = 'label label-warning';
		else if (result.cable_status === 'cross') statusClass = 'label label-warning';

		document.getElementById('result-cable-status').innerHTML = '<span class="' + statusClass + '">' + statusText + '</span>';
		document.getElementById('result-cable-length').textContent = result.cable_length || '-';
		document.getElementById('result-fault-desc').textContent = result.fault_description || '-';

		// 线对状态
		var pairs = result.pairs;
		if (pairs) {
			updatePair('a', pairs.pair_a);
			updatePair('b', pairs.pair_b);
			updatePair('c', pairs.pair_c);
			updatePair('d', pairs.pair_d);
		}
	}

	function updatePair(name, pair) {
		var statusMap = {
			'ok': '<span class="label label-success">正常</span>',
			'open': '<span class="label label-danger">断路</span>',
			'short': '<span class="label label-warning">短路</span>',
			'cross': '<span class="label label-warning">线序错误</span>'
		};
		document.getElementById('pair-' + name + '-status').innerHTML = statusMap[pair.status] || pair.status;
		document.getElementById('pair-' + name + '-distance').textContent = pair.fault_distance || '-';
	}

	// 初始化
	document.addEventListener('DOMContentLoaded', function() {
		loadPorts();
	});
	</script>

	<style>
	#cable-test-result .box-body h4 {
		margin-top: 20px;
		margin-bottom: 10px;
		color: #666;
	}
	</style>
	`

	boxContent := template.HTML(cableTestForm)

	cableBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-wrench"></i> 虚拟电缆检测 (VCT)`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(cableBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "虚拟电缆检测",
		Description: "网络 → 虚拟电缆检测",
	}, nil
}
