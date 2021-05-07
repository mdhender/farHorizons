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
	"github.com/mdhender/farHorizons/internal/fh"
	"github.com/spf13/cobra"
)

// createSpeciesCmd implements the create species command
var createSpeciesCmd = &cobra.Command{
	Use:   "species",
	Short: "Create a new species",
	Long: `This command creates a new species record using information
from ???something???`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		g, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}

		// seed random number generator
		fh.Seed(0xC0FFEE)

		return g.Write(name)
	},
}

func init() {
	createCmd.AddCommand(createSpeciesCmd)
	createSpeciesCmd.Flags().StringP("species-file", "f", "D:/GoLand/farHorizons/testdata/sp1.json", "file containing species data as a JSON object")
}
