package gamescube

import (
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/models"
)

func TestVendor_parseProduct(t *testing.T) {
	type args struct {
		s        *goquery.Selection
		cardName string
	}

	sampleProduct, err := os.Open("testdata/sample_product.html")
	if err != nil {
		t.Fatalf("failed to open testdata file: %v", err)
	}

	sampleDoc, err := goquery.NewDocumentFromReader(sampleProduct)
	if err != nil {
		t.Fatalf("failed to read testdata file: %v", err)
	}

	tests := []struct {
		name    string
		v       *Vendor
		args    args
		want    []models.Offering
		wantErr bool
	}{
		{
			name: "Should parse product",
			v: &Vendor{
				logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
			},
			args: args{
				s:        sampleDoc.Find("li.product"),
				cardName: "Army of the Damned",
			},
			want: []models.Offering{
				{
					Name:       "Army of the Damned",
					Price:      196,
					Quantity:   12,
					Set:        "Commander 2016",
					Condition:  "NM-Mint, English",
					ImgURI:     "https://crystalcommerce-assets.nyc3.cdn.digitaloceanspaces.com/photos/6324633/medium/gnAACpUTlw_EN.png",
					ProductURI: "https://www.thegamescube.com/catalog/magic_singles-commander_singles-commander_2016/army_of_the_damned/394253",
					StoreSKU:   "5302723",
					VendorID:   "gamescube",
					Properties: []string{},
				},
				{
					Name:       "Army of the Damned",
					Price:      186,
					Quantity:   4,
					Set:        "Commander 2016",
					Condition:  "Light Play, English",
					ImgURI:     "https://crystalcommerce-assets.nyc3.cdn.digitaloceanspaces.com/photos/6324633/medium/gnAACpUTlw_EN.png",
					ProductURI: "https://www.thegamescube.com/catalog/magic_singles-commander_singles-commander_2016/army_of_the_damned/394253",
					StoreSKU:   "8466752",
					VendorID:   "gamescube",
					Properties: []string{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.v.parseProducts(tt.args.s, tt.args.cardName)
			// Ignore parse time
			for i := range got {
				got[i].CreatedAt = 0
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Vendor.parseProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vendor.parseProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkCardNameMatch(t *testing.T) {
	type args struct {
		productName string
		cardName    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Should match exact name",
			args: args{
				productName: "Swords to Plowshares",
				cardName:    "Swords to Plowshares",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should not match different name",
			args: args{
				productName: "Swords to Plowshares",
				cardName:    "Forest",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Should not match slightly different name",
			args: args{
				productName: "Swords to Plowshares",
				cardName:    "Swords",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Should match name with foil",
			args: args{
				productName: "Damn - Foil",
				cardName:    "Damn",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should match name with foil and borderless",
			args: args{
				productName: "Damn - Foil - Borderless",
				cardName:    "Damn",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should match name with parenthesis",
			args: args{
				productName: "Forest (Full art)",
				cardName:    "Forest",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should match name with parenthesis and numbers",
			args: args{
				productName: "Island (B08/010) - Foil Double Sided",
				cardName:    "Island",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should not match name with parenthesis",
			args: args{
				productName: "Island (Full art)",
				cardName:    "Forest",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Should match name with dash",
			args: args{
				productName: "Accident-Prone Apprentice - Foil",
				cardName:    "Accident-Prone Apprentice",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should match name with non-ascii characters",
			args: args{
				productName: "Tale of Tinúviel - Foil",
				cardName:    "Tale of Tinuviel",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Should match name with non-ascii characters (inverse)",
			args: args{
				productName: "Tale of Tinuviel - Foil",
				cardName:    "Tale of Tinúviel",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkCardNameMatch(tt.args.productName, tt.args.cardName)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkCardNameMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkCardNameMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
