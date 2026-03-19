package system

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// GetUserTable 获取用户管理表格（已有功能，占位）
func GetUserTable(ctx *context.Context) (userTable table.Table) {

	userTable = table.NewDefaultTable(ctx, table.Config{
		Driver:     db.DriverSqlite,
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: table.DefaultPrimaryKeyName,
		},
	})

	info := userTable.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("用户名", "username", db.Varchar)
	info.AddField("邮箱", "email", db.Varchar)
	info.AddField("手机号", "phone", db.Varchar)
	info.AddField("注册时间", "created_at", db.Timestamp)

	info.SetTable("users").SetTitle("用户管理").SetDescription("用户管理列表")

	formList := userTable.GetForm()
	formList.AddField("用户名", "username", db.Varchar, form.Default)
	formList.AddField("邮箱", "email", db.Varchar, form.Email)
	formList.AddField("手机号", "phone", db.Varchar, form.Text)
	formList.AddField("密码", "password", db.Varchar, form.Password)

	formList.SetTable("users").SetTitle("用户管理").SetDescription("用户管理")

	return
}
