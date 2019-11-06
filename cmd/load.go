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
	"fmt"
	"log"
	"os"

	"upper.io/db.v3/ql"

	"github.com/dsnezhkov/deepsea/global"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SourceFile string
var IdentifierRegex string

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load Marks from a file",
	Long:  `LOAD: Help here`,
	Run: func(cmd *cobra.Command, args []string) {
		loadDriver(cmd, args)
	},
}

func init() {

	loadCmd.Flags().StringVarP(&DBFile, "DBFile", "d",
		"", "Path to QL DB file")
	loadCmd.Flags().StringVarP(&SourceFile, "SourceFile", "s",
		"", "Path to Source of marks file")
	loadCmd.Flags().StringVarP(&IdentifierRegex, "IdentifierRegex", "r",
		"", "<dynamic> RegEx pattern")

	if err = viper.BindPFlag(
		"storage.DBFile", storageCmd.Flags().Lookup("DBFile")); err != nil {
		_ = storageCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"storage.load.SourceFile", loadCmd.Flags().Lookup("SourceFile")); err != nil {
		_ = loadCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"storage.load.IdentifierRegex", loadCmd.Flags().Lookup("IdentifierRegex")); err != nil {
		_ = loadCmd.Help()
		os.Exit(2)
	}
	storageCmd.AddCommand(loadCmd)
}

func loadDriver(cmd *cobra.Command, args []string) {

	var settings = ql.ConnectionURL{
		Database: viper.GetString("storage.DBFile"), // Path to database file.
	}
	// Attemping to open the "example.db" database file.
	sess, err := ql.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close() // Remember to close the database session.

	// Option A: Remove Mark table
	log.Printf("Dropping table Mark if exists\n")
	_, err = sess.Exec(`DROP TABLE IF EXISTS mark`)
	log.Printf("Creating Marks table\n")
	_, err = sess.Exec(`CREATE TABLE mark ( 
			identifier string, 
			email string,
			firstname string,
			lastname string )`)

	// Option B: Truncate Mark table data
	// Pointing to the "mark" table.
	log.Printf("Pointing to mark table \n")
	markCollection := sess.Collection("mark")

	// Attempt to remove existing rows (if any).
	log.Printf("Removing existing rows if any \n")
	err = markCollection.Truncate()
	if err != nil {
		log.Printf("Truncate(): %q\n", err)
	}

	// Marks in CSV file
	if global.CSVFileRe.MatchString(viper.GetString("storage.load.SourceFile")) {
		var marks []global.Mark

		marksJson, err := global.CSV2Json(viper.GetString("storage.load.SourceFile"))
		if err != nil {
			fmt.Printf("ERROR: Could not parse CSV %v", err)
			os.Exit(3)
		}
		//fmt.Println(string(marksJson))

		err = json.Unmarshal(marksJson, &marks)
		if err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			os.Exit(2)
		}

		for k := range marks {
			// TODO: Preprocess rules
			if marks[k].Identifier == "<dynamic>" {
				// Regex generate
				marks[k].Identifier, err = global.RegToString(
					viper.GetString("storage.load.IdentifierRegex"))
				if err != nil {
					// Fallback to random strings
					marks[k].Identifier = global.RandString(8)
				}
			}

			// Inserting rows into the "Mark" table.
			log.Printf("Inserting a row\n")
			_, err = markCollection.Insert(global.Mark{
				Identifier: marks[k].Identifier,
				Email:      marks[k].Email,
				Firstname:  marks[k].Firstname,
				Lastname:   marks[k].Lastname,
			})
		}
	}

	// Let's query for the results we've just inserted.
	log.Printf("Querying for result : find()\n")
	res := markCollection.Find()

	// Query all results and fill the mark variable with them.
	var marks []global.Mark

	log.Printf("Getting all results\n")
	err = res.All(&marks)
	if err != nil {
		log.Fatalf("res.All(): %q\n", err)
	}

	// Printing to stdout.
	log.Printf("Printing Marks\n")
	for _, mark := range marks {
		fmt.Printf("%s, %s, %s, %s.\n",
			mark.Identifier,
			mark.Email,
			mark.Firstname,
			mark.Lastname,
		)
	}
}
