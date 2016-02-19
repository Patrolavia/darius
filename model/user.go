// This file is part of Darius. See License.txt for license information.

package model

// User record
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"-"`
	Image string `json:"image"`
}
