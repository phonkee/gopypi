/*
filter provides set of FilterFunc functions.
*/
package core

import "github.com/jinzhu/gorm"

// FilterFunc is callback which is called in some methods (Get, List)
type FilterFunc func(*gorm.DB) *gorm.DB

/*
ApplyFilterFuncs applies all FilterFuncs and returns db
*/
func ApplyFilterFuncs(db *gorm.DB, funcs ...FilterFunc) *gorm.DB {
	for _, f := range funcs {
		db = f(db)
	}

	return db
}

/*
FilterUsername filters user by username
*/
func FFUsername(username string) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}

/*
FFID adds clause to search by id
*/
func FFID(id interface{}) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

/*
FFOrderBy orders by
*/
func FFOrderBy(order string) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(order, false)
	}
}


/*
FFPackagesFor filters packages by given user, user is either author or maintainer
*/
func FFPackagesFor(user User) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		maintainers := []User{}
		other := db.New()

		ids := make([]uint, 0, len(maintainers))

		var pid uint
		rows, err := other.Table("package_maintainers").Where("user_id = ?", user.ID).Select("package_id").Rows()
		if err != nil {
			db.Error = err
			return db
		}
		for rows.Next() {
			rows.Scan(&pid)
			ids = append(ids, pid)
		}

		// if we have found packages that are maintained by user
		if len(ids) > 0 {
			db = db.Where("author_id = ? OR id IN (?)", user.ID, ids)
		} else {
			db = db.Where("author_id = ?", user.ID)
		}

		return db
	}
}

/*
FFPreload add preloads to que
*/
func FFPreload(preloads ...string) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range preloads {
			db = db.Preload(preload)
		}
		return db
	}
}

/*
FFDownloadStatsPackage filters stats by package
*/
func FFDownloadStatsPackage(pack *Package) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		versions := []PackageVersion{}

		// first select all versions for IN clause
		other := db.New()
		other.Model(pack).Related(&versions)

		// prepare list of ids
		ids := make([]uint, 0, len(versions))
		for _, version := range versions {
			ids = append(ids, version.ID)
		}

		// if ids, add where clause, otherwise select none
		if len(ids) > 0 {
			db = db.Where("package_version_id IN (?)", ids)
		} else {
			db = db.Where("package_version_id IN (0)", ids)
		}

		return db
	}
}

/*
FFDownloadStatsPackage filters stats by package
*/
func FFDownloadStatsPackageVersion(pv *PackageVersion) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("package_version_id = ?", pv.ID)
		return db
	}
}


/*
FFWhere is shorthand for gorm Where
 */
func FFWhere(query interface{}, args ...interface{}) FilterFunc {
	return func(db *gorm.DB) *gorm.DB {
		queryset := db.Where(query, args...)
		return queryset
	}
}
