/*
 * farHorizons - a clone of Far Horizons
 * Copyright (C) 2021  Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/farHorizons/internal/fh"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// createHomesCmd implements the create homes command
var createHomesCmd = &cobra.Command{
	Use:   "homes",
	Short: "Create home systems",
	Long:  `The command line interface to create home systems.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// seed random number generator
		fh.Seed(0xC0FFEE)

		//name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		if len(args) != 0 {
			return fmt.Errorf("create homes: unknown arguments")
		}

		for num_planets := 3; num_planets < 10; num_planets++ {
			filename := fmt.Sprintf("D:/GoLand/farHorizons/testdata/HS%d.json", num_planets)
			fmt.Printf("Now doing file '%s'...\n", filename)
			var planets []*fh.PlanetData
			for planets == nil {
				planets = fh.GenerateEarthLikePlanet(num_planets)
			}
			if data, err := json.MarshalIndent(&struct {
				Name    string
				Planets []*fh.PlanetData
			}{
				Name:    fmt.Sprintf("HS%d", num_planets),
				Planets: planets,
			}, "  ", "  "); err != nil {
				return err
			} else if err := ioutil.WriteFile(filename, data, 0644); err != nil {
				return err
			}
			fmt.Printf("Created %q.\n", filename)
		}

		return nil
	},
}

func init() {
	createCmd.AddCommand(createHomesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createHomesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createHomesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
