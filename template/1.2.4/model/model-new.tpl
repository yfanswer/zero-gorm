
func New{{.upperStartCamelObject}}Model(gdb *gorm.DB, c cache.CacheConf, opts ...cache.Option) {{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		dbConn: db.NewDBConn(gdb, c, opts...),
		table:{{.table}},
	}
}