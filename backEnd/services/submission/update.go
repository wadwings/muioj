package submission

import (
	"MuiOJ-backEnd/models"
	"MuiOJ-backEnd/services/db"
)

func Update(sid, timeUsed, spaceUsed uint32, status string, judge []models.JudgeResult) error {
	updateObj := models.Submission{
		TimeUsed:  timeUsed,
		SpaceUsed: spaceUsed,
		Status:    status,
		Judge:     judge,
	}

	if _, err := db.DB.Table("submission").
		Where("sid = ?", sid).
		Cols("time_used", "space_used", "status", "judge").
		Update(&updateObj); err != nil {
		return err
	}

	return nil
}
