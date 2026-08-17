package repository

import (
	"errors"

	"github.com/rxcheck/rxcheck/internal/model"
	"gorm.io/gorm"
)

// ReportRepository 审核报告仓储。
type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// DB 返回底层数据库句柄。
func (r *ReportRepository) DB() *gorm.DB { return r.db }

// CreateWithItemsTx 在事务中创建报告及审核意见。
func (r *ReportRepository) CreateWithItemsTx(tx *gorm.DB, report *model.ReviewReport) error {
	items := report.Items
	report.Items = nil // 避免 GORM 自动保存 has-many 关联导致重复插入
	if err := tx.Create(report).Error; err != nil {
		return err
	}
	for i := range items {
		items[i].ReportID = report.ID
	}
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
	}
	report.Items = items
	return nil
}

// ReplaceWithItemsTx 在事务中删除旧报告并写入新报告（重新审核时使用）。
func (r *ReportRepository) ReplaceWithItemsTx(tx *gorm.DB, report *model.ReviewReport) error {
	if err := tx.Where("report_id IN (SELECT id FROM review_reports WHERE prescription_id = ?)", report.PrescriptionID).
		Delete(&model.ReviewItem{}).Error; err != nil {
		return err
	}
	if err := tx.Where("prescription_id = ?", report.PrescriptionID).Delete(&model.ReviewReport{}).Error; err != nil {
		return err
	}
	return r.CreateWithItemsTx(tx, report)
}

// FindByPrescriptionID 按处方查询报告（含审核意见）。
func (r *ReportRepository) FindByPrescriptionID(prescriptionID uint) (*model.ReviewReport, error) {
	var report model.ReviewReport
	err := r.db.Where("prescription_id = ?", prescriptionID).First(&report).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &report, err
}

// List 分页查询报告，支持审核结论过滤。
func (r *ReportRepository) List(page, pageSize int, status string) ([]model.ReviewReport, int64, error) {
	var list []model.ReviewReport
	var total int64
	q := r.db.Model(&model.ReviewReport{}).Preload("Items").Preload("Prescription")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return list, total, err
}

// CountByStatus 按审核结论统计。
func (r *ReportRepository) CountByStatus() (map[string]int64, error) {
	return groupCount(r.db, &model.ReviewReport{}, "status")
}
