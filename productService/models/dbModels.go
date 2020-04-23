package models

import (
	"database/sql"
)

type DBProducts struct {
	DBID            sql.NullString
	DBName          sql.NullString
	DBDescription   sql.NullString
	DBPrice         sql.NullFloat64
	DBDeliveryPrice sql.NullFloat64
}

type DBProductOptions struct {
	DBID          sql.NullString
	DBProductID   sql.NullString
	DBName        sql.NullString
	DBDescription sql.NullString
}
