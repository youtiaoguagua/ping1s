package main

import (
	"database/sql"
	"os"
)

type Poetry struct {
	ID            sql.NullInt32  `json:"id"`
	Author        sql.NullString `json:"author"`
	Dynasty       sql.NullString `json:"dynasty"`
	Title         sql.NullString `json:"title"`
	Paragraphs    sql.NullString `json:"paragraphs"`
	ParagraphList []string
	Collection    sql.NullString `json:"collection"`
}

type Hitokoto struct {
	ID       sql.NullInt32  `json:"id"`
	Type     sql.NullString `json:"type"`
	Hitokoto sql.NullString `json:"hitokoto"`
	From     sql.NullString `json:"from"`
	FromWho  sql.NullString `json:"from_who"`
}

type CommandArgs struct {
	Author         string
	CollectionType int
	Collection     string
	Type           string
	Num            int
	Version        bool
}

type Logger struct {
	file *os.File
	size int64
}
