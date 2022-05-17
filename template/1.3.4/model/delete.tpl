
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}
{{end}}	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.dbConn.ExecCtx(ctx, func(ctx context.Context, conn *gorm.DB) (result interface{}, err error) {
		if err := conn.Delete(&{{.upperStartCamelObject}}{}, {{.lowerStartCamelPrimaryKey}}).Error; err != nil {
		    return nil, perrors.WithStack(err)
		}
		return result, nil
	}, {{.keyValues}}){{end}}
	return err
}
