package data

import "testing"

type input struct {
	titleInput       string
	artistNamesInput []string
}

func TestCreateHash(t *testing.T) {
	tests := []struct {
		name        string
		inputs      [2]input
		shouldMatch bool
	}{
		{
			"Same tracks should have same hashes",
			[2]input{
				{
					"get him back!",
					[]string{"Olivia Rodrigo"},
				},
				{
					"get him back!",
					[]string{"Olivia Rodrigo"},
				},
			},
			true,
		},
		{
			"Artist order shouldn't matter",
			[2]input{
				{
					"Lean Beef Patty",
					[]string{"JPEGMAFIA", "Danny Brown"},
				},
				{
					"Lean Beef Patty",
					[]string{"Danny Brown", "JPEGMAFIA"},
				},
			},
			true,
		},
		{
			"All artists should have to be the same",
			[2]input{
				{
					"CAROUSEL",
					[]string{"Aries"},
				},
				{
					"CAROUSEL",
					[]string{"Travis Scott"},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := CreateHash(tt.inputs[0].titleInput, tt.inputs[0].artistNamesInput)
			b := CreateHash(tt.inputs[1].titleInput, tt.inputs[1].artistNamesInput)

			if (a == b && !tt.shouldMatch) || (a != b && tt.shouldMatch) {
				t.Fail()
			}
		})
	}
}
