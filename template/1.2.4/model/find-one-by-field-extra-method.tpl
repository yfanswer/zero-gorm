
func (m *default{{.upperStartCamelObject}}Model) FormatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", {{.primaryKeyLeft}}, primary)
}
