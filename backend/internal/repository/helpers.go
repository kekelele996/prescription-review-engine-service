package repository

import "gorm.io/gorm"

// groupCount 按字段分组统计（被处方/报告统计接口复用）。
func groupCount(db *gorm.DB, model any, field string) (map[string]int64, error) {
	type row struct {
		GroupKey string
		Count    int64
	}
	var rows []row
	err := db.Model(model).Select(field + " AS group_key, COUNT(*) AS count").Group(field).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(rows))
	for _, rw := range rows {
		out[rw.GroupKey] = rw.Count
	}
	return out, nil
}
