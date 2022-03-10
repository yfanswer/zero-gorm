
func (m *default{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
    if err!=nil{
        return err
    }{{end}}
    
	{{.keys}}{{end}}
	err = m.dbConn.DelIndex(func(conn *gorm.DB) error {
		return conn.Delete(&{{.upperStartCamelObject}}{},{{.lowerStartCamelPrimaryKey}}).Error
	}{{if .withCache}},{{.keyValues}}{{end}})
	if err != nil {
		return perrors.WithStack(err)
	}
	return nil
}
