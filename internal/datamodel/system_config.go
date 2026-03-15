package datamodel

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// GetSystemConfigPage 获取系统配置页面
func GetSystemConfigPage(ctx *context.Context) (types.Panel, error) {
	fmt.Println("GetSystemConfigPage called")

	content := template.HTML(`
	<div class="row">
		<div class="col-md-12">
			<div class="box box-primary">
				<div class="box-header with-border">
					<h3 class="box-title">系统配置</h3>
				</div>
				<div class="box-body">
					<table class="table table-bordered">
						<tr>
							<th width="150">运行模式</th>
							<td>
								<span class="label label-info">交换机模式</span>
								<span class="text-muted">(switch)</span>
							</td>
						</tr>
						<tr>
							<th>数据库</th>
							<td>SQLite3 (data/admin.db)</td>
						</tr>
						<tr>
							<th>GoAdmin 版本</th>
							<td>v1.2.26</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	</div>
	`)

	fmt.Println("GetSystemConfigPage returning panel")
	return types.Panel{
		Content:     content,
		Title:       "系统配置",
		Description: "查看和管理系统配置",
	}, nil
}
