
func (m *default{{.upperStartCamelObject}}Model) Update(data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
	err := m.dbConn.UpdateIndex(func(conn *gorm.DB) error {
		return conn.Updates(data).Error
	}, {{.keyValues}})
	if err != nil {
		return perrors.WithStack(err)
	}{{else}}err := m.dbConn.UpdateIndex(func(conn *gorm.DB) error {
        return conn.Updates(data).Error
    })
    if err != nil {
        return perrors.WithStack(err)
    }{{end}}
	return nil
}
