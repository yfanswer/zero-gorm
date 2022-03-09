
func (m *default{{.upperStartCamelObject}}Model) FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	var resp {{.upperStartCamelObject}}

	{{if .withCache}}{{.cacheKey}}
	err := m.dbConn.FindIndex(&resp, {{.cacheKeyVariable}}, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).First(&resp).Error
	})
	if err != nil {
		return nil, perrors.WithStack(err)
	}{{else}}
	err := m.dbConn.FindNoCache(&resp, func(conn *gorm.DB, v interface{}) error {
        return conn.Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).First(&resp).Error
    })
    if err != nil {
        return nil, perrors.WithStack(err)
    }{{end}}
	return &resp, nil
}
