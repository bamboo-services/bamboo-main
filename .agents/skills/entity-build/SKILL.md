---
name: entity-build
description: >-
  根据口头描述快速生成符合 bamboo-base-go-template 项目规范的 Go Entity 代码。
  当用户说"创建一个实体"、"添加 entity"、"新建数据模型"、"加张表"、"建个 Model"
  或任何涉及创建 GORM 实体的场景时，都应使用此技能。即使用户只是说"加个 User 表"
  或"我需要存玩家数据"，也应触发。
---

# Entity 构建

口头描述实体需求，自动生成符合 bamboo-base-go-template 项目规范的 Go 实体代码。

实体是数据模型的单一真相来源——GORM 标签决定表结构，JSON 标签决定接口序列化，行尾注释是团队协作的语义说明。如果实体定义不规范，后续 handler/logic/repository 层都会跟着跑偏。所以这个技能的每条规范都指向一个目标：**让实体一次写对，后续不需要回头改**。

---

## 交互流程

当你说「创建一个实体」时，会按以下步骤进行：

1. **收集基本信息**
   - 实体名称（自动转为 PascalCase）
   - 实体描述（中文）

2. **询问字段定义**
   - 字段名、类型、是否可空
   - 特殊要求（唯一索引、默认值等）

3. **询问关联关系**（可选）
   - 是否属于其他实体（belongs_to）
   - 是否拥有多个子实体（has_many）

4. **生成代码**
   - 输出到 `internal/entity/<snake_case>.go`
   - 所有字段行尾追加中文注释（`// 中文说明`）
   - 提醒 Gene 常量定义和 AutoMigrate 注册

---

## 支持的字段类型

| 描述                       | Go 类型                    | 说明        |
|--------------------------|--------------------------|-----------|
| `string`                 | `string`                 | 字符串       |
| `int`                    | `int`                    | 整数        |
| `int64`                  | `int64`                  | 64 位整数    |
| `uint`                   | `uint`                   | 无符号整数     |
| `xSnowflake.SnowflakeID` | `xSnowflake.SnowflakeID` | 雪花算法 ID   |
| `bool`                   | `bool`                   | 布尔值       |
| `float`                  | `float64`                | 浮点数       |
| `time`                   | `time.Time`              | 时间戳       |
| `decimal`                | `float64`                | 小数        |
| 可空字符串                    | `*string`                | 指针类型，允许 null |

---

## 使用 AskUserQuestion 收集信息

当用户的描述不够完整时，主动询问。以下是一些常见场景：

```yaml
# 询问字段类型
questions:
  - question: "UUID 字段需要唯一约束吗？"
    header: "唯一约束"
    options:
      - label: "是，唯一"
        description: "添加 unique 约束，防止重复"
      - label: "否，可重复"
        description: "允许相同值存在"
    multiSelect: false

# 询问是否可空
questions:
  - question: "LastSeen 字段是否可空？"
    header: "可空类型"
    options:
      - label: "可空"
        description: "使用 *time.Time 指针类型"
      - label: "不可空"
        description: "使用 time.Time 类型"
    multiSelect: false

# 询问关联关系
questions:
  - question: "Player 需要关联哪些实体？"
    header: "关联关系"
    options:
      - label: "属于 User"
        description: "添加 UserID 外键，属于一个用户"
      - label: "拥有多个 GameProfile"
        description: "一对多关系"
    multiSelect: true
```

| 情况       | 询问内容                              |
|----------|-----------------------------------|
| 字段类型不明确  | 确认 Go 类型（string/int/bool 等）       |
| 字段约束不明确  | 确认是否 unique、not null、默认值          |
| 关系不明确    | 确认是否属于其他实体、是否有一对多关系               |
| Gene 不明确 | 确认使用内置 Gene 还是自定义                 |

能从上下文推断的信息直接判断，不询问用户。

---

## 常用字段模板

| 场景     | GORM 标签                                                 | JSON 标签                       |
|--------|---------------------------------------------------------|-------------------------------|
| 非空字符串  | `gorm:"not null;type:varchar(255);comment:说明"`          | `json:"field_name"`           |
| 可空字符串  | `gorm:"type:varchar(512);comment:说明"`                   | `json:"field_name,omitempty"` |
| 唯一字符串  | `gorm:"unique;not null;type:varchar(36);comment:说明"`    | `json:"field_name"`           |
| 整数     | `gorm:"not null;default:1;comment:说明"`                  | `json:"field_name"`           |
| 布尔     | `gorm:"not null;type:boolean;default:false;comment:说明"` | `json:"field_name"`           |
| 时间戳    | `gorm:"type:timestamptz;comment:说明"`                    | `json:"field_name,omitempty"` |
| 外键     | `gorm:"not null;index:idx_user_id;comment:说明"`          | `json:"user_id"`              |
| 密码（敏感） | `gorm:"not null;type:varchar(255);comment:说明"`          | `json:"-"`                    |

---

## 字段行尾注释

生成实体时，所有字段行尾都要有中文注释。行尾注释是团队协作中最快速的字段语义说明——看代码时不需要跳转到 GORM 标签或数据库注释去理解字段含义。

```go
FieldName FieldType `gorm:"...;comment:字段说明" json:"field_name"` // 字段说明
```

规则：

1. 结构体中每一行字段定义都要有行尾注释（包括普通字段、外键字段、切片关联字段）。
2. 行尾注释语义和字段含义一致，与 `gorm comment` 保持一致。
3. 行尾注释使用中文，格式统一为 `// 中文说明`。
4. 不省略行尾注释，即使字段名看起来很清晰——团队约定统一风格，不逐个判断"这个够不够明显"。

