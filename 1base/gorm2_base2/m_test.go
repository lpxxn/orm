package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	db := getDB()
	db.AutoMigrate(&SysDept{}, &SysRole{})
}

func TestInitData(t *testing.T) {
	db := getDB()
	//tables, _ := db.Migrator().GetTables()
	tables := []string{"sys_dept", "sys_role_dept", "sys_roles"}
	for _, t := range tables {
		db.Exec("TRUNCATE TABLE " + t + " RESTART IDENTITY CASCADE")
	}
	initData()
}

func TestQuery(t *testing.T) {
	db := getDB()
	var role1 SysRole
	db.First(&role1, 1) // RoleId=1
	t.Logf("role1: %+v", role1)

	var deptsForRole1 []SysDept
	db.Where("dept_id IN ?", []int{10, 20}).Find(&deptsForRole1)

	//db.Model(&role1).Association("SysDept").Append([]int{10, 20})
	db.Model(&role1).Association("SysDept").Append([]SysDept{{DeptId: 10}, {DeptId: 20}})
	// 建立关联
	// 它会在 中间表 sys_role_dept 插入对应的记录。
	//	•	不会重复插入（如果记录已存在）。
	//	•	不会更新 SysDept 表本身。
	//	•	只是建立 role 和 dept 的 关联关系。
	//db.Model(&role1).Association("SysDept").Append(deptsForRole1)

	//depts := []SysDept{{DeptId: 10}, {DeptId: 20}}
	//db.Model(&role1).Association("SysDept").Append(depts)
	fmt.Println("---------")
	var deptsForRole2 []SysDept
	db.Model(&role1).Association("SysDept").Find(&deptsForRole2)
	t.Logf("deptsForRole2: %+v", deptsForRole2)

	allRoles := []SysRole{}
	err := db.Preload("SysDept").Find(allRoles).Error
	require.Nil(t, err)
	t.Logf("allRoles: %+v", allRoles)
}

func initData() {
	db := getDB()

	roles := []SysRole{
		{RoleId: 1}, // 假设手动指定主键，或者留空让自增
		{RoleId: 2},
	}
	db.Create(&roles)

	depts := []SysDept{
		{DeptId: 10, ParentId: 0, DeptPath: "0/10", DeptName: "总部"},
		{DeptId: 20, ParentId: 10, DeptPath: "0/10/20", DeptName: "研发部"},
		{DeptId: 30, ParentId: 10, DeptPath: "0/10/30", DeptName: "销售部"},
	}
	db.Create(&depts)
}
