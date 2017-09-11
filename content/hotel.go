package content

import (
	"fmt"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Hotel struct {
	item.Item

	Available   bool   `json:"available"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Rate        int    `json:"rate"`
	Address string `json:"address"`
	// Count       int    `json:"count"`
	Total    int    `json:"total"`
	PhotoURL string `json:"photo_u_r_l"`
	Phone    string `json:"phone"`

	Rooms []Room
}

// MarshalEditor writes a buffer of html to edit a Hotel within the CMS
// and implements editor.Editable
func (h *Hotel) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(h,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Hotel field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Name", h, map[string]string{
				"label":       "Name",
				"type":        "text",
				"placeholder": "Enter the Name here",
			}),
		},
		editor.Field{
			View: editor.Checkbox("Available", h, map[string]string{}, map[string]string{
				"true": "Available",
			}),
		},
		editor.Field{
			View: editor.Richtext("Description", h, map[string]string{
				"label":       "Description",
				"placeholder": "Enter the Description here",
			}),
		},
		editor.Field{
			View: editor.File("PhotoURL", h, map[string]string{
				"label":       "Photo",
				"type":        "text",
				"placeholder": "Enter the Photo here",
			}),
		},
		// editor.Field{
		// 	View: editor.Input("Rate", h, map[string]string{
		// 		"label":       "Rate",
		// 		"type":        "text",
		// 		"placeholder": "Enter the Rate here",
		// 	}),
		// },
		editor.Field{
			View: editor.Input("Address", h, map[string]string{
				"label":       "Address",
				"type":        "text",
				"placeholder": "Enter the Address here",
			}),
		},

		editor.Field{
			View: editor.Input("Phone", h, map[string]string{
				"label":       "Phone",
				"type":        "text",
				"placeholder": "Enter the Phone here",
			}),
		},
		// editor.Field{
		// 	View: editor.Input("Count", h, map[string]string{
		// 		"label":       "Count",
		// 		"type":        "text",
		// 		"placeholder": "Enter the Count here",
		// 	}),
		// },
		// editor.Field{
		// 	View: editor.Input("Total", h, map[string]string{
		// 		"label":       "Total",
		// 		"type":        "text",
		// 		"placeholder": "Enter the Total here",
		// 	}),
		// },
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Hotel editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Hotel"] = func() interface{} { return new(Hotel) }
}

// String defines how a Hotel is printed. Update it using more descriptive
// fields from the Hotel struct type
func (h *Hotel) String() string {
	return fmt.Sprintf("Hotel: %s", h.Name)
}
