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

// checkCmd implements the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the integrity of the galaxy file",
	Long:  `Runs optional checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "D:/GoLand/farHorizons/testdata/galaxy.json"
		// load all the data
		galaxy, err := fh.GetGalaxy(name)
		if err != nil {
			return err
		}

		allChecks, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		duplicateStars, err := cmd.Flags().GetBool("duplicate-stars")
		if err != nil {
			return err
		}

		var errors []error

		if allChecks || duplicateStars {
			stars := make(map[string]bool)
			for i, star := range galaxy.Stars {
				coords := fmt.Sprintf("%04d/%04d/%04d", star.X, star.Y, star.Z)
				if ok := stars[coords]; ok {
					errors = append(errors, fmt.Errorf("duplicate star %d: %s", i, coords))
				}
				stars[coords] = true
			}
		}

		if errors != nil {
			for _, err := range errors {
				fmt.Printf("%+v\n", err)
			}
			return fmt.Errorf("errors found")
		}

		fmt.Printf("no errors found")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().Bool("all", false, "perform all checks")
	checkCmd.Flags().BoolP("duplicate-stars", "s", false, "check for duplicate stars")
}
