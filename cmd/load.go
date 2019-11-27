// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"log"
	"os"
	"upper.io/db.v3"

	"upper.io/db.v3/ql"

	"deepsea/global"
	"github.com/spf13/cobra"
	jlog "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var SourceFile string
var IdentifierRegex string
var DropTable bool

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load Marks from a file",
	Long:  `LOAD: Help here`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Println("loadDriver()")
		loadDriver(cmd, args)
	},
}

func init() {

	loadCmd.Flags().StringVarP(&DBFile,
		"DBFile",
		"d",
		"",
		"Path to QL DB file")

	loadCmd.Flags().StringVarP(
		&SourceFile,
		"SourceFile",
		"s",
		"",
		"Path to Source of marks file")

	loadCmd.Flags().StringVarP(
		&IdentifierRegex,
		"IdentifierRegex",
		"r",
		"",
		"<dynamic> RegEx pattern")

	if err = viper.BindPFlag(
		"storage.DBFile",
		storageCmd.Flags().Lookup("DBFile")); err != nil {
		_ = storageCmd.Help()
		jlog.DEBUG.Println("Setting DBFile")
		os.Exit(2)
	}

	loadCmd.Flags().BoolVarP(&DropTable,
		"DropTable",
		"D",
		false,
		"Drop Table, do not truncate data")

	if err = viper.BindPFlag(
		"storage.load.SourceFile",
		loadCmd.Flags().Lookup("SourceFile")); err != nil {
		jlog.DEBUG.Println("Setting SourceFile")
		_ = loadCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"storage.load.IdentifierRegex",
		loadCmd.Flags().Lookup("IdentifierRegex")); err != nil {
		jlog.DEBUG.Println("Setting IdentifierRegex")
		_ = loadCmd.Help()
		os.Exit(2)
	}

	storageCmd.AddCommand(loadCmd)
}

func loadDriver(cmd *cobra.Command, args []string) {

	var markCollection db.Collection

	var settings = ql.ConnectionURL{
		Database: viper.GetString("storage.DBFile"),
	}

	sess, err := ql.Open(settings)
	if err != nil {
		jlog.ERROR.Printf("db.Open(): %q\n", err)
		os.Exit(2)
	}
	defer sess.Close()

	// Option A: Remove Mark table
	if DropTable {
		jlog.DEBUG.Printf("Dropping table Mark if exists\n")
		_, err = sess.Exec(`DROP TABLE IF EXISTS mark`)

		jlog.DEBUG.Printf("Creating Marks table\n")
		_, err = sess.Exec(`CREATE TABLE mark ( 
			identifier string, 
			email string,
			firstname string,
			lastname string )`)
	}

	// Option B: Truncate Mark table data
	// Pointing to the "mark" table.
	jlog.DEBUG.Printf("Selecting the mark table \n")
	markCollection = sess.Collection("mark")

	// Attempt to remove existing rows (if any).
	jlog.DEBUG.Printf("Removing existing rows if any \n")
	err = markCollection.Truncate()
	if err != nil {
		jlog.TRACE.Printf("Truncate(): %q\n", err)
	}

	// Marks in CSV file
	if global.CSVFileRe.MatchString(
		viper.GetString("storage.load.SourceFile")) {

		jlog.DEBUG.Println("Matched Source File as a CSV file")
		var marks []global.Mark

		jlog.DEBUG.Println("Converting CSV2JSON for DB Load")
		marksJson, err := global.CSV2Json(
			viper.GetString("storage.load.SourceFile"))
		if err != nil {
			jlog.ERROR.Printf("Could not parse Source File CSV %v", err)
			os.Exit(3)
		}
		jlog.TRACE.Println(string(marksJson))

		jlog.DEBUG.Println("Unmarshalling JSON into Marks")
		err = json.Unmarshal(marksJson, &marks)
		if err != nil {
			jlog.ERROR.Println("JSON Unmarshal error:", err)
			os.Exit(2)
		}

		jlog.INFO.Println("Loading Marks into DB")
		ix := 1
		for k := range marks {

			jlog.INFO.Printf("[%d] Loading Mark\n", ix)
			if len(marks[k].Email) == 0 {
				jlog.WARN.Println("	Mark has no email. Skip...")
				continue
			}

			if !global.EmailRe.MatchString(marks[k].Email) {
				jlog.WARN.Println("	Mark email format is invalid. Skip...")
				continue
			}

			if marks[k].Identifier == "<dynamic>" {
				jlog.TRACE.Println("Mark has dynamic ID")
				// Regex generate
				marks[k].Identifier, err = global.RegToString(
					viper.GetString("storage.load.IdentifierRegex"))
				if err != nil {
					// Fallback to random strings
					jlog.WARN.Println(
						"IdentifierRegex problem? Setting a random string id")
					marks[k].Identifier = global.RandString(8)
				}
			}

			// Inserting rows into the "Mark" table.
			jlog.DEBUG.Printf("Checks Passed. Inserting record.\n")
			_, err = markCollection.Insert(global.Mark{
				Identifier: marks[k].Identifier,
				Email:      marks[k].Email,
				Firstname:  marks[k].Firstname,
				Lastname:   marks[k].Lastname,
			})
		}
	}

	// Query for the results we've just inserted.
	jlog.TRACE.Println("Finding all marks")
	res := markCollection.Find()

	// Query all results and fill the mark variable with them.
	var marks []global.Mark

	jlog.TRACE.Println("Getting all marks from collection")
	err = res.All(&marks)
	if err != nil {
		jlog.ERROR.Printf("res.All(): %q\n", err)
		os.Exit(3)
	}

	// Printing to stdout.
	jlog.INFO.Println("-= = = = Mark Database = = = =-")
	for _, mark := range marks {
		log.Printf("%s, %s, %s, %s\n",
			mark.Identifier,
			mark.Email,
			mark.Firstname,
			mark.Lastname,
		)
	}
}
