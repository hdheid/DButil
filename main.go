package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"
)

const crudTemplate = `
/*
请确保导入 gorm.io/gorm 库
*/

const (
	IsDeletedNo = iota
	IsDeletedYes
)

// InsertOne ...
func (m *{{.ModelName}}) InsertOne(tx *gorm.DB, row *{{.ModelName}}) (*{{.ModelName}}, error) {
	if err := tx.Table(m.TableName()).Create(&row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

// InsertMulti ...
func (m *{{.ModelName}}) InsertMulti(tx *gorm.DB, row []*{{.ModelName}}) ([]*{{.ModelName}}, error) {
	if err := tx.Table(m.TableName()).Create(&row).Error; err != nil {
		return nil, err
	}
	return row, nil
}

// FindOne ...
func (m *{{.ModelName}}) FindOne(tx *gorm.DB, where map[string]interface{}) (*{{.ModelName}}, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	oneData := new({{.ModelName}})
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
func (m *{{.ModelName}}) FindOneWithOrder(tx *gorm.DB, where map[string]interface{}, orderBy string) (*{{.ModelName}}, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	oneData := new({{.ModelName}})
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
func (m *{{.ModelName}}) FindMulti(tx *gorm.DB, where map[string]interface{}, orderBy string) ([]*{{.ModelName}}, error) {
	for k, v := range where {
		tx = tx.Where(k, v)
	}
	var rows []*{{.ModelName}}
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
func (m *{{.ModelName}}) FindList(tx *gorm.DB, where map[string]interface{}, page, pageSize int, orderBy string) ([]*{{.ModelName}}, int64, error) {
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

	var rows []*{{.ModelName}}
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
func (m *{{.ModelName}}) UpdateData(tx *gorm.DB, where map[string]interface{}, data map[string]interface{}) (int64, error) {
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
func (m *{{.ModelName}}) UpdateMultiData(tableName string, tx *gorm.DB, idName string, update []*UpdatesMulti) (int64, error) {
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
`

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("正确命令为: curd gen <filename>")
	}

	filename := args[1]
	if filepath.Ext(filename) != ".go" {
		log.Fatal("请在文件所在目录使用该命令或输入正确文件名")
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var modelName string
	var found bool

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		if fn.Name.Name == "TableName" {
			if recv := fn.Recv; recv != nil && len(recv.List) > 0 {
				if star, ok := recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						modelName = ident.Name
						found = true
						return false
					}
				}
			}
		}
		return true
	})

	if !found {
		fmt.Println("没有找到 TableName() 函数")
		return
	}

	tmpl, err := template.New("crud").Parse(crudTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct{ ModelName string }{ModelName: modelName})
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	newContent := append(content, []byte(buf.String())...)
	err = ioutil.WriteFile(filename, newContent, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("成功为 %s 生成CRUD方法\n", modelName)
}
