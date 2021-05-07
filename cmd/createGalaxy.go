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
	"time"
)

// createGalaxyCmd implements the create galaxy command
var createGalaxyCmd = &cobra.Command{
	Use:   "galaxy",
	Short: "Create a new galaxy",
	Long: `This commands loads setup data from a
configuration file, then creates a new galaxy file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		started := time.Now()
		fh.Seed(0xC0FFEE) // seed random number generator

		galaxyFileName, err := cmd.Flags().GetString("galaxy-file")
		if err != nil {
			return err
		} else if galaxyFileName == "" {
			return fmt.Errorf("you must specify a valid file name to create")
		}
		setupFileName, err := cmd.Flags().GetString("setup-file")
		if err != nil {
			return err
		} else if setupFileName == "" {
			return fmt.Errorf("you must specify a valid setup file name")
		}

		setupData, err := fh.GetSetup(setupFileName)
		if err != nil {
			return err
		}

		g, err := fh.GenerateGalaxy(setupData)
		if err != nil {
			return err
		}
		err = g.Write(galaxyFileName)
		if err != nil {
			return err
		}

		fmt.Printf("Created file %q in %v\n", galaxyFileName, time.Now().Sub(started))
		return nil
	},
}

func init() {
	createCmd.AddCommand(createGalaxyCmd)
	createGalaxyCmd.Flags().StringP("galaxy-file", "g", "", "name of galaxy file to create")
	_ = createGalaxyCmd.MarkFlagRequired("galaxy-file")
	createGalaxyCmd.Flags().StringP("setup-file", "i", "", "name of configuration file to load")
	_ = createGalaxyCmd.MarkFlagRequired("setup-file")
}
