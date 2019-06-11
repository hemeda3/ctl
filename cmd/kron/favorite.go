package kron

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	KronCmd.AddCommand(favoriteCmd)
	viper.SetDefault("favorites", make(map[string]location))
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.kron")
	err := viper.ReadInConfig()
	if err != nil {
		// Write config file
		fmt.Println("Creating new config file")
		createConfig()
		// panic(err.Error())
	}
	favoriteCmd.Flags().StringSliceP("context", "c", []string{}, "Specific contexts to list cronjobs from")
	favoriteCmd.Flags().StringP("namespace", "n", "", "Specific namespace to list cronjobs from within contexts")
}

var favoriteCmd = &cobra.Command{
	Use:   "favorite [jobs]",
	Short: "Adds jobs to favorite list",
	Long:  "Adds specified job(s) to the favorite list. If no job was specified the selected job is added.",
	Run: func(cmd *cobra.Command, args []string) {
		// args/flags
		ctxs, _ := cmd.Flags().GetStringSlice("context")
		nss, _ := cmd.Flags().GetString("namespace")

		f, err := getFavorites()
		if err != nil {
			panic(err.Error())
		}

		if len(args) == 0 {
			selected, err := getSelected()
			if err != nil {
				panic(err.Error())
			}
			if l, ok := f[selected.Name]; ok {
				fmt.Println(overrideFavoriteMessage(selected.Name, l))
			}
			f[selected.Name] = selected.Location
		} else {
			for _, job := range args {
				if l, ok := f[job]; ok {
					fmt.Println(overrideFavoriteMessage(job, l))
				}
				f[job] = location{ctxs, nss}
			}
		}

		viper.Set("favorites", f)
		viper.WriteConfig()
	},
}
