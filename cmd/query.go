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
	"fmt"
	"github.com/dsnezhkov/deepsea/global"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/ql"
)

var DBTask string

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query storage",
	Long:  `QUERY: Help`,
	Run: func(cmd *cobra.Command, args []string) {
		queryDriver(cmd, args)
	},
}

//var optDBTask = []string{"showmarks", "somethingelse"}

var optDBTaskMap = map[string]interface{}{
	"showmarks":    qShowMarks,
	"truncate":     qTruncateMarks,
	"droptable":    qDropMarks,
	"createtable":  qCreateMarks,
	"recycletable": qRecycleMarks,
}

func init() {

	optDBTaskMapKeys := make([]string, 0)
	for key := range optDBTaskMap {
		optDBTaskMapKeys = append(optDBTaskMapKeys, key)
	}

	queryCmd.Flags().StringVarP(&DBFile, "DBFile", "d",
		"", "Path to QL DB file")
	queryCmd.Flags().StringVarP(&DBTask, "DBTask", "t",
		"", "Tasks to run: \n"+strings.Join(optDBTaskMapKeys, "\n"))

	if err = viper.BindPFlag(
		"storage.DBFile", queryCmd.Flags().Lookup("DBFile")); err != nil {
		_ = queryCmd.Help()
		os.Exit(2)
	}
	if err = viper.BindPFlag(
		"storage.query.DBTask", queryCmd.Flags().Lookup("DBTask")); err != nil {
		_ = queryCmd.Help()
		os.Exit(2)
	}
	storageCmd.AddCommand(queryCmd)
}

func queryDriver(cmd *cobra.Command, args []string) {

	optDBTaskMapKeys := make([]string, 0)
	for key := range optDBTaskMap {
		optDBTaskMapKeys = append(optDBTaskMapKeys, key)
	}

	var settings = ql.ConnectionURL{
		Database: viper.GetString("storage.DBFile"), // Path to database file.
	}

	sess, err := ql.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close() // Remember to close the database session.

	markCollection := sess.Collection("mark")

	dt := viper.GetString("storage.query.DBTask")
	if val, ok := optDBTaskMap[dt]; ok {
		log.Printf("Task: %s", dt)
		// Convert string to a function call:
		val.(func(database sqlbuilder.Database, collection db.Collection))(
			sess, markCollection.(db.Collection))

	} else {
		log.Printf("Task:%s undefined. Valid options:\n", dt)
		log.Printf("%s\n", strings.Join(optDBTaskMapKeys, "|"))
	}
}

func qShowMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	// Let's query for the results we've just inserted.
	log.Printf("Querying for result : find()\n")
	res := markCollection.Find()

	// Query all results and fill the mark variable with them.
	var marks []global.Mark

	err = res.All(&marks)
	if err != nil {
		log.Fatalf("res.All(): %q\n", err)
	}

	// Printing to stdout.
	fmt.Printf("-= Table: Marks =-\n")
	for _, mark := range marks {
		fmt.Printf("%s, %s, %s, %s.\n",
			mark.Identifier,
			mark.Email,
			mark.Firstname,
			mark.Lastname,
		)
	}
}

func qTruncateMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	log.Printf("Removing existing rows (if any) \n")
	err = markCollection.Truncate()
	if err != nil {
		log.Printf("Truncate(): %q\n", err)
	}
}
func qDropMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	log.Printf("Dropping table Mark if exists\n")
	_, err = sess.Exec(`DROP TABLE IF EXISTS mark`)
}
func qCreateMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	log.Printf("Creating Marks table\n")
	_, err = sess.Exec(`CREATE TABLE mark (
			identifier string,
			email string,
			firstname string,
			lastname string )`)
}
func qRecycleMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	qDropMarks(sess, markCollection)
	qCreateMarks(sess, markCollection)
}
