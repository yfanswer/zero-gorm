
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}(ctx context.Context, {{.in}}) (*{{.upperStartCamelObject}}, error) {
	var resp {{.upperStartCamelObject}}
	{{if .withCache}}{{.cacheKey}}
	err := m.dbConn.QueryRowIndexCtx(ctx, &resp, {{.cacheKeyVariable}}, m.formatPrimary, func(ctx context.Context, conn *gorm.DB, v interface{}) (i interface{}, e error) {
		if err := conn.Where("{{.originalField}}", {{.lowerStartCamelField}}).First(&resp).Error; err != nil {
		    return nil, perrors.WithStack(err)
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	if err != nil {
	    return nil, err
	}
	return &resp, nil
}{{else}}if err := conn.Where("{{.originalField}}", {{.lowerStartCamelField}}).First(&resp).Error; err != nil {
        return nil, perrors.WithStack(err)
    }
    return &resp, nil
}{{end}}
