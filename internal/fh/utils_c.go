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

package fh

import (
	"time"
)

/* This routine will return a random int between 1 and max, inclusive.
   It uses the so-called "Algorithm M" method, which is a combination
   of the congruential and shift-register methods. */

var last_random uint64 = 1924085713 /* Random seed. */

func rnd(max int) int {
	var a, b, c, cong_result, shift_result uint64

	/* For congruential method, multiply previous value by the prime number 16417. */
	a = last_random
	b = last_random << 5
	c = last_random << 14
	cong_result = a + b + c /* Effectively multiply by 16417. */

	/* For shift-register method, use shift-right 15 and shift-left 17 with no-carry addition (i.e., exclusive-or). */
	a = last_random >> 15
	shift_result = a ^ last_random
	a = shift_result << 17
	shift_result ^= a

	last_random = cong_result ^ shift_result

	a = last_random & 0x0000FFFF

	return int((a*uint64(max))>>16) + 1
}

// Seed random number generator
func Seed(seed uint64) {
	last_random = seed
	n := rnd(100) + rnd(200) + rnd(300)
	for i := 0; i < n; i++ {
		rnd(10)
	}
}

// SeedFromTime
func SeedFromTime() {
	Seed(uint64(time.Now().UnixNano()))
}
