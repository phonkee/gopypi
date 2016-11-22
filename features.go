/*
Features

gopypi provides optional features that can be enabled/disabled directly from administration

List of features:

* auto_maintainers - automatically assign maintainers from package emails (users must exist in database)
* download_stats - enable download statistics for downloads

*/
package core

import "github.com/jinzhu/gorm"

const (
	FEATURE_AUTO_MAINTAINERS = "auto_maintainers"
	FEATURE_DOWNLOAD_STATS   = "download_stats"
)

const (
	DESCRIPTION_FEATURE_AUTO_MAINTAINERS = "Automatically assign maintainers on uploading package"
	DESCRIPTION_FEATURE_DOWNLOAD_STATS   = "Store download statistics of packages"
)

/*
createFeatures creates all needed features in database if they don't exist, it's run during all migrations
*/
func createFeatures(db *gorm.DB) (err error) {

	allFeatures := []Feature{
		{
			ID:          FEATURE_AUTO_MAINTAINERS,
			Description: DESCRIPTION_FEATURE_AUTO_MAINTAINERS,
		},
		{
			ID:          FEATURE_DOWNLOAD_STATS,
			Description: DESCRIPTION_FEATURE_DOWNLOAD_STATS,
		},
	}

	for _, feature := range allFeatures {
		target := Feature{}
		if err = db.First(&target, "id = ?", feature.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err = db.Create(&feature).Error; err != nil {
					return err
				}
				err = nil
			} else {
				return
			}
		} else {
			// Update description if changed
			if target.Description != feature.Description {
				target.Description = feature.Description
				if err = db.Save(&feature).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
