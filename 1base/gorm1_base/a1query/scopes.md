# GORM Scopes 详细说明

## 什么是 Scopes

Scopes 是 GORM 提供的一个强大功能，它允许你定义常用的查询逻辑作为可重用的函数，然后将这些函数链式地应用到查询中。本质上，Scopes 是一种封装和重用查询逻辑的方法，有助于保持代码的干净和可维护性。

## 基本语法

Scope 函数的标准签名：

```go
func(db *gorm.DB) *gorm.DB
```

每个 Scope 函数接收一个 `*gorm.DB` 对象，对其应用某些查询条件，然后返回修改后的 `*gorm.DB` 对象。

## 使用场景

### 1. 定义常用查询条件

```go
// 定义 Scope 函数
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func OrderedByCreatedAt(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC")
}

func RecentlyRegistered(db *gorm.DB) *gorm.DB {
    return db.Where("created_at > ?", time.Now().Add(-7*24*time.Hour))
}

// 使用 Scope 函数
db.Scopes(ActiveUsers, RecentlyRegistered).Find(&users)
// 等同于: db.Where("active = ?", true).Where("created_at > ?", time.Now().Add(-7*24*time.Hour)).Find(&users)
```

### 2. 分页查询

```go
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

// 使用分页 Scope
db.Scopes(Paginate(2, 10)).Find(&users)
// 等同于: db.Offset(10).Limit(10).Find(&users)
```

### 3. 动态查询条件

```go
func Filter(filter map[string]interface{}) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        for key, value := range filter {
            if value != nil {
                db = db.Where(key+" = ?", value)
            }
        }
        return db
    }
}

// 使用动态过滤
filter := map[string]interface{}{
    "status": "active",
    "role":   "admin",
}
db.Scopes(Filter(filter)).Find(&users)
```

### 4. 软删除过滤

```go
func NotDeleted(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL")
}

// 使用软删除过滤
db.Scopes(NotDeleted).Find(&users)
```

### 5. 权限控制

```go
func VisibleToUser(userID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("public = ? OR user_id = ?", true, userID)
    }
}

// 使用权限过滤
db.Scopes(VisibleToUser(currentUserID)).Find(&posts)
```

## 高级用法

### 1. 组合多个 Scopes

```go
// 组合使用多个 Scopes
db.Scopes(ActiveUsers, RecentlyRegistered, OrderedByCreatedAt).Find(&users)
```

### 2. 条件 Scopes

```go
func ConditionalScope(condition bool, trueScope, falseScope func(*gorm.DB) *gorm.DB) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if condition {
            return trueScope(db)
        }
        return falseScope(db)
    }
}

// 使用条件 Scope
isAdmin := true
db.Scopes(ConditionalScope(isAdmin, AdminView, RegularUserView)).Find(&data)
```

### 3. 带参数的 Scopes

```go
func WithinDateRange(start, end time.Time) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("created_at BETWEEN ? AND ?", start, end)
    }
}

// 使用带参数的 Scope
startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
db.Scopes(WithinDateRange(startDate, endDate)).Find(&records)
```

## Scopes 的优势

1. **代码复用**：避免在多个地方重复编写相同的查询逻辑
2. **可读性**：使代码更加清晰和自文档化
3. **可维护性**：集中修改查询逻辑，而不必更新所有使用该逻辑的地方
4. **模块化**：将复杂查询拆分为更小、更易管理的部分
5. **测试性**：可以单独测试每个 Scope 函数

## 最佳实践

1. **命名明确**：为 Scope 函数使用描述性名称，清晰表达其功能
2. **单一职责**：每个 Scope 应该只负责一个特定的查询条件
3. **文档化**：为复杂的 Scope 函数添加注释，解释其作用和用法
4. **避免副作用**：Scope 函数不应修改传入的参数或执行查询操作
5. **按模型组织**：将特定于模型的 Scope 函数放在相应的模型文件中

## 与 Where 子句的区别

虽然简单的 Scope 可能看起来只是 Where 子句的包装，但 Scopes 提供了更高级的抽象和组合能力。它们允许封装完整的查询逻辑，而不仅仅是条件。

## 总结

GORM 的 Scopes 功能是一个强大的工具，用于组织和重用数据库查询逻辑。通过将常用查询封装为 Scope 函数，可以显著提高代码的可读性、可维护性和可测试性，同时减少重复代码。在复杂应用中，合理使用 Scopes 可以使数据访问层更加清晰和结构化。