
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data *{{.upperStartCamelObject}}) error {
	{{if .withCache}}{{.keys}}
    _, err := m.dbConn.ExecCtx(ctx, func(ctx context.Context, conn *gorm.DB) (result interface{}, err error) {
		if err := conn.Updates(data).Error; err != nil {
		    return nil, perrors.WithStack(err)
		}
		return nil, nil
	}, {{.keyValues}})
	if err !=nil {
	    return err
	}{{else}}if err := conn.Updates(data).Error; err != nil {
        return perrors.WithStack(err)
    }{{end}}
    return nil
}
