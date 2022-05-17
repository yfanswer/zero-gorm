
func (m *default{{.upperStartCamelObject}}Model) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}

func (m *default{{.upperStartCamelObject}}Model) queryPrimary(ctx context.Context, conn *gorm.DB, v, primary interface{}) error {
	if err := conn.Where("{{.originalPrimaryField}} = ?", primary).First(v).Error; err != nil {
        return perrors.WithStack(err)
    }
    return nil
}
