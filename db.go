package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"regexp"
	"strings"
)

func startQueryDb() (*[]string, *[]string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(
				color.Error,
				"[ %v ] %s\n",
				color.New(color.FgRed, color.Bold).Sprint("ERROR"),
				err,
			)
			os.Exit(int(exitCodeErrDB))
		}
	}()
	db, err := sql.Open("sqlite", path.Join(homeDir, "/.ping1s", dbName))
	if err != nil {
		panic(err)
	}
	return query(db)

}

func query(db *sql.DB) (*[]string, *[]string) {
	poetry := queryPoetry(db)
	hitokoto := queryHitokoto(db)

	return poetry, hitokoto
}

func queryPoetry(db *sql.DB) *[]string {
	commands := ""
	if commandArgs.Author != "-1" {
		commands += fmt.Sprintf("and author = '%s' ", commandArgs.Author)
	}
	if commandArgs.CollectionType != -1 {
		commands += fmt.Sprintf("and collection = %s ", commandArgs.Collection)
	}

	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM poetry where 1=1  %s ORDER BY RANDOM() limit %d`, commands, commandArgs.Num))

	defer rows.Close()

	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		rows, err = db.Query(fmt.Sprintf(`SELECT * FROM poetry  ORDER BY RANDOM() limit %d`, commandArgs.Num))
	}

	var result []string
	for rows.Next() {
		poetry := Poetry{}

		err := rows.Scan(&poetry.ID, &poetry.Author, &poetry.Dynasty, &poetry.Title, &poetry.Paragraphs, &poetry.Collection)
		if err != nil {
			panic(err)
		}

		var paragraphListTmp []string
		err = json.Unmarshal([]byte(poetry.Paragraphs.String), &paragraphListTmp)
		if err != nil {
			panic(err)
		}

		var paragraphList []string

		for _, item := range paragraphListTmp {
			item = strings.TrimSpace(item)
			re := regexp.MustCompile(`([。？！])`)
			newStr := re.ReplaceAllString(item, "$1-")
			newStr = strings.TrimRight(newStr, "-")
			split := strings.Split(newStr, "-")
			paragraphList = append(paragraphList, split...)
		}
		poetry.ParagraphList = paragraphList

		t := strings.Join([]string{poetry.Dynasty.String, poetry.Author.String}, " ")
		strList := []string{poetry.Title.String, t}
		strList = append(strList, paragraphList...)

		result = append(result, strList...)

		result = append(result, "", "")
	}

	max := 0
	for _, item := range result {
		tLength := runewidth.StringWidth(item)
		if tLength > max {
			max = tLength
		}
	}
	if max%2 != 0 {
		max += 1
	}

	for idx, item := range result {
		result[idx] = fillSpace(item, max)
	}

	return &result
}

func queryHitokoto(db *sql.DB) *[]string {
	commands := ""
	if commandArgs.Type != "-1" {
		commands += fmt.Sprintf("and type = '%s' ", commandArgs.Type)
	}

	hitokoto := Hitokoto{}

	err := db.QueryRow(fmt.Sprintf(`SELECT * FROM hitokoto where 1=1  %s ORDER BY RANDOM() limit 1`, commands)).
		Scan(&hitokoto.ID, &hitokoto.Hitokoto, &hitokoto.Type, &hitokoto.From, &hitokoto.FromWho)

	if err != nil {
		err := db.QueryRow(fmt.Sprintf(`SELECT * FROM hitokoto ORDER BY RANDOM() limit 1`)).
			Scan(&hitokoto.ID, &hitokoto.Hitokoto, &hitokoto.Type, &hitokoto.From, &hitokoto.FromWho)
		if err != nil {
			panic(err)
		}
	}

	result := []string{hitokoto.Hitokoto.String, hitokoto.From.String}
	return &result
}

func fillSpace(s string, length int) string {
	const space = " "
	for runewidth.StringWidth(s) < length {
		s = space + s + space
	}

	return s
}
