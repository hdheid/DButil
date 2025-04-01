该项目简单实现了一个数据库的CURD代码生成。

## 使用方式

该程序生成的代码一般是在 model 层下，存放表结构的结构体的代码下。进入到该层后，输入命令 `curd curd <filename>` 即可生成对应代码。需要注意的是该代码需要有 `TableName()` 方法以及 `gorm` 库。

### 示例

例如我的 model 层下有一个 doc_field.go 文件：
```go
package model

import (
	"time"
)

type LearningEcDocField struct {
	Id         int64     `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`                         // 自增ID
	IsDeleted  int       `gorm:"column:is_deleted;default:0;NOT NULL" json:"is_deleted"`                 // 是否删除：0-正常，1-删除
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 更新时间
}

func (m *LearningEcDocField) TableName() string {
	return "ec_doc_field"
}
```

在 model 层下输入 `curd curd doc_field.go`，即可生成以下代码：
```go

/*
请确保导入 gorm.io/gorm 库
*/

const (
	IsDeletedNo = iota
	IsDeletedYes
)

// InsertOne ...
func (m *LearningEcDocField) InsertOne(tx *gorm.DB, row *LearningEcDocField) (*LearningEcDocField, error) {
	if err := tx.Table(m.TableName()).Create(&row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

// InsertMulti ...
func (m *LearningEcDocField) InsertMulti(tx *gorm.DB, row []*LearningEcDocField) ([]*LearningEcDocField, error) {
	if err := tx.Table(m.TableName()).Create(&row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

// FindOne ...
func (m *LearningEcDocField) FindOne(tx *gorm.DB, where map[string]interface{}) (*LearningEcDocField, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	oneData := new(LearningEcDocField)
	if err := tx.Table(m.TableName()).Where("is_deleted = ?", IsDeletedNo).First(oneData).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, err
		default:
			return nil, err
		}
	}
	return oneData, nil
}

// FindOneWithOrder ...
func (m *LearningEcDocField) FindOneWithOrder(tx *gorm.DB, where map[string]interface{}, orderBy string) (*LearningEcDocField, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	oneData := new(LearningEcDocField)
	if err := tx.Table(m.TableName()).Where("is_deleted = ?", IsDeletedNo).Order(orderBy).First(oneData).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, err
		default:
			return nil, err
		}
	}
	return oneData, nil
}

// FindMulti ...
func (m *LearningEcDocField) FindMulti(tx *gorm.DB, where map[string]interface{}, orderBy string) ([]*LearningEcDocField, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	var rows []*LearningEcDocField
	if err := tx.Table(m.TableName()).Where("is_deleted = ?", IsDeletedNo).Order(orderBy).Find(&rows).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, err
		default:
			return nil, err
		}
	}
	return rows, nil
}

// FindList ...
func (m *LearningEcDocField) FindList(tx *gorm.DB, where map[string]interface{}, page, pageSize int, orderBy string) ([]*LearningEcDocField, int64, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	tx = tx.Table(m.TableName()).Where("is_deleted = ?", IsDeletedNo)

	if len(orderBy) > 0 {
		tx = tx.Order(orderBy)
	}

	var total int64
	tx.Count(&total)

	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize > 500:
		pageSize = 500
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var rows []*LearningEcDocField
	if err := tx.Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, 0, err
		default:
			return nil, 0, err
		}
	}
	return rows, total, nil
}

// UpdateData ...
func (m *LearningEcDocField) UpdateData(tx *gorm.DB, where map[string]interface{}, data map[string]interface{}) (int64, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	query := tx.Table(m.TableName()).Where("is_deleted = ?", IsDeletedNo).Updates(data)
	if err := query.Error; err != nil {
		return 0, err
	}
	return query.RowsAffected, nil
}

type UpdatesMulti struct {
	IdVal     string
	FieldName string
	FieldVal  string
}

// UpdateMultiData ...
func (m *LearningEcDocField) UpdateMultiData(tableName string, tx *gorm.DB, idName string, update []*UpdatesMulti) (int64, error) {
	if update == nil || len(update) <= 0 {
		return 0, nil
	}

	var FieldMap = make(map[string][]*UpdatesMulti)
	var fieldNameList []string
	idsMp := make(map[string]struct{})
	var ids []string
	var args []interface{}

	for _, u := range update {
		if _, ok := FieldMap[u.FieldName]; !ok {
			FieldMap[u.FieldName] = make([]*UpdatesMulti, 0)
			fieldNameList = append(fieldNameList, u.FieldName)
		}

		FieldMap[u.FieldName] = append(FieldMap[u.FieldName], u)
		idsMp[u.IdVal] = struct{}{}
	}

	for k := range idsMp {
		ids = append(ids, k)
	}

	sql := fmt.Sprintf("UPDATE %s SET", tableName)

	for i, fieldName := range fieldNameList {
		var str strings.Builder
		var set string

		for _, v := range FieldMap[fieldName] {
			str.WriteString("WHEN ? THEN ? ")
			args = append(args, v.IdVal, v.FieldVal)
		}

		if i == len(fieldNameList)-1 {
			set = fmt.Sprintf("%s = CASE %s %s ELSE %s END", fieldName, idName, str.String(), fieldName)
		} else {
			set = fmt.Sprintf("%s = CASE %s %s ELSE %s END,", fieldName, idName, str.String(), fieldName)
		}

		sql = fmt.Sprintf("%s %s", sql, set)
	}

	sql = fmt.Sprintf("%s WHERE %s IN (%s) AND is_deleted = 0",
		sql,
		idName,
		strings.Join(strings.Split(strings.Repeat("?", len(ids)), ""), ","),
	)

	for _, id := range ids {
		args = append(args, id)
	}

	query := tx.Exec(sql, args...)
	if query.Error != nil {
		return 0, query.Error
	}
	return query.RowsAffected, nil
}
```
