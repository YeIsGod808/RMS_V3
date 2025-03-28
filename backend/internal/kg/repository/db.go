package repository

import (
	"fmt"

	"github.com/RMS_V3/internal/kg/repository/models"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/db"
)

func AddVideo(pointId int64, title string, videoUrl string, coverUrl string, description *string) error {
	db := db.GetDB()
	video := models.Video{
		Title:            title,
		PlayURL:          videoUrl,
		CoverURL:         coverUrl,
		Description:      description,
		KnowledgePointID: pointId,
	}
	return db.Create(&video).Error
}
func DeleteVideo(knowledgePointId int64, videoId int64) error {
	db := db.GetDB()
	// 构建查询条件
	conditions := map[string]interface{}{
		"knowledge_point_id": knowledgePointId,
		"id":                 videoId,
	}

	// 执行删除操作
	result := db.Where(conditions).Delete(&models.Video{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到符合条件的视频记录")
	}

	return nil
}
func AddCourseware(pointId int64, title string, coursewareUrl string, description *string) error {
	db := db.GetDB()
	courseware := models.Courseware{
		Title:            title,
		CoursewareURL:    coursewareUrl,
		Description:      description,
		KnowledgePointID: pointId,
	}
	return db.Create(&courseware).Error
}
func DeleteCourseware(knowledgePointId int64, coursewareId int64) error {
	db := db.GetDB()
	// 构建查询条件
	conditions := map[string]interface{}{
		"knowledge_point_id": knowledgePointId,
		"id":                 coursewareId,
	}

	// 执行删除操作
	result := db.Where(conditions).Delete(&models.Courseware{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到符合条件的视频记录")
	}

	return nil
}
func AddExercise(pointId int64, title string, exerciseUrl string, difficulty string, description *string) error {
	db := db.GetDB()
	exercise := models.Exercise{
		Title:            title,
		ExerciseURL:      exerciseUrl,
		Difficulty:       difficulty,
		Description:      description,
		KnowledgePointID: pointId,
	}
	return db.Create(&exercise).Error
}
func DeleteExercise(knowledgePointId int64, exerciseId int64) error {
	db := db.GetDB()
	// 构建查询条件
	conditions := map[string]interface{}{
		"knowledge_point_id": knowledgePointId,
		"id":                 exerciseId,
	}

	// 执行删除操作
	result := db.Where(conditions).Delete(&models.Exercise{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到符合条件的视频记录")
	}
	return nil
}

func GetVideoByPointId(pointId int64) ([]models.Video, error) {
	db := db.GetDB()
	var videos []models.Video
	err := db.Where("knowledge_point_id = ?", pointId).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
func GetFileByPointId(pointId int64, resourceType string) (interface{}, error) {
	switch resourceType {
	case "exercise":
		exercises, err := getExerciseByPointId(pointId)
		if err != nil {
			return nil, err
		}
		return exercises, nil
	case "courseware":
		coursewares, err := getCoursewareByPointId(pointId)
		if err != nil {
			return nil, err
		}
		return coursewares, nil
	default:
		log.Errorf("unsupported resource type: %s", resourceType)
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}
func getExerciseByPointId(pointId int64) ([]models.Exercise, error) {
	db := db.GetDB()
	log.Infof("pointId: %d", pointId)
	var exercises []models.Exercise
	err := db.Where("knowledge_point_id = ?", pointId).Find(&exercises).Error
	if err != nil {
		return nil, err
	}
	// log.Infof("exercises: ", exercises)
	return exercises, nil
}
func getCoursewareByPointId(pointId int64) ([]models.Courseware, error) {
	db := db.GetDB()
	var coursewares []models.Courseware
	err := db.Where("knowledge_point_id = ?", pointId).Find(&coursewares).Error
	if err != nil {
		return nil, err
	}
	return coursewares, nil
}
