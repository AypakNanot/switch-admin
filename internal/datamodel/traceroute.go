package datamodel

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetTracerouteContent Traceroute 诊断页面
func GetTracerouteContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	tracerouteForm := `
	<form id="traceroute-form" class="form-inline" role="form">
		<div class="form-group">
			<label>VRF ID</label>
			<select class="form-control" id="vrf_id">
				<option value="mgmt vrf">mgmt vrf</option>
				<option value="default vrf">default vrf</option>
			</select>
		</div>
		<div class="form-group">
			<label>目标 IP 地址</label>
			<input type="text" class="form-control" id="target" placeholder="10.10.25.30" style="width:250px;">
		</div>
		<button type="button" class="btn btn-primary" id="trace-btn" onclick="startTraceroute()">
			<i class="fa fa-play"></i> Traceroute
		</button>
	</form>

	<div id="traceroute-result" style="margin-top:20px; display:none;">
		<pre id="traceroute-output" style="background:#f5f5f5;padding:15px;border-radius:4px;font-size:12px;min-height:200px;"></pre>
	</div>

	<script>
	var currentTaskId = null;
	var pollTimer = null;

	function startTraceroute() {
		var vrfId = document.getElementById('vrf_id').value;
		var target = document.getElementById('target').value;

		if (!target) {
			alert('请输入目标 IP 地址');
			return;
		}

		document.getElementById('trace-btn').disabled = true;
		document.getElementById('trace-btn').innerHTML = '<i class="fa fa-spinner fa-spin"></i> 执行中...';
		document.getElementById('traceroute-result').style.display = 'block';
		document.getElementById('traceroute-output').textContent = '正在执行 Traceroute 到 ' + target + '，最多 30 跳...\n\n';

		fetch('/api/v1/diagnostic/traceroute', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({
				vrf_id: vrfId,
				target: target,
				max_hops: 30,
				timeout: 2,
				probes: 3
			})
		})
		.then(res => res.json())
		.then(data => {
			if (data.code === 200) {
				currentTaskId = data.data.task_id;
				pollTimer = setInterval(() => pollTracerouteResult(), 2000);
			} else {
				alert(data.message || '创建任务失败');
				resetButton();
			}
		});
	}

	function pollTracerouteResult() {
		if (!currentTaskId) return;

		fetch('/api/v1/diagnostic/traceroute/' + currentTaskId)
			.then(res => res.json())
			.then(data => {
				if (data.code === 200 && data.data.status === 'completed') {
					clearInterval(pollTimer);
					displayTracerouteResult(data.data);
					resetButton();
				} else if (data.code === 200 && data.data.hops) {
					displayTracerouteResult(data.data);
				}
			});
	}

	function displayTracerouteResult(result) {
		var output = document.getElementById('traceroute-output');
		var text = 'traceroute to ' + result.target + ' (' + result.target + '), ' + result.total_hops + ' hops max\n\n';

		result.hops.forEach(function(hop) {
			var hopNum = hop.hop.toString().padStart(2, ' ');
			var ip = hop.ip || '*';
			var times = hop.times && hop.times.length > 0 ? hop.times.join('  ') : '*  *  *';
			text += hopNum + '  ' + ip + '  ' + times + '\n';
		});

		if (result.error) {
			text += '\nError: ' + result.error + '\n';
		}

		output.textContent = text;
	}

	function resetButton() {
		document.getElementById('trace-btn').disabled = false;
		document.getElementById('trace-btn').innerHTML = '<i class="fa fa-play"></i> Traceroute';
	}

	window.addEventListener('beforeunload', function() {
		if (currentTaskId) {
			navigator.sendBeacon('/api/v1/diagnostic/traceroute/' + currentTaskId);
		}
	});
	</script>
	`

	boxContent := template.HTML(tracerouteForm)

	traceBox := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-road"></i> Traceroute - 网络路径追踪`).
		SetHeadColor("#f7f7f7").
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(traceBox).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	return types.Panel{
		Content:     rowContent,
		Title:       "Traceroute 诊断",
		Description: "网络 → Traceroute",
	}, nil
}
