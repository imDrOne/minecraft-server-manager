package pagination

import "testing"

func TestNewPageMetadata(t *testing.T) {
	tests := []struct {
		name    string
		page    uint64
		size    uint64
		total   uint64
		wantErr bool
		want    PageMetadata
	}{
		{
			name:    "Valid metadata with exact pages",
			page:    1,
			size:    10,
			total:   100,
			wantErr: false,
			want:    PageMetadata{page: 1, size: 10, total: 100, pages: 10},
		},
		{
			name:    "Valid metadata with extra page",
			page:    1,
			size:    7,
			total:   100,
			wantErr: false,
			want:    PageMetadata{page: 1, size: 7, total: 100, pages: 15}, // 100 / 7 = 14.28 → 15
		},
		{
			name:    "Invalid: Page is zero",
			page:    0,
			size:    10,
			total:   100,
			wantErr: true,
		},
		{
			name:    "Invalid: Size is zero",
			page:    1,
			size:    0,
			total:   100,
			wantErr: true,
		},
		{
			name:    "Total records zero",
			page:    1,
			size:    10,
			total:   0,
			wantErr: false,
			want:    PageMetadata{page: 1, size: 10, total: 0, pages: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageMetadata(tt.page, tt.size, tt.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.page != tt.want.page || got.size != tt.want.size || got.total != tt.want.total || got.pages != tt.want.pages {
					t.Errorf("NewPageMetadata() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestCalcPages(t *testing.T) {
	tests := []struct {
		name  string
		size  uint64
		total uint64
		want  uint64
	}{
		{"Exact pages", 10, 100, 10},
		{"Extra page needed", 7, 100, 15},
		{"One item per page", 1, 100, 100},
		{"Zero total records", 10, 0, 0},
		{"Size greater than total", 50, 30, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result uint64
			calcPages(&result, tt.size, tt.total)
			if result != tt.want {
				t.Errorf("calcPages() = %d, want %d", result, tt.want)
			}
		})
	}
}

func TestPageMetadataGetters(t *testing.T) {
	pm := PageMetadata{page: 2, size: 15, total: 45, pages: 3}

	if pm.Page() != 2 {
		t.Errorf("Page() = %d, want %d", pm.Page(), 2)
	}
	if pm.Size() != 15 {
		t.Errorf("Size() = %d, want %d", pm.Size(), 15)
	}
	if pm.Total() != 45 {
		t.Errorf("Total() = %d, want %d", pm.Total(), 45)
	}
	if pm.Pages() != 3 {
		t.Errorf("Pages() = %d, want %d", pm.Pages(), 3)
	}
}

func TestNewPageRequest(t *testing.T) {
	tests := []struct {
		name    string
		page    uint64
		size    uint64
		wantErr bool
		want    PageRequest
	}{
		{"Valid request", 1, 10, false, PageRequest{page: 1, size: 10}},
		{"Invalid: Page is zero", 0, 10, true, PageRequest{}},
		{"Invalid: Size is zero", 1, 0, true, PageRequest{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageRequest(tt.page, tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.page != tt.want.page || got.size != tt.want.size) {
				t.Errorf("NewPageRequest() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestToPageMeta(t *testing.T) {
	tests := []struct {
		name  string
		req   PageRequest
		total uint64
		want  PageMetadata
	}{
		{"Exact pages", PageRequest{page: 1, size: 10}, 100, PageMetadata{page: 1, size: 10, total: 100, pages: 10}},
		{"Extra page needed", PageRequest{page: 2, size: 7}, 100, PageMetadata{page: 2, size: 7, total: 100, pages: 15}},
		{"Zero total records", PageRequest{page: 1, size: 10}, 0, PageMetadata{page: 1, size: 10, total: 0, pages: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.req.ToPageMeta(tt.total)
			if got.page != tt.want.page || got.size != tt.want.size || got.total != tt.want.total {
				t.Errorf("ToPageMeta() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	tests := []struct {
		name string
		req  PageRequest
		want uint64
	}{
		{"First page", PageRequest{page: 1, size: 10}, 0},
		{"Second page", PageRequest{page: 2, size: 10}, 10},
		{"Third page", PageRequest{page: 3, size: 5}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.Offset(); got != tt.want {
				t.Errorf("Offset() = %d, want %d", got, tt.want)
			}
		})
	}
}
