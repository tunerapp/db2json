/*
db2json for Tunerapp
Copyright (C) 2021  Louis Brauer <louis@brauer.family>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	DatabaseURL string `cli:"*d,database-url" usage:"connection string to mysql db, format username:password@protocol(address)/dbname?param=value"`
	OutputDir   string `cli:"*o,output-dir" usage:"existing directory where the JSON files will be emitted"`
}

type Station struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Homepage    string `json:"homepage"`
	Favicon     string `json:"favicon"`
	Creation    string `json:"creation"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Language    string `json:"language"`
	Tags        string `json:"tags"`
	Subcountry  string `json:"state"`
	Bitrate     int    `json:"bitrate"`
}

func main() {
	printLicense()

	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		return runExport(argv.DatabaseURL, argv.OutputDir)
	}))
}

func runExport(dbURL, outputDir string) error {
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		log.Fatalf("Output directory '%s' does not exist, aborting.", outputDir)
	}
	fmt.Println("Output Dir   :", outputDir)

	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("Database URL : ", dbURL)

	stmt := `SELECT StationUuid, Name, Url, Homepage, 
	Favicon, Creation, Country, 
	CountryCode, Language, Tags,
	Subcountry, Bitrate FROM Station`
	rows, err := db.Query(stmt)
	if err != nil {
		log.Fatalf("Cannot query database: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var entry Station
		// TODO Export all columns
		if err := rows.Scan(&entry.UUID,
			&entry.Name,
			&entry.URL,
			&entry.Homepage,
			&entry.Favicon,
			&entry.Creation,
			&entry.Country,
			&entry.CountryCode,
			&entry.Language,
			&entry.Tags,
			&entry.Subcountry,
			&entry.Bitrate); err != nil {
			log.Fatalf("Cannot scan row: %v", err)
		}
		log.Printf("id %s\n", entry.UUID)
		raw, err := json.MarshalIndent(entry, "", "    ")
		if err != nil {
			log.Fatalf("Cannot marshal row: %v", err)
		}
		fpath := filepath.Join(outputDir, fmt.Sprintf("%s.json", entry.UUID))
		if err := ioutil.WriteFile(fpath, raw, 0644); err != nil {
			log.Fatalf("Cannot write file %s: %v", fpath, err)
		}
	}

	return nil
}

func printLicense() {
	fmt.Println(`
db2json        
	
    Copyright (C) 2021  Tunerapp developers
    This program comes with ABSOLUTELY NO WARRANTY; for details see LICENSE'.
    This is free software, and you are welcome to redistribute it
    under certain conditions.
	`)
}
