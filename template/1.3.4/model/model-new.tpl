
func new{{.upperStartCamelObject}}Model(gdb *gorm.DB, c cache.CacheConf) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		dbConn: db.NewDBConn(gdb, c),
		table:{{.table}},
	}
}
