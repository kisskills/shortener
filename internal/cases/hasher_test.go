package cases

import "testing"

func TestHasher(t *testing.T) {
	type args struct {
		originalLink string
	}
	tests := []struct {
		name         string
		originalLink string
		want         string
	}{
		{
			name:         "OK Google:)",
			originalLink: "google.com",
			want:         "kwUIdA9TOq",
		},
		{
			name:         "OK Ozon",
			originalLink: "ozon.ru",
			want:         "naNLoN5x10",
		},
		{
			name:         "Job Ozon ?!",
			originalLink: "job.ozon.ru",
			want:         "Vmvi9FRmmK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hasher(tt.originalLink); got != tt.want {
				t.Errorf("Hasher() = %v, want %v", got, tt.want)
			}
		})
	}
}
