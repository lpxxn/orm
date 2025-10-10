package main

import (
	"time"

	"gorm.io/gorm"
)

type SysRole struct {
	RoleId  int       `json:"roleId" gorm:"primaryKey;autoIncrement"` // 角色编码
	SysDept []SysDept `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;"`
	ModelTime
}

type SysDept struct {
	DeptId   int       `json:"deptId" gorm:"primaryKey;column:dept_id;autoIncrement;"` //部门编码
	ParentId int       `json:"parentId" gorm:""`                                       //上级部门
	DeptPath string    `json:"deptPath" gorm:"size:255;"`                              //
	DeptName string    `json:"deptName"  gorm:"size:128;"`                             //部门名称
	Children []SysDept `json:"children" gorm:"-"`
	ModelTime
}

type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (*SysDept) TableName() string {
	return "sys_dept"
}

/*

	•	SysDept []SysDept：声明了多对多关系。
	•	many2many:sys_role_dept → 中间表表名叫 sys_role_dept。
	•	foreignKey:RoleId → 当前模型（SysRole）的主键字段。
	•	joinForeignKey:role_id → 中间表里对应 RoleId 的列名。
	•	references:DeptId → 对端模型（SysDept）的主键字段。
	•	joinReferences:dept_id → 中间表里对应 DeptId 的列名。

👉 最终映射关系：
sys_role.role_id  <--->  sys_role_dept.role_id
sys_dept.dept_id  <--->  sys_role_dept.dept_id

*/
