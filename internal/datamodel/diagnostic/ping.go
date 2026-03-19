package diagnostic

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetPingContent Ping 诊断页面
func GetPingContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	pingForm := `
	<form id="ping-form" class="form-inline" role="form">
		<div class="form-group">
			<label>VRF ID</label>
			<select class="form-control" id="vrf_id">
				<option value="mgmt vrf">mgmt vrf</option>
				<option value="default vrf">default vrf</option>
			</select>
		</div>
		<div class="form-group">
			<label>目标 IP 地址</label>
			<input type="text" class="form-control" id="target" placeholder="192.168.1.1 或域名" style="width:250px;">
		</div>
		<div class="form-group">
			<label>次数</label>
			<input type="number" class="form-control" id="count" value="4" min="1" max="100" style="width:80px;">
		</div>
		<div class="form-group">
			<label>超时 (秒)</label>
			<input type="number" class="form-control" id="timeout" value="2" min="1" max="60" style="width:80px;">
		</div>
		<div class="form-group">
			<label>间隔 (秒)</label>
			<input type="number" class="form-control" id="interval" value="1" min="1" max="60" style="width:80px;">
		</div>
		<button type="button" class="btn btn-primary" id="ping-btn" onclick="startPing()">
			<i class="fa fa-play"></i> Ping
		</button>
	</form>

	<div id="ping-result" style="margin-top:20px; display:none;">
		<pre id="ping-output" style="background:#f5f5f5;padding:15px;border-radius:4px;font-size:12px;min-height:200px;"></pre>
	</div>

	<script>
	var currentTaskId = null;
	var pollTimer = null;

	function startPing() {
		var vrfId = document.getElementById('vrf_id').value;
		var target = document.getElementById('target').value;
		var count = parseInt(document.getElementById('count').value) || 4;
		var timeout = parseInt(document.getElementById('timeout').value) || 2;
		var interval = parseInt(document.getElementById('interval').value) || 1;

		// 参数验证
		if (!target) {
			alert('请输入目标 IP 地址或域名');
			return;
		}
		if (count < 1 || count > 100) {
			alert('次数必须在 1-100 之间');
			return;
		}
		if (timeout < 1 || timeout > 60) {
			alert('超时时间必须在 1-60 秒之间');
			return;
		}
		if (interval < 1 || interval > 60) {
			alert('间隔时间必须在 1-60 秒之间');
			return;
		}

		document.getElementById('ping-btn').disabled = true;
		document.getElementById('ping-btn').innerHTML = '<i class="fa fa-spinner fa-spin"></i> 执行中...';
		document.getElementById('ping-result').style.display = 'block';
		document.getElementById('ping-output').textContent = '正在执行 Ping ' + target + ' (次数=' + count + ', 超时=' + timeout + 's, 间隔=' + interval + 's)...\n';

		fetch('/api/v1/diagnostic/ping', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({
				vrf_id: vrfId,
				target: target,
				count: count,
				timeout: timeout,
				interval: interval
			})
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				currentTaskId = data.data.task_id;
				pollTimer = setInterval(() => pollPingResult(), 1000);
			} else {
				alert(data.message || '创建任务失败');
				resetButton();
			}
		});
	}

	function pollPingResult() {
		if (!currentTaskId) return;

		fetch('/api/v1/diagnostic/ping/' + currentTaskId)
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data.status === 'completed') {
					clearInterval(pollTimer);
					displayPingResult(data.data);
					resetButton();
				} else {
					document.getElementById('ping-output').textContent += '任务执行中...\n';
				}
			});
	}

	function displayPingResult(result) {
		var output = document.getElementById('ping-output');
		var text = 'PING ' + result.target + ' (' + result.target + '): 56 data bytes\n\n';

		result.results.forEach(function(r) {
			if (r.status === 'success') {
				text += '64 bytes from ' + result.target + ': seq=' + r.seq + ' ttl=' + r.ttl + ' time=' + r.time + '\n';
			} else {
				text += 'Request timeout (seq=' + r.seq + ')\n';
			}
		});

		text += '\n--- ' + result.target + ' ping statistics ---\n';
		text += result.statistics.sent + ' packets transmitted, ' + result.statistics.received + ' received, ' + result.statistics.loss_rate + ' packet loss\n';
		text += 'round-trip min/avg/max = ' + result.statistics.min_time + '/' + result.statistics.avg_time + '/' + result.statistics.max_time + '\n';

		if (result.error) {
			text += '\nError: ' + result.error + '\n';
		}

		output.textContent = text;
	}

	function resetButton() {
		document.getElementById('ping-btn').disabled = false;
		document.getElementById('ping-btn').innerHTML = '<i class="fa fa-play"></i> Ping';
	}

	// 页面离开时清理任务
	window.addEventListener('beforeunload', function() {
		if (currentTaskId) {
			navigator.sendBeacon('/api/v1/diagnostic/ping/' + currentTaskId);
		}
	});
	</script>
	`

	boxContent := template.HTML(pingForm)

	pingBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-plane"></i> Ping - 网络连通性测试`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(pingBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "Ping 诊断",
		Description: "网络 → Ping",
	}, nil
}
