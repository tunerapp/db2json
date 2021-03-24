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
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	DatabaseURL string `cli:"*d,database-url" usage:"connection string to mysql db, format username:password@protocol(address)/dbname?param=value"`
	OutputDir   string `cli:"*o,output-dir" usage:"existing directory where the JSON files will be emitted"`
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
