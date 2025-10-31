package repositories

import (
	"chatbot/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileAnalysisRepository handles database operations for file analysis
type FileAnalysisRepository struct {
	db *gorm.DB
}

// NewFileAnalysisRepository creates a new file analysis repository
func NewFileAnalysisRepository(db *gorm.DB) *FileAnalysisRepository {
	return &FileAnalysisRepository{db: db}
}

// Create saves a new file analysis record
func (r *FileAnalysisRepository) Create(analysis *models.FileAnalysis) error {
	return r.db.Create(analysis).Error
}

// GetByID retrieves a file analysis by ID
func (r *FileAnalysisRepository) GetByID(id uuid.UUID) (*models.FileAnalysis, error) {
	var analysis models.FileAnalysis
	err := r.db.Where("id = ?", id).First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// GetAll retrieves all file analyses with pagination
func (r *FileAnalysisRepository) GetAll(limit, offset int) ([]models.FileAnalysis, int64, error) {
	var analyses []models.FileAnalysis
	var total int64

	// Count total records
	if err := r.db.Model(&models.FileAnalysis{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// GetByFileType retrieves file analyses filtered by file type
func (r *FileAnalysisRepository) GetByFileType(fileType string, limit, offset int) ([]models.FileAnalysis, int64, error) {
	var analyses []models.FileAnalysis
	var total int64

	query := r.db.Model(&models.FileAnalysis{}).Where("file_type = ?", fileType)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// Update updates a file analysis record
func (r *FileAnalysisRepository) Update(analysis *models.FileAnalysis) error {
	return r.db.Save(analysis).Error
}

// Delete soft deletes a file analysis record
func (r *FileAnalysisRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.FileAnalysis{}, "id = ?", id).Error
}

// IncrementReanalysisCount increments the reanalysis counter
func (r *FileAnalysisRepository) IncrementReanalysisCount(id uuid.UUID) error {
	return r.db.Model(&models.FileAnalysis{}).
		Where("id = ?", id).
		UpdateColumn("reanalysis_count", gorm.Expr("reanalysis_count + 1")).
		Error
}

// GetRecentAnalyses retrieves the most recent file analyses
func (r *FileAnalysisRepository) GetRecentAnalyses(limit int) ([]models.FileAnalysis, error) {
	var analyses []models.FileAnalysis
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Find(&analyses).Error

	if err != nil {
		return nil, err
	}

	return analyses, nil
}

// SearchByFileName searches file analyses by filename (partial match)
func (r *FileAnalysisRepository) SearchByFileName(filename string, limit, offset int) ([]models.FileAnalysis, int64, error) {
	var analyses []models.FileAnalysis
	var total int64

	query := r.db.Model(&models.FileAnalysis{}).Where("file_name ILIKE ?", "%"+filename+"%")

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// GetStatistics returns analysis statistics
func (r *FileAnalysisRepository) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total analyses
	var totalCount int64
	if err := r.db.Model(&models.FileAnalysis{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_analyses"] = totalCount

	// Count by file type
	type FileTypeCount struct {
		FileType string
		Count    int64
	}
	var fileTypeCounts []FileTypeCount
	if err := r.db.Model(&models.FileAnalysis{}).
		Select("file_type, COUNT(*) as count").
		Group("file_type").
		Scan(&fileTypeCounts).Error; err != nil {
		return nil, err
	}
	stats["by_file_type"] = fileTypeCounts

	// Average processing time
	var avgProcessTime float64
	if err := r.db.Model(&models.FileAnalysis{}).
		Select("AVG(process_time_ms) as avg_time").
		Scan(&avgProcessTime).Error; err != nil {
		return nil, err
	}
	stats["avg_process_time_ms"] = avgProcessTime

	// Total tokens used
	var totalTokens int64
	if err := r.db.Model(&models.FileAnalysis{}).
		Select("SUM(tokens_used) as total").
		Scan(&totalTokens).Error; err != nil {
		return nil, err
	}
	stats["total_tokens_used"] = totalTokens

	return stats, nil
}