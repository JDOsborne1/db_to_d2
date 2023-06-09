package mysql

import (
	"core"
	"database/sql"
)

// userColumnPermission is a struct that contains the information about a user's permission
// on a column.
type userColumnPermission struct {
	User   string
	Table  string
	Column string
	Select bool
}

// The `get_column_level_permissions` function applies only at the column permission level.
// It does not apply at the table level, and it does not apply at the database level.
// A different function will be needed to apply at those levels.
// The function does not discriminate between users. It returns a slice of column permissions
// for all users.
func get_column_level_permissions(db *sql.DB) ([]userColumnPermission, error) {
	var permissions []userColumnPermission

	rows, err := db.Query(`
		SELECT
			Grantee,
			C.Table_Name,
			C.Column_Name,
			(SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMN_PRIVILEGES
				WHERE TABLE_NAME = C.Table_Name AND COLUMN_NAME = C.Column_Name AND PRIVILEGE_TYPE = 'SELECT'
				AND GRANTEE = CP.GRANTEE) AS Has_Select_Permission
		FROM
			INFORMATION_SCHEMA.COLUMNS C
			INNER JOIN INFORMATION_SCHEMA.TABLES T ON C.Table_Name = T.Table_Name
			LEFT JOIN INFORMATION_SCHEMA.COLUMN_PRIVILEGES CP ON C.Table_Name = CP.Table_Name AND C.Column_Name = CP.Column_Name
		WHERE
			T.TABLE_TYPE = 'BASE TABLE'
		AND Grantee IS NOT NULL
		ORDER BY
			Grantee, Table_Name, Column_Name;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		var table string
		var column string
		var hasSelectPermission sql.NullInt64

		err = rows.Scan(&user, &table, &column, &hasSelectPermission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, userColumnPermission{
			User:   user,
			Table:  table,
			Column: column,
			Select: hasSelectPermission.Valid && hasSelectPermission.Int64 > 0,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// userTablePermission is a struct that contains the information about a user's permission
// on a table.
type userTablePermission struct {
	User   string
	Table  string
	Select bool
}

// The `get_table_level_permissions` function applies only at the column permission level.
// It does not apply at the table level, and it does not apply at the database level.
// A different function will be needed to apply at those levels.
// This will not discriminate based on user, so will return a slice of all table permissions
// for all users.
func get_table_level_permissions(db *sql.DB) ([]userTablePermission, error) {
	var permissions []userTablePermission

	rows, err := db.Query(`
	SELECT Grantee,
		T.Table_Name,
		COUNT(*) AS Has_Select_Permission
	FROM INFORMATION_SCHEMA.TABLES T
		LEFT JOIN INFORMATION_SCHEMA.TABLE_PRIVILEGES TP ON T.Table_Name = TP.Table_Name
	WHERE T.TABLE_TYPE = 'BASE TABLE'
		AND TP.PRIVILEGE_TYPE = 'SELECT'
		AND Grantee IS NOT NULL
	GROUP BY Grantee,
		Table_Name
	ORDER BY Grantee,
		Table_Name;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user string
		var table string
		var hasSelectPermission sql.NullInt64

		err = rows.Scan(&user, &table, &hasSelectPermission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, userTablePermission{
			User:   user,
			Table:  table,
			Select: hasSelectPermission.Valid && hasSelectPermission.Int64 > 0,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func permission_driven_restrictor(_table_permissions []userTablePermission, _column_permissions []userColumnPermission, _for_user string) core.Restrictor {
	table_permission_map := make(map[string]bool)
	for _, permission := range _table_permissions {
		if permission.User == _for_user {
			table_permission_map[permission.Table] = permission.Select
		}
	}

	column_permission_map := make(map[string]map[string]bool)
	for _, permission := range _column_permissions {
		if permission.User == _for_user {
			if _, ok := column_permission_map[permission.Table]; !ok {
				column_permission_map[permission.Table] = make(map[string]bool)
			}
			column_permission_map[permission.Table][permission.Column] = permission.Select
		}
	}

	return func(_table core.Table, _column core.Column) bool {
		allowed := table_permission_map[_table.Name] || column_permission_map[_table.Name][_column.Name]
		return !allowed
	}
}

func Restrict_to_table_for_user(_username string) (core.Restrictor, error) {
	db, err := connect_to_db()
	if err != nil {
		return core.Standard, err
	}

	table_permissions, err := get_table_level_permissions(db)
	if err != nil {
		return core.Standard, err
	}

	column_permissions, err := get_column_level_permissions(db)
	if err != nil {
		return core.Standard, err
	}

	return permission_driven_restrictor(table_permissions, column_permissions, _username), nil
}
