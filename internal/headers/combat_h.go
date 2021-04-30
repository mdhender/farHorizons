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

package headers

/* Maximum number of battle locations for all players. */
const MAX_BATTLES = 50

/* Maximum number of ships at a single battle. */
const MAX_SHIPS = 200

/* Maximum number of engagement options that a player may specify
   for a single battle. */
const MAX_ENGAGE_OPTIONS = 20

type char = byte
type long = int64
type cstring struct {
	b []byte
}

type battle_data struct {
	x, y, z, num_species_here char
	spec_num                  [MAX_SPECIES]char
	summary_only              [MAX_SPECIES]char
	transport_withdraw_age    [MAX_SPECIES]char
	warship_withdraw_age      [MAX_SPECIES]char
	fleet_withdraw_percentage [MAX_SPECIES]char
	haven_x                   [MAX_SPECIES]char
	haven_y                   [MAX_SPECIES]char
	haven_z                   [MAX_SPECIES]char
	special_target            [MAX_SPECIES]char
	hijacker                  [MAX_SPECIES]char
	can_be_surprised          [MAX_SPECIES]char
	enemy_mine                [MAX_SPECIES][MAX_SPECIES]char
	num_engage_options        [MAX_SPECIES]char
	engage_option             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]char
	engage_planet             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]char
	ambush_amount             [MAX_SPECIES]long
}

/* Types of combatants. */
const SHIP = 1
const NAMPLA = 2
const GENOCIDE_NAMPLA = 3
const BESIEGED_NAMPLA = 4

/* Types of special targets. */
const TARGET_WARSHIPS = 1
const TARGET_TRANSPORTS = 2
const TARGET_STARBASES = 3
const TARGET_PDS = 4

/* Types of actions. */
const DEFENSE_IN_PLACE = 0
const DEEP_SPACE_DEFENSE = 1
const PLANET_DEFENSE = 2
const DEEP_SPACE_FIGHT = 3
const PLANET_ATTACK = 4
const PLANET_BOMBARDMENT = 5
const GERM_WARFARE = 6
const SIEGE = 7

/* Special types. */
const NON_COMBATANT = 1

type action_data struct {
	num_units_fighting     int
	fighting_species_index [MAX_SHIPS]int
	num_shots              [MAX_SHIPS]int
	shots_left             [MAX_SHIPS]int
	weapon_damage          [MAX_SHIPS]long
	shield_strength        [MAX_SHIPS]long
	shield_strength_left   [MAX_SHIPS]long
	original_age_or_PDs    [MAX_SHIPS]long
	bomb_damage            [MAX_SHIPS]long
	surprised              [MAX_SHIPS]char
	unit_type              [MAX_SHIPS]char
	fighting_unit          [MAX_SHIPS]*cstring
}
