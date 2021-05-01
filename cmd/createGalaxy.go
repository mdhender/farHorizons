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
	"io/ioutil"

	"github.com/spf13/cobra"
)

// createGalaxyCmd represents the galaxy command
var createGalaxyCmd = &cobra.Command{
	Use:   "galaxy",
	Short: "A galaxy manager?",
	Long: `The command line interface to manage a galaxy.

The number of species must be between ` + fmt.Sprintf("%d and %d", fh.MIN_SPECIES, fh.MAX_SPECIES) + `, inclusive.

The number of stars is based on the number of species, and will be somewhere between ` + fmt.Sprintf("%d and %d", fh.MIN_STARS, fh.MAX_STARS) + `.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("create galaxy called")
		if len(args) != 0 {
			return fmt.Errorf("create galaxy: unknown arguments")
		}
		numberOfSpecies, err := cmd.Flags().GetInt("number-of-species")
		if err != nil {
			return err
		}
		g, err := fh.GenerateGalaxy(numberOfSpecies)
		if err != nil {
			return err
		}
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		if b, err := json.MarshalIndent(g, "  ", "  "); err != nil {
			return err
		} else if err := ioutil.WriteFile(name, b, 0644); err != nil {
			return err
		}
		fmt.Printf("created %q\n", name)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createGalaxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createGalaxyCmd.PersistentFlags().Int("number-of-species", 9, "number of species to create in the galaxy")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createGalaxyCmd.Flags().IntP("number-of-species", "n", 9, "number of species to create in the galaxy")
	_ = createGalaxyCmd.MarkFlagRequired("number-of-species")
}
