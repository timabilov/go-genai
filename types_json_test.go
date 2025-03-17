package genai

import (
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    any
		wantErr bool
		target  string // "Schema", "Citation", "TokensInfo"
	}{
		// Schema tests
		{
			name:    "Schema empty",
			jsonStr: `{}`,
			want:    &Schema{},
			wantErr: false,
			target:  "Schema",
		},
		{
			name:    "Schema all fields",
			jsonStr: `{"maxLength": "10", "minLength": "5", "minProperties": "2", "maxProperties": "4", "maxItems": "8", "minItems": "1", "maximum": 10.0, "minimum": 2.0}`,
			want: &Schema{
				MaxLength:     Ptr[int64](10),
				MinLength:     Ptr[int64](5),
				MinProperties: Ptr[int64](2),
				MaxProperties: Ptr[int64](4),
				MaxItems:      Ptr[int64](8),
				MinItems:      Ptr[int64](1),
				Maximum:       Ptr[float64](10.0),
				Minimum:       Ptr[float64](2.0),
			},
			wantErr: false,
			target:  "Schema",
		},
		{
			name:    "Schema invalid maxLength",
			jsonStr: `{"maxLength": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid minLength",
			jsonStr: `{"minLength": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid minProperties",
			jsonStr: `{"minProperties": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid maxProperties",
			jsonStr: `{"maxProperties": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid maxItems",
			jsonStr: `{"maxItems": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid minItems",
			jsonStr: `{"minItems": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid maximum",
			jsonStr: `{"maximum": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid minimum",
			jsonStr: `{"minimum": "abc"}`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},
		{
			name:    "Schema invalid json",
			jsonStr: `{"minimum": "abc"`,
			want:    nil,
			wantErr: true,
			target:  "Schema",
		},

		// Citation tests
		{
			name:    "Citation empty",
			jsonStr: `{}`,
			want:    &Citation{},
			wantErr: false,
			target:  "Citation",
		},
		{
			name:    "Citation all fields",
			jsonStr: `{"endIndex": 10, "license": "MIT", "publicationDate": {"year": 2023, "month": 10, "day": 26}, "startIndex": 5, "title": "Test Title", "uri": "https://example.com"}`,
			want: &Citation{
				EndIndex:        10,
				License:         "MIT",
				PublicationDate: civil.Date{Year: 2023, Month: 10, Day: 26},
				StartIndex:      5,
				Title:           "Test Title",
				URI:             "https://example.com",
			},
			wantErr: false,
			target:  "Citation",
		},
		{
			name:    "Citation missing year",
			jsonStr: `{"publicationDate": {"month": 10, "day": 26}}`,
			want:    nil,
			wantErr: true,
			target:  "Citation",
		},
		{
			name:    "Citation only year",
			jsonStr: `{"publicationDate": {"year": 2023}}`,
			want: &Citation{
				PublicationDate: civil.Date{Year: 2023},
			},
			wantErr: false,
			target:  "Citation",
		},
		{
			name:    "Citation only year and month",
			jsonStr: `{"publicationDate": {"year": 2023, "month": 10}}`,
			want: &Citation{
				PublicationDate: civil.Date{Year: 2023, Month: 10},
			},
			wantErr: false,
			target:  "Citation",
		},
		{
			name:    "Citation invalid json",
			jsonStr: `{"publicationDate": {"year": 2023`,
			want:    nil,
			wantErr: true,
			target:  "Citation",
		},

		// TokensInfo tests
		{
			name:    "TokensInfo empty",
			jsonStr: `{}`,
			want:    &TokensInfo{},
			wantErr: false,
			target:  "TokensInfo",
		},
		{
			name:    "TokensInfo all fields",
			jsonStr: `{"role": "user", "tokenIds": ["1", "2", "3"], "tokens": ["YQ==", "Yg==", "Yw=="]}`,
			want: &TokensInfo{
				Role:     "user",
				TokenIDs: []int64{1, 2, 3},
				Tokens:   [][]byte{[]byte("a"), []byte("b"), []byte("c")},
			},
			wantErr: false,
			target:  "TokensInfo",
		},
		{
			name:    "TokensInfo invalid token id",
			jsonStr: `{"tokenIds": ["1", "a", "3"]}`,
			want:    nil,
			wantErr: true,
			target:  "TokensInfo",
		},
		{
			name:    "TokensInfo invalid json",
			jsonStr: `{"tokenIds": ["1", "2", "3"`,
			want:    nil,
			wantErr: true,
			target:  "TokensInfo",
		},

		// CreateCachedContentConfig tests
		{
			name:    "CreateCachedContentConfig empty",
			jsonStr: `{}`,
			want:    &CreateCachedContentConfig{},
			wantErr: false,
			target:  "CreateCachedContentConfig",
		},
		{
			name:    "CreateCachedContentConfig with expireTime",
			jsonStr: `{"expireTime": "2024-12-31T23:59:59Z"}`,
			want: &CreateCachedContentConfig{
				ExpireTime: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			wantErr: false,
			target:  "CreateCachedContentConfig",
		},
		{
			name:    "CreateCachedContentConfig invalid json",
			jsonStr: `{"expireTime": "2024-12-31T23:59:59Z`,
			want:    nil,
			wantErr: true,
			target:  "CreateCachedContentConfig",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var got any
			switch tt.target {
			case "Schema":
				s := &Schema{}
				err = s.UnmarshalJSON([]byte(tt.jsonStr))
				got = s
			case "Citation":
				c := &Citation{}
				err = c.UnmarshalJSON([]byte(tt.jsonStr))
				got = c
			case "TokensInfo":
				ti := &TokensInfo{}
				err = ti.UnmarshalJSON([]byte(tt.jsonStr))
				got = ti
			case "CreateCachedContentConfig":
				c := &CreateCachedContentConfig{}
				err = json.Unmarshal([]byte(tt.jsonStr), c)
				got = c
			default:
				t.Fatalf("unknown target: %s", tt.target)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("%s.UnmarshalJSON() error = %v, wantErr %v", tt.target, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("%s.UnmarshalJSON() = %v, want %v. Diff: %s", tt.target, got, tt.want, diff)
				}
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
		target  string // "Schema", "Citation", "TokensInfo"
	}{
		// Schema tests
		{
			name:    "Schema empty",
			input:   &Schema{},
			want:    `{}`,
			wantErr: false,
			target:  "Schema",
		},
		{
			name: "Schema all fields",
			input: &Schema{
				MaxLength:     Ptr[int64](10),
				MinLength:     Ptr[int64](5),
				MinProperties: Ptr[int64](2),
				MaxProperties: Ptr[int64](4),
				MaxItems:      Ptr[int64](8),
				MinItems:      Ptr[int64](1),
				Maximum:       Ptr[float64](10.0),
				Minimum:       Ptr[float64](2.0),
			},
			want:    `{"maxLength":"10","minLength":"5","minProperties":"2","maxProperties":"4","maxItems":"8","minItems":"1","maximum":10,"minimum":2}`,
			wantErr: false,
			target:  "Schema",
		},

		// Citation tests
		{
			name:    "Citation empty",
			input:   &Citation{},
			want:    `{}`,
			wantErr: false,
			target:  "Citation",
		},
		{
			name: "Citation all fields",
			input: &Citation{
				EndIndex:        10,
				License:         "MIT",
				PublicationDate: civil.Date{Year: 2023, Month: 10, Day: 26},
				StartIndex:      5,
				Title:           "Test Title",
				URI:             "https://example.com",
			},
			want:    `{"publicationDate":"2023-10-26","endIndex":10,"license":"MIT","startIndex":5,"title":"Test Title","uri":"https://example.com"}`,
			wantErr: false,
			target:  "Citation",
		},

		// TokensInfo tests
		{
			name:    "TokensInfo empty",
			input:   &TokensInfo{},
			want:    `{}`,
			wantErr: false,
			target:  "TokensInfo",
		},
		{
			name: "TokensInfo all fields",
			input: &TokensInfo{
				Role:     "user",
				TokenIDs: []int64{1, 2, 3},
				Tokens:   [][]byte{},
			},
			want:    `{"tokenIds":["1","2","3"],"role":"user"}`,
			wantErr: false,
			target:  "TokensInfo",
		},
		// CreateCachedContentConfig tests
		{
			name:    "CreateCachedContentConfig empty",
			input:   &CreateCachedContentConfig{},
			want:    `{}`,
			wantErr: false,
			target:  "CreateCachedContentConfig",
		},
		{
			name: "CreateCachedContentConfig with expireTime",
			input: &CreateCachedContentConfig{
				ExpireTime: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			want:    `{"expireTime":"2024-12-31T23:59:59Z"}`,
			wantErr: false,
			target:  "CreateCachedContentConfig",
		},
		// GenerateContentResponse tests
		{
			name:    "GenerateContentResponse empty",
			input:   &GenerateContentResponse{},
			want:    `{}`,
			wantErr: false,
			target:  "GenerateContentResponse",
		},
		{
			name: "GenerateContentResponse with createTime",
			input: &GenerateContentResponse{
				CreateTime: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			want:    `{"createTime":"2024-12-31T23:59:59Z"}`,
			wantErr: false,
			target:  "GenerateContentResponse",
		},
		// TunedModelInfo tests
		{
			name:    "TunedModelInfo empty",
			input:   &TunedModelInfo{},
			want:    `{}`,
			wantErr: false,
			target:  "TunedModelInfo",
		},
		{
			name: "TunedModelInfo with createTime and updateTime",
			input: &TunedModelInfo{
				CreateTime: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
				UpdateTime: time.Date(2024, 12, 30, 23, 59, 59, 0, time.UTC),
			},
			want:    `{"createTime":"2024-12-31T23:59:59Z","updateTime":"2024-12-30T23:59:59Z"}`,
			wantErr: false,
			target:  "TunedModelInfo",
		},
		// CachedContent tests
		{
			name:    "CachedContent empty",
			input:   &CachedContent{},
			want:    `{}`,
			wantErr: false,
			target:  "CachedContent",
		},
		// UpdateCachedContentConfig tests
		{
			name:    "UpdateCachedContentConfig empty",
			input:   &UpdateCachedContentConfig{},
			want:    `{}`,
			wantErr: false,
			target:  "UpdateCachedContentConfig",
		},
		{
			name: "UpdateCachedContentConfig with expireTime",
			input: &UpdateCachedContentConfig{
				ExpireTime: time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			want:    `{"expireTime":"2024-12-31T23:59:59Z"}`,
			wantErr: false,
			target:  "UpdateCachedContentConfig",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []byte
			var err error
			switch tt.target {
			case "Schema":
				got, err = json.Marshal(tt.input.(*Schema))
			case "Citation":
				got, err = json.Marshal(tt.input.(*Citation))
			case "TokensInfo":
				got, err = json.Marshal(tt.input.(*TokensInfo))
			case "CreateCachedContentConfig":
				got, err = json.Marshal(tt.input.(*CreateCachedContentConfig))
			case "GenerateContentResponse":
				got, err = json.Marshal(tt.input.(*GenerateContentResponse))
			case "TunedModelInfo":
				got, err = json.Marshal(tt.input.(*TunedModelInfo))
			case "CachedContent":
				got, err = json.Marshal(tt.input.(*CachedContent))
			case "UpdateCachedContentConfig":
				got, err = json.Marshal(tt.input.(*UpdateCachedContentConfig))
			default:
				t.Fatalf("unknown target: %s", tt.target)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("%s.MarshalJSON() error = %v, wantErr %v", tt.target, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if string(got) != tt.want {
					t.Errorf("%s.MarshalJSON() = %v, want %v", tt.target, string(got), tt.want)
				}
			}
		})
	}
}
