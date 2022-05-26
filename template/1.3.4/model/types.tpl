
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
        dbConn *db.DBConn
        table string
    }

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
