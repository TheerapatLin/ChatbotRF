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

// FindByID retrieves a file analysis by ID string
func (r *FileAnalysisRepository) FindByID(id string) (*models.FileAnalysis, error) {
	fileID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return r.GetByID(fileID)
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
	err := r.db.Order("uploaded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// GetByFileType retrieves file analyses filtered by MIME type
func (r *FileAnalysisRepository) GetByFileType(mimeType string, limit, offset int) ([]models.FileAnalysis, int64, error) {
	var analyses []models.FileAnalysis
	var total int64

	query := r.db.Model(&models.FileAnalysis{}).Where("mime_type = ?", mimeType)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := query.Order("uploaded_at DESC").
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

// GetRecentAnalyses retrieves the most recent file uploads
func (r *FileAnalysisRepository) GetRecentAnalyses(limit int) ([]models.FileAnalysis, error) {
	var analyses []models.FileAnalysis
	err := r.db.Order("uploaded_at DESC").
		Limit(limit).
		Find(&analyses).Error

	if err != nil {
		return nil, err
	}

	return analyses, nil
}

// SearchByFileName searches file uploads by filename (partial match)
func (r *FileAnalysisRepository) SearchByFileName(filename string, limit, offset int) ([]models.FileAnalysis, int64, error) {
	var analyses []models.FileAnalysis
	var total int64

	query := r.db.Model(&models.FileAnalysis{}).Where("file_name ILIKE ?", "%"+filename+"%")

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	err := query.Order("uploaded_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&analyses).Error

	if err != nil {
		return nil, 0, err
	}

	return analyses, total, nil
}

// GetStatistics returns file upload statistics
func (r *FileAnalysisRepository) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total uploads
	var totalCount int64
	if err := r.db.Model(&models.FileAnalysis{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_uploads"] = totalCount

	// Count by MIME type
	type MimeTypeCount struct {
		MimeType string
		Count    int64
	}
	var mimeTypeCounts []MimeTypeCount
	if err := r.db.Model(&models.FileAnalysis{}).
		Select("mime_type, COUNT(*) as count").
		Group("mime_type").
		Scan(&mimeTypeCounts).Error; err != nil {
		return nil, err
	}
	stats["by_mime_type"] = mimeTypeCounts

	// Total storage size
	var totalSize int64
	if err := r.db.Model(&models.FileAnalysis{}).
		Select("SUM(file_size) as total").
		Scan(&totalSize).Error; err != nil {
		return nil, err
	}
	stats["total_storage_bytes"] = totalSize

	return stats, nil
}