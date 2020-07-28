package cmd

import "testing"

func Test_massage(t *testing.T) {
	type args struct {
		csvPath  string
		jsonPath string
		zipCount int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{"testdata/data.csv", "/tmp/data.json", 3}, false},
		{"noinput", args{"nonexistent.csv", "/tmp/data.json", 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := massage(tt.args.csvPath, tt.args.jsonPath); (err != nil) != tt.wantErr {
				t.Errorf("massage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			data, err := loadData(tt.args.jsonPath)
			if err != nil {
				t.Errorf("loadData() error %v", err)
			}
			if len(data[keyCounters]) != tt.args.zipCount {
				t.Errorf("loadData() expected %v got %v ", tt.args.zipCount, len(data[keyCounters]))
			}

		})
	}
}
