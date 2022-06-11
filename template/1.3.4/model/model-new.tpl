
func new{{.upperStartCamelObject}}Model(conn *db.DBConn) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		dbConn: conn,
		table:{{.table}},
	}
}
