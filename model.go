package main

//
//func InitDB() (*gorm.DB, error) {
//	dsn := "root:Wuwang222@tcp(127.0.0.1:3306)/db_ex_learning_condition?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info), // 打印所有 SQL 日志
//	})
//	if err != nil {
//		return nil, err
//	}
//	return db, nil
//}
//
//func main() {
//	tx, err := InitDB()
//
//	var tableName = "t_learning_ec_doc_fields"
//	var idName = "field_id"
//	u := make([]*UpdatesMulti, 0)
//	u = append(u, &UpdatesMulti{
//		IdVal:     "1531046071",
//		FieldName: "field_name",
//		FieldVal:  "单选框1-test",
//	})
//	u = append(u, &UpdatesMulti{
//		IdVal:     "5427147458",
//		FieldName: "field_name",
//		FieldVal:  "复选框1-test",
//	})
//	u = append(u, &UpdatesMulti{
//		IdVal:     "1531046071",
//		FieldName: "field_type",
//		FieldVal:  "text_single_line",
//	})
//	u = append(u, &UpdatesMulti{
//		IdVal:     "1302044886",
//		FieldName: "field_name",
//		FieldVal:  "订单验证-test",
//	})
//	u = append(u, &UpdatesMulti{
//		IdVal:     "6218336788",
//		FieldName: "field_value",
//		FieldVal:  "666",
//	})
//
//	_, err = UpdateMultiData(tableName, tx, idName, u)
//	if err != nil {
//		return
//	}
//}
