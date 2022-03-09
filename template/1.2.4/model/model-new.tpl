
func New{{.upperStartCamelObject}}Model(conn db.DBConn) {{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		dbConn: conn,
		table:{{.table}},
	}
}
