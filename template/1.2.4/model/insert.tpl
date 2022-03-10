
func (m *default{{.upperStartCamelObject}}Model) Insert(data *{{.upperStartCamelObject}}) error {
    {{if .withCache}}{{if .containsIndexCache}}{{.keys}}
    err := m.dbConn.InsertIndex(func(conn *gorm.DB) error {
        return conn.Create(data).Error
    },{{.keyValues}})
    if err != nil {
        return perrors.WithStack(err)
    }{{else}}err := m.dbConn.InsertIndex(func(conn *gorm.DB) error {
        return conn.Create(data).Error
    })
    if err != nil {
        return perrors.WithStack(err)
    }{{end}}{{else}}err := m.dbConn.InsertIndex(func(conn *gorm.DB, v interface{}) error {
        return conn.Create(data).Error
    })
    if err != nil {
        return perrors.WithStack(err)
    }{{end}}
    return nil
}
