package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInitConfigure(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: testing.CoverMode()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_config_GetDBConfig(t *testing.T) {
	type fields struct {
		database Database
		storage  Storage
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "can get db info", fields: fields{
			database: Database{Path: "./Database.db"},
			storage:  Storage{},
		}, want: "./Database.db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := config{
				Database: tt.fields.database,
				Storage:  tt.fields.storage,
			}
			got := c.GetDBConfig()
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetDBConfig differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func Test_config_GetStorageConfig(t *testing.T) {
	type fields struct {
		database Database
		storage  Storage
	}
	tests := []struct {
		name        string
		fields      fields
		wantDist    string
		wantSubdist string
	}{
		{name: testing.CoverMode(), fields: fields{
			database: Database{},
			storage:  Storage{Dist: "test", SubDist: "test2"},
		}, wantDist: "test", wantSubdist: "test2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config{
				Database: tt.fields.database,
				Storage:  tt.fields.storage,
			}
			gotDist, gotSubdist := c.GetStorageConfig()
			if diff := cmp.Diff(gotDist, tt.wantDist); diff != "" {
				t.Errorf("GetStorageConfig differs: (-got +want)\n%s", diff)
			}
			if diff := cmp.Diff(gotSubdist, tt.wantSubdist); diff != "" {
				t.Errorf("GetStorageConfig differs: (-got +want)\n%s", diff)
			}
		})
	}
}
