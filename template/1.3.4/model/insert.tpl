
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (interface{},error) {
	{{if .withCache}}{{.keys}}
    ret, err := m.dbConn.ExecCtx(ctx, func(ctx context.Context, conn *gorm.DB) (interface{}, error) {
        var res int64
        if err := conn.Create(data).Row().Scan(&res); err != nil {
            return nil, perrors.WithStack(err)
        }
		return res, nil
	}, {{.keyValues}})
	return ret, err{{else}}var ret {{.dataType}}
    if err := conn.Create(data).Row().Scan(&ret); err != nil {
      return nil, perrors.WithStack(err)
    }
	return ret,err{{end}}
}
