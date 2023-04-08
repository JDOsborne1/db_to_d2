package main

import "database/sql"

type UserPermission struct {
	User   string
	Table  string
	Column string
	Select bool
}

// The `get_user_permissions` function applies only at the column permission level.
// It does not apply at the table level, and it does not apply at the database level.
// A different function will be needed to apply at those levels.
func get_user_permissions(db *sql.DB) ([]UserPermission, error) {
	var permissions []UserPermission

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

		permissions = append(permissions, UserPermission{
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
