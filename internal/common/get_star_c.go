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

package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

/* In case gamemaster creates new star systems with Edit program. */
const NUM_EXTRA_STARS = 20

var num_stars int
var star_data_modified bool
var StarBase []*StarData

type StarData struct {
	X, Y, Z             int                    /* Coordinates. */
	Type                int                    /* Dwarf, degenerate, main sequence or giant. */ // was `type`
	Color               int                    /* Star color. Blue, blue-white, etc. */
	Size                int                    /* Star size, from 0 thru 9 inclusive. */
	NumPlanets          int                    /* Number of usable planets in star system. */
	HomeSystem          bool                   /* TRUE if this is a good potential home system. */
	WormHere            bool                   /* TRUE if wormhole entry/exit. */
	WormX, WormY, WormZ int                    /* Coordinates. */
	PlanetIndex         int                    /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	Message             int                    /* Message associated with this star system, if any. */
	VisitedBy           [NUM_CONTACT_WORDS]int /* A bit is set if corresponding species has been here. */
}

/* GetStarData loads star data from a JSON file. */
func GetStarData() {
	b, err := ioutil.ReadFile("D:/GoLand/farHorizon/testdata/stars.dat")
	if err != nil {
		panic(fmt.Sprintf("Cannot open file 'stars.dat': %+v", err))
	}
	if err := json.Unmarshal(b, &StarBase); err != nil {
		panic(fmt.Sprintf("Cannot parse json data from 'stars.data': %+v", err))
	}
	num_stars = len(StarBase)
	star_data_modified = false
}
