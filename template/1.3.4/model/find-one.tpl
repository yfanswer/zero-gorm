
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	var {{.lowerStartCamelObject}} {{.upperStartCamelObject}}
	{{if .withCache}}{{.cacheKey}}
	err := m.dbConn.QueryRowCtx(ctx, &{{.lowerStartCamelObject}}, {{.cacheKeyVariable}}, func(ctx context.Context, conn *gorm.DB, v interface{}) error {
		if err := conn.Where("`{{.lowerStartCamelPrimaryKey}}` = ?", {{.lowerStartCamelPrimaryKey}}).First(&{{.lowerStartCamelObject}}).Error; err != nil {
		    return perrors.WithStack(err)
		}
		return nil
	})
	if err != nil {
        return nil, err
    }
    return &{{.lowerStartCamelObject}}, nil{{else}}err := m.dbConn.db.Where("`{{.lowerStartCamelPrimaryKey}}` = ?", {{.lowerStartCamelPrimaryKey}}).First(&{{.lowerStartCamelObject}}).Error; err != nil {
            return perrors.WithStack(err)
        }
    })
    if err != nil {
        return nil, err
    }
    return &{{.lowerStartCamelObject}}, nil{{end}}
}
