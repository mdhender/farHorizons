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
	"fmt"
	"github.com/mdhender/farHorizons/internal/fh"
	"github.com/spf13/cobra"
)

// showTurnCmd implements the show turn command
var showTurnCmd = &cobra.Command{
	Use:   "turn",
	Short: "Show turn number",
	Long:  `The command line interface to show turn information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"

		// Get galaxy data.
		galaxy, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}

		// print the current turn number
		fmt.Printf("%d\n", galaxy.TurnNumber)

		return nil
	},
}

func init() {
	showCmd.AddCommand(showTurnCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showTurnCmd.PersistentFlags().Int("number-of-species", 9, "number of species to create in the galaxy")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showTurnCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
