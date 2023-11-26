package field

import "testing"

type mapping struct {
	field, expect string
}

var useCases = []mapping{
	{"id", "id = ?"},
	{"idGt", "id > ?"},
	{"idGe", "id >= ?"},
	{"idLt", "id < ?"},
	{"idLe", "id <= ?"},
	{"idNot", "id != ?"},
	{"idNe", "id <> ?"},
	{"idEq", "id == ?"},
}

func TestProcess(t *testing.T) {

	for _, useCase := range useCases {
		t.Run(useCase.field, func(t *testing.T) {
			actual := Process(useCase.field)
			expect := useCase.expect
			if actual != expect {
				t.Errorf("Expected: %s, but got %s", expect, actual)
			}
		})
	}

}
