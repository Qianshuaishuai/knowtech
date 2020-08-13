package models

import "regexp"

var (
	Regexp_dbConfig_Sql                = regexp.MustCompile(`\?`)
	Regexp_dbConfig_NumericPlaceHolder = regexp.MustCompile(`\$\d+`)
)
