package todo

import (
	"errors"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	type Test struct {
		name        string
		listID      uint32
		description string
		comments    string
		dueDate     time.Time
		labels      []string
		done        bool
		wantErr     error
	}

	tests := []Test{
		{
			name:        "SuccessMakeTheBed",
			listID:      0,
			description: "Make the bed.",
			comments:    "",
			dueDate:     parseTime(t, "2021-02-01T00:00:00Z"),
			labels:      []string{"", ""},
			done:        false,
			wantErr:     nil,
		},
		{
			name:        "ErrorEmptyDescription",
			listID:      0,
			description: "",
			comments:    "",
			dueDate:     parseTime(t, "2021-02-01T00:00:00Z"),
			labels:      []string{"", ""},
			done:        false,
			wantErr:     ErrEmptyDescription,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			todo := Todo{
				ListID:      test.listID,
				Description: test.description,
				Comments:    test.comments,
				DueDate:     test.dueDate,
				Labels:      test.labels,
			}

			err := Create(todo)

			if !errors.Is(err, test.wantErr) {
				t.Errorf("got error %v; want %v", err, test.wantErr)
			}
		})
	}
}

func parseTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