---

## 外键关系模板

### belongs_to（多对一）

```go
UserID xSnowflake.SnowflakeID `gorm:"not null;index:idx_user_id;comment:关联用户ID" json:"user_id"` // 关联用户ID
User   User                   `gorm:"constraint:OnDelete:CASCADE;comment:关联用户" json:"user,omitempty"` // 关联用户
```

### has_many（一对多）

```go
GameProfiles []GameProfile `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;comment:游戏档案关联" json:"game_profiles,omitempty"` // 游戏档案关联
```

---

## GetGene 方法

每个实体（除 `Role` 等使用字符串主键的特例外）都要实现 `GetGene()` 方法，返回该实体在雪花算法 ID 生成时使用的基因类型。基因类型决定了 ID 的命名空间隔离——不同实体类型的 ID 不会冲突。

```go
// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *EntityName) GetGene() xSnowflake.Gene {
    return xSnowflake.GeneUser       // 内置类型
    // return bConst.GeneForXXX      // 自定义类型（需要在 constant 中定义）
}
```

### 常用 Gene 类型

| Gene 值                  | 用途          |
|-------------------------|-------------|
| `xSnowflake.GeneUser`   | 用户实体        |
| `xSnowflake.GeneDefault`| 默认/通用实体     |
| `bConst.GeneForXXX`     | 自定义（需在 constant 中定义） |

`Role` 实体是例外：它使用 `Name`（字符串主键）而非雪花 ID，所以不需要 `GetGene()`。如果新实体也是字符串主键，同样跳过此方法。

---

## 完整生成示例

### 用户输入

```
创建一个 Player 实体，包含：
- UUID（唯一）
- Name（游戏内玩家名）
- Level（等级，默认1）
- LastSeen（最后在线时间，可空）
属于 User
```

### 生成结果

**文件**: `internal/entity/player.go`

```go
package entity

import (
	"time"

	bConst "github.com/xiaolfeng/coding-plan-manage/internal/constant"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xModels "github.com/bamboo-services/bamboo-base-go/major/models"
)

// Player 玩家实体，包含 UUID、名称、等级等游戏内信息。
type Player struct {
	xModels.BaseEntity                        // 嵌入基础实体字段（雪花ID + 时间戳）
	UserID             xSnowflake.SnowflakeID `gorm:"not null;index:idx_user_id;comment:关联用户ID" json:"user_id"` // 关联用户ID
	UUID               string                 `gorm:"unique;not null;type:varchar(36);comment:Minecraft UUID" json:"uuid"` // Minecraft UUID
	Name               string                 `gorm:"not null;type:varchar(32);comment:游戏内玩家名" json:"name"` // 游戏内玩家名
	Level              int                    `gorm:"not null;default:1;comment:玩家等级" json:"level"` // 玩家等级
	LastSeen           *time.Time             `gorm:"type:timestamptz;comment:最后在线时间" json:"last_seen,omitempty"` // 最后在线时间

	// ----------
	//  外键约束
	// ----------
	User User `gorm:"constraint:OnDelete:CASCADE;comment:关联用户" json:"user,omitempty"` // 关联用户
}

// GetGene 返回 xSnowflake.Gene，用于标识该实体在 ID 生成时使用的基因类型。
func (_ *Player) GetGene() xSnowflake.Gene {
	return bConst.GeneForPlayer // 需要在 internal/constant/gene_number.go 中定义
}
```

**提醒**: 记得在 `internal/constant/gene_number.go` 中添加 Gene 常量：

```go
const (
    GeneForPlayer xSnowflake.Gene = 64 // 玩家实体
)
```

并注册到 `main.go` 的 `xOption.WithDatabase(...)` 中（按外键依赖顺序追加实体指针）：

```go
xOption.WithDatabase(
    xOptionDB.FromEnv(),
    xOptionDB.WithAutoMigrate(
        &entity.User{},
        &entity.Role{},
        &entity.Player{},  // 新增
    ),
)

---

## 注意事项

1. **Gene 常量**: 自定义 Gene 需要在 `internal/constant/gene_number.go` 中定义，并确保编号不与已有常量冲突。
2. **外键删除策略**: 默认使用 `OnDelete:CASCADE`，如果业务要求数据保留则改为 `RESTRICT`（如 `Role` 实体）。
3. **字段行尾注释**: 所有字段定义追加 `// 中文说明`，不省略。
4. **指针类型**: 可空字段使用指针类型（`*string`、`*time.Time`），JSON 标签自动添加 `omitempty`。
5. **敏感字段**: 密码等使用 `json:"-"` 隐藏，不暴露到接口。
6. **不确定时**: 使用 `AskUserQuestion` 询问用户，不擅自猜测——实体是数据模型的根，猜错了后续全跟着错。
7. **新表注册**: 新建实体后，追加到 `main.go` 的 `xOptionDB.WithAutoMigrate(...)` 中（按外键依赖顺序），否则不会自动建表。

---

## 参考资料

- **bamboo-base-go 全局文档**: https://doc.x-lf.com/llms.txt
- **具体路径查询**: https://doc.x-lf.com/llms.mdx/<search_path>

### 查询示例

| 需要查找的内容           | 查询 URL                                              |
|-------------------|-----------------------------------------------------|
| BaseEntity 定义     | https://doc.x-lf.com/llms.mdx/models/base_entity.go |
| Snowflake Gene 类型 | https://doc.x-lf.com/llms.mdx/snowflake/gene.go     |
| 所有可导出类型           | https://doc.x-lf.com/llms.txt                       |
