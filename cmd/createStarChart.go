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
	"github.com/spf13/cobra"
	"io/ioutil"
)

// createStarChartCmd implements the create homes command
var createStarChartCmd = &cobra.Command{
	Use:   "star-chart",
	Short: "Graph nearby stars",
	Long:  `Create a graph showing stars within a tolerable jump factor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		graviticsLevel, err := cmd.Flags().GetInt("gravitics-level")
		if err != nil {
			return err
		}
		mishapLimit, err := cmd.Flags().GetInt("mishap-limit")
		if err != nil {
			return err
		}
		shipAge, err := cmd.Flags().GetInt("ship-age")
		if err != nil {
			return err
		}
		x, err := cmd.Flags().GetInt("x-origin")
		if err != nil {
			return err
		}
		y, err := cmd.Flags().GetInt("y-origin")
		if err != nil {
			return err
		}
		z, err := cmd.Flags().GetInt("z-origin")
		if err != nil {
			return err
		}

		name := "D:/GoLand/farHorizons/testdata/starlist.json"
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		type Star struct {
			X, Y, Z int
			Name    string
		}
		type Link struct {
			From     string `json:"source"`
			To       string `json:"target"`
			Distance int    `json:"value"`
		}
		var stars []*Star
		err = json.Unmarshal(b, &stars)
		if err != nil {
			return err
		}

		// find the origin and ensure that all stars have unique names
		origin := stars[0]
		names := make(map[string]bool)
		for _, star := range stars {
			if star.X == x && star.Y == y && star.Z == z {
				origin = star
				if star.Name == "" {
					star.Name = "Origin"
				}
			} else if star.Name == "" {
				star.Name = fmt.Sprintf("%d,%d,%d", star.X, star.Y, star.Z)
			}
			if ok := names[star.Name]; ok {
				return fmt.Errorf("duplicate name %q", star.Name)
			}
			names[star.Name] = true
		}
		type node struct {
			Id      string `json:"id"`
			Group   int    `json:"group"`
			star    *Star
			plotted bool
		}
		var links []*Link
		nodes := make(map[string]*node)
		group := 1
		nodes[origin.Name] = &node{Id: origin.Name, Group: 1, star: origin}
		for {
			group++
			added := 0
			for _, this := range nodes {
				if this.plotted {
					continue
				}
				from := this.star
				for _, to := range stars {
					if from == to {
						continue
					}
					deltaX, deltaY, deltaZ := from.X-to.X, from.Y-to.Y, from.Z-to.Z
					mishap_age := shipAge
					mishap_gv := graviticsLevel
					mishap_chance := (100 * (((deltaX) * (deltaX)) + ((deltaY) * (deltaY)) + ((deltaZ) * (deltaZ)))) / mishap_gv
					if mishap_chance > 10000 {
						mishap_chance = 10000
					} else if mishap_age > 0 {
						/* Add aging effect. */
						success_chance := 10000 - mishap_chance
						success_chance -= (2 * mishap_age * success_chance) / 100
						if success_chance < 0 {
							success_chance = 0
						}
						mishap_chance = 10000 - success_chance
					}
					if mishap_chance > (mishapLimit * 100) {
						continue
					}
					if _, ok := nodes[to.Name]; ok {
						// don't add again
						continue
					}
					nodes[to.Name] = &node{Id: to.Name, Group: group, star: to}
					added++
					links = append(links, &Link{
						From:     from.Name,
						To:       to.Name,
						Distance: mishap_chance / 100,
					})
				}
			}
			if added == 0 {
				break
			}
		}

		ds3 := struct {
			Nodes []*node `json:"nodes"`
			Links []*Link `json:"links"`
		}{
			Links: links,
		}
		for _, n := range nodes {
			ds3.Nodes = append(ds3.Nodes, n)
		}

		name = "D:/GoLand/farHorizons/testdata/starChart.json"
		if b, err := json.MarshalIndent(&ds3, "  ", "  "); err != nil {
			return err
		} else if err := ioutil.WriteFile(name, b, 0644); err != nil {
			return err
		}

		fmt.Printf("Created %q.\n", name)
		return nil
	},
}

func init() {
	createCmd.AddCommand(createStarChartCmd)
	createStarChartCmd.Flags().IntP("gravitics-level", "g", 1, "gravitics level for jump calculations")
	createStarChartCmd.Flags().IntP("mishap-limit", "l", 40, "maximum mishap threshold to map")
	createStarChartCmd.Flags().IntP("ship-age", "a", 0, "age of ship jumping")
	createStarChartCmd.Flags().IntP("x-origin", "x", 1, "x coordinate to begin from")
	createStarChartCmd.Flags().IntP("y-origin", "y", 1, "y coordinate to begin from")
	createStarChartCmd.Flags().IntP("z-origin", "z", 1, "z coordinate to begin from")
}
