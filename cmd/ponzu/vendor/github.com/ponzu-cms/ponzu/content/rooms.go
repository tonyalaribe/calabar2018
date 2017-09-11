package content

import (
	"fmt"
	"net/http"

	"github.com/bosssauce/reference"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Room struct {
	item.Item

	Type        string `json:"type"`
	Hotel       string `json:"hotel"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	RackRate    int    `json:"rack_rate"`
	LionsRate   int    `json:"lions_rate"`
	Distance    int    `json:"distance"`
	// RoomPhoto   string `json:"room_photo"`
}

// MarshalEditor writes a buffer of html to edit a Room within the CMS
// and implements editor.Editable
func (a *Room) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(a,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Room field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("Type", a, map[string]string{
				"label":       "Type",
				"type":        "text",
				"placeholder": "Enter the Type here",
			}),
		},
		editor.Field{
			View: reference.Select("Hotel", a, map[string]string{
				"label": "Hotel",
			},
				"Hotel",
				`{{ .name  }} `,
			),
		},
		editor.Field{
			View: editor.Textarea("Description", a, map[string]string{
				"label":       "Description",
				"placeholder": "Enter the Description here",
			}),
		},
		editor.Field{
			View: editor.Input("Quantity", a, map[string]string{
				"label":       "Quantity",
				"type":        "number",
				"placeholder": "Enter the Quantity here",
			}),
		},
		editor.Field{
			View: editor.Input("RackRate", a, map[string]string{
				"label":       "RackRate",
				"type":        "number",
				"placeholder": "Enter the Rack Rate here",
			}),
		},
		editor.Field{
			View: editor.Input("LionsRate", a, map[string]string{
				"label":       "LionsRate",
				"type":        "number",
				"placeholder": "Enter the Lions Rate here",
			}),
		},
		editor.Field{
			View: editor.Input("Distance", a, map[string]string{
				"label":       "Distance",
				"type":        "number",
				"placeholder": "Enter the Distance here",
			}),
		},
		// editor.Field{
		// 	View: editor.File("RoomPhoto", a, map[string]string{
		// 		"label":       "Photo",
		// 		"type":        "text",
		// 		"placeholder": "Enter the Photo here",
		// 	}),
		// },
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Room editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Room"] = func() interface{} { return new(Room) }
}

// func SearchMapping() (*mapping.IndexMappingImpl, error) {

// }

func IndexContent() bool {
	return true
}
func (a *Room) String() string {
	return a.Type
}

func (a *Room) Update(res http.ResponseWriter, req *http.Request) error {
	return nil
}
