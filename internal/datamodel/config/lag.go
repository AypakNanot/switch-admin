package config

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	tmpl "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// getLinkAggregationContent 链路聚合页面
func getLinkAggregationContent(ctx *context.Context) (types.Panel, error) {
	components := tmpl.Default(ctx)
	colComp := components.Col()
	boxComp := components.Box()

	laContent := `
<style>
	.config-section { margin-bottom: 30px; }
	.status-normal { color: #28a745; }
	.status-degraded { color: #ffc107; }
	.status-down { color: #dc3545; }
</style>

<div class="config-section">
	<div class="box box-default">
		<div class="box-header with-border">
			<h3 class="box-title">链路聚合</h3>
			<div class="box-tools">
				<button type="button" class="btn btn-success btn-sm" onclick="showCreateModal()">
					<i class="fa fa-plus"></i> 新建
				</button>
				<button type="button" class="btn btn-danger btn-sm" onclick="deleteSelected()">
					<i class="fa fa-trash"></i> 删除
				</button>
				<button type="button" class="btn btn-primary btn-sm" onclick="refreshData()">
					<i class="fa fa-refresh"></i> 刷新
				</button>
			</div>
		</div>
		<div class="box-body">
			<table class="table table-bordered table-hover">
				<thead>
					<tr>
						<th width="50"><input type="checkbox" id="select-all"></th>
						<th>聚合组</th>
						<th>模式</th>
						<th>成员端口</th>
						<th>负载均衡</th>
						<th>最小活跃</th>
						<th>状态</th>
						<th width="150">操作</th>
					</tr>
				</thead>
				<tbody id="la-table-body">
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag1</td>
						<td>LACP</td>
						<td>GE1/0/1-4</td>
						<td>src-dst-ip</td>
						<td>2</td>
						<td><span class="status-normal">正常</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag1')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag1')">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag2</td>
						<td>Static</td>
						<td>GE1/0/5-6</td>
						<td>src-dst-mac</td>
						<td>1</td>
						<td><span class="status-normal">正常</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag2')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag2')">删除</button>
						</td>
					</tr>
					<tr>
						<td><input type="checkbox" class="row-checkbox"></td>
						<td>Ag3</td>
						<td>LACP</td>
						<td>GE1/0/9-12</td>
						<td>src-dst-mac</td>
						<td>2</td>
						<td><span class="status-degraded">降级</span></td>
						<td>
							<button class="btn btn-sm btn-primary" onclick="editLA('Ag3')">编辑</button>
							<button class="btn btn-sm btn-danger" onclick="deleteLA('Ag3')">删除</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>

<!-- 新建/编辑聚合组弹窗 -->
<div class="modal fade" id="laModal" tabindex="-1">
	<div class="modal-dialog modal-lg">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="laModalTitle">新建链路聚合组</h4>
			</div>
			<div class="modal-body">
				<form id="la-form">
					<input type="hidden" id="la-id">

					<div class="box box-default">
						<div class="box-header with-border">
							<h4 class="box-title">基本信息</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">聚合组 ID</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-group-id" min="1" max="128" placeholder="范围：1-128">
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">模式</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-mode">
										<option value="LACP">LACP（动态）</option>
										<option value="Static">Static（静态）</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">描述</label>
								<div class="col-sm-9">
									<input type="text" class="form-control" id="la-description" placeholder="可选">
								</div>
							</div>
						</div>
					</div>

					<div class="box box-success">
						<div class="box-header with-border">
							<h4 class="box-title">负载均衡策略</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">负载均衡</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-load-balance">
										<option value="src-mac">src-mac</option>
										<option value="dst-mac">dst-mac</option>
										<option value="src-dst-mac">src-dst-mac</option>
										<option value="src-ip">src-ip</option>
										<option value="dst-ip">dst-ip</option>
										<option value="src-dst-ip" selected>src-dst-ip</option>
									</select>
								</div>
							</div>
						</div>
					</div>

					<div class="box box-info">
						<div class="box-header with-border">
							<h4 class="box-title">成员端口配置</h4>
						</div>
						<div class="box-body">
							<div class="form-group">
								<label>选择端口：</label>
								<div id="port-selection">
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/1"> GE1/0/1</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/2"> GE1/0/2</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/3"> GE1/0/3</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/4"> GE1/0/4</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/5"> GE1/0/5</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/6"> GE1/0/6</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/7"> GE1/0/7</label>
									<label class="checkbox-inline"><input type="checkbox" class="port-checkbox" value="GE1/0/8"> GE1/0/8</label>
								</div>
							</div>
							<div class="alert alert-warning">
								<i class="fa fa-exclamation-triangle"></i> 需确保：
								<ul style="margin: 5px 0 0 20px;">
									<li>所有端口速率/双工配置一致</li>
									<li>所有端口未加入其他聚合组</li>
									<li>所有端口未配置镜像/安全等功能</li>
								</ul>
							</div>
						</div>
					</div>

					<div class="box box-warning">
						<div class="box-header with-border">
							<h4 class="box-title">高级选项（LACP 模式）</h4>
						</div>
						<div class="box-body">
							<div class="form-group row">
								<label class="col-sm-3 control-label">最小活跃链路</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-min-active" min="1" value="1">
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">LACP 超时</label>
								<div class="col-sm-9">
									<select class="form-control" id="la-lacp-timeout">
										<option value="long">长超时 (90s)</option>
										<option value="short">短超时 (3s)</option>
									</select>
								</div>
							</div>
							<div class="form-group row">
								<label class="col-sm-3 control-label">LACP 优先级</label>
								<div class="col-sm-9">
									<input type="number" class="form-control" id="la-priority" min="1" max="65535" value="32768">
								</div>
							</div>
						</div>
					</div>
				</form>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
				<button type="button" class="btn btn-primary" onclick="saveLA()">确定</button>
			</div>
		</div>
	</div>
</div>

<script>
$('#select-all').change(function() {
	$('.row-checkbox').prop('checked', this.checked);
});

function refreshData() {
	location.reload();
}

function showCreateModal() {
	$('#laModalTitle').text('新建链路聚合组');
	$('#la-form')[0].reset();
	$('#la-id').val('');
	$('#laModal').modal('show');
}

function editLA(laId) {
	$('#laModalTitle').text('编辑链路聚合组 - ' + laId);
	$('#laModal').modal('show');
}

function deleteLA(laId) {
	layer.confirm('确定要删除聚合组 ' + laId + ' 吗？成员端口将恢复独立管理。', {
		btn: ['确定', '取消']
	}, function() {
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function deleteSelected() {
	var selected = [];
	$('.row-checkbox:checked').each(function() {
		selected.push($(this).closest('tr').find('td:eq(1)').text());
	});
	if (selected.length === 0) {
		layer.msg('请先选择要删除的聚合组', {icon: 2});
		return;
	}
	layer.confirm('确定要删除选中的 ' + selected.length + ' 个聚合组吗？', {
		btn: ['确定', '取消']
	}, function() {
		layer.msg('删除成功', {icon: 1});
		setTimeout(function() { location.reload(); }, 1000);
	});
}

function saveLA() {
	var laId = $('#la-id').val();
	var selectedPorts = [];
	$('.port-checkbox:checked').each(function() {
		selectedPorts.push($(this).val());
	});

	var data = {
		group_id: parseInt($('#la-group-id').val()),
		mode: $('#la-mode').val(),
		description: $('#la-description').val(),
		load_balance: $('#la-load-balance').val(),
		member_ports: selectedPorts,
		min_active: parseInt($('#la-min-active').val()),
		lacp_timeout: $('#la-lacp-timeout').val(),
		priority: parseInt($('#la-priority').val())
	};

	layer.msg('保存成功', {icon: 1});
	$('#laModal').modal('hide');
	setTimeout(function() { location.reload(); }, 1000);
}
</script>
`

	boxContent := template.HTML(laContent)
	box := boxComp.
		WithHeadBorder().
		SetHeader(`<i class="fa fa-link"></i> 链路聚合`).
		SetBody(boxContent).
		GetContent()

	colContent := colComp.SetSize(types.SizeMD(12)).SetContent(box).GetContent()
	rowContent := components.Row().SetContent(colContent).GetContent()

	panel := types.Panel{
		Content:     rowContent,
		Title:       "链路聚合",
		Description: "将多个物理端口捆绑成一个逻辑端口，提升带宽和可靠性。支持 LACP 动态聚合和静态聚合。",
	}

	return panel, nil
}
