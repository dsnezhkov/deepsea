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
	"deepsea/global"
	"fmt"
	"github.com/spf13/cobra"
	jlog "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"os"
	"strings"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/ql"
)

var DBTask string

// managerCmd represents the db management command
var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Manage information in marks database",
	Long:  `MANAGER: Manage information in marks database`,
	Run: func(cmd *cobra.Command, args []string) {
		jlog.DEBUG.Println("managerDriver()")
		managerDriver(cmd, args)
	},
}

var optDBTaskMap = map[string]interface{}{
	"showmarks":    ShowMarks,
	"truncate":     TruncateMarks,
	"droptable":    DropMarks,
	"createtable":  CreateMarks,
	"recycletable": RecycleMarks,
}

func init() {

	optDBTaskMapKeys := make([]string, 0)
	for key := range optDBTaskMap {
		optDBTaskMapKeys = append(optDBTaskMapKeys, key)
	}

	managerCmd.Flags().StringVarP(
		&DBTask,
		"DBTask",
		"T",
		"",
		"Tasks to run: \n"+strings.Join(optDBTaskMapKeys, "\n"))

	if err = viper.BindPFlag(
		"storage.manager.DBTask",
		managerCmd.Flags().Lookup("DBTask")); err != nil {
		jlog.ERROR.Println("Setting DBTask")
		_ = managerCmd.Help()
		os.Exit(2)
	}

	storageCmd.AddCommand(managerCmd)
}

func managerDriver(cmd *cobra.Command, args []string) {

	optDBTaskMapKeys := make([]string, 0)
	for key := range optDBTaskMap {
		optDBTaskMapKeys = append(optDBTaskMapKeys, key)
	}

	jlog.DEBUG.Println("Setting DBFile link")
	var dbFile = viper.GetString("storage.DBFile")
	if ! global.FileExists(dbFile)  {
		jlog.ERROR.Fatalf("Database file does not exist: %s", dbFile )
	}
	var settings = ql.ConnectionURL{
		Database: dbFile,
	}

	sess, err := ql.Open(settings)
	if err != nil {
		jlog.ERROR.Printf("db.Open(): %q\n", err)
	}
	defer sess.Close()

	jlog.TRACE.Printf("Making a Collection")
	markCollection := sess.Collection("mark")

	jlog.TRACE.Printf("Getting a DBTask")
	dt := viper.GetString("storage.manager.DBTask")
	if val, ok := optDBTaskMap[dt]; ok {
		jlog.DEBUG.Printf("Found Valid Task: %s", dt)

		//
		jlog.TRACE.Println("Converting Task to a function call")
		val.(func(database sqlbuilder.Database, collection db.Collection))(
			sess, markCollection.(db.Collection))

	} else {
		jlog.ERROR.Fatalf(
			"Task:%s undefined (-T). Valid options: %s\n",
			dt, strings.Join(optDBTaskMapKeys, "|"))
	}
}

func ShowMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	// Query for the results we've just inserted.
	jlog.DEBUG.Printf("Querying for result : find()\n")
	res := markCollection.Find()

	// Query all results and fill the mark variable with them.
	var marks []global.Mark

	err = res.All(&marks)
	if err != nil {
		jlog.ERROR.Fatalf("res.All(): %q\n", err)
	}

	for _, mark := range marks {
		fmt.Printf("%s, %s, %s, %s.\n",
			mark.Identifier,
			mark.Email,
			mark.Firstname,
			mark.Lastname,
		)
	}
}

func TruncateMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	jlog.INFO.Println("Removing existing rows (if any)")
	err = markCollection.Truncate()
	if err != nil {
		jlog.ERROR.Fatalf("Truncating data in table failed: %q\n", err)
	}
}

func DropMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	jlog.INFO.Println("Dropping table Mark if exists")
	_, err = sess.Exec(`DROP TABLE IF EXISTS mark`)
	if err != nil {
		jlog.ERROR.Fatalf("Dropping table failed: %q\n", err)
	}
}

func CreateMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	jlog.INFO.Println("Creating Marks table")
	_, err = sess.Exec(`CREATE TABLE mark (
			identifier string,
			email string,
			firstname string,
			lastname string )`)
	if err != nil {
		jlog.ERROR.Fatalf("Creating table failed: %q\n", err)
	}
}

func RecycleMarks(sess sqlbuilder.Database, markCollection db.Collection) {
	jlog.DEBUG.Println("Recycling Marks table")
	DropMarks(sess, markCollection)
	CreateMarks(sess, markCollection)
}
