package cmd

import (
	"fmt"
	"os"

	"github.com/pengux/cruder/generator"
	"github.com/spf13/cobra"
)

// var cfgFile string
var (
	pkgName         string
	funcs           []string
	skipFuncSuffix  bool
	readFields      []string
	writeFields     []string
	primaryField    string
	softDeleteField string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cruder",
	Short: "Generate code for CRUD functions from a Go struct",
	Long: `cruder is a tool to generate code for Create, Read, Update, Delete functions
from a Go struct. It supports multiple generators which are listed in the 'Available
Commands section'. Functions that can be generated are:
- Create: Adds an entry
- Read: Gets an entry using an ID
- List: Gets multiple entries
- Update: Updates an entry
- Delete: Deletes an entry using an ID
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	//
	// // Here you will define your flags and configuration settings.
	// // Cobra supports persistent flags, which, if defined here,
	// // will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	//
	// // Cobra also supports local flags, which will only run
	// // when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&pkgName, "pkg", "", "package name for the generated code, default to the same package from input")
	RootCmd.PersistentFlags().StringSliceVar(&funcs, "fn", []string{
		string(generator.Create),
		string(generator.Get),
		string(generator.List),
		string(generator.Update),
		string(generator.Delete),
	}, `CRUD functions to generate, e.g. --fn "create" --fn "delete". Default to all functions`)
	RootCmd.PersistentFlags().BoolVar(&skipFuncSuffix, "skipsuffix", false, "Skip adding the struct name as suffix to the generated functions")
	RootCmd.PersistentFlags().StringVar(&primaryField, "primaryfield", "", "the field to use as primary key. Default to 'ID' if it exists in the <struct>")
	RootCmd.PersistentFlags().StringVar(&softDeleteField, "softdeletefield", "", "the field to use for softdelete (should be of type nullable datetime field). Default to 'DeletedAt' if it exists in the <struct>")
	RootCmd.PersistentFlags().StringSliceVar(&readFields, "readfields", []string{}, "Fields in the struct that should be used for read operations (get,list). Default to all fields except the one used for softdelete")
	RootCmd.PersistentFlags().StringSliceVar(&writeFields, "writefields", []string{}, "Fields in the struct that should be used for write operations (create,update). Default to all fields")

}
