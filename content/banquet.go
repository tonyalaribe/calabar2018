package content

import (
	"fmt"
	"net/http"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Banquet struct {
	item.Item

	RegistrationId string `json:"registration_id"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Amount         int    `json:"amount"`
}

// MarshalEditor writes a buffer of html to edit a Banquet within the CMS
// and implements editor.Editable
func (b *Banquet) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(b,
		// Take note that the first argument to these Input-like functions
		// is the string version of each Banquet field, and must follow
		// this pattern for auto-decoding and auto-encoding reasons:
		editor.Field{
			View: editor.Input("RegistrationId", b, map[string]string{
				"label":       "RegistrationId",
				"type":        "text",
				"placeholder": "Enter the RegistrationId here",
			}),
		},
		editor.Field{
			View: editor.Input("Phone", b, map[string]string{
				"label":       "Phone",
				"type":        "text",
				"placeholder": "Enter the Phone here",
			}),
		},
		editor.Field{
			View: editor.Input("Email", b, map[string]string{
				"label":       "Email",
				"type":        "text",
				"placeholder": "Enter the Email here",
			}),
		},
		editor.Field{
			View: editor.Input("Amount", b, map[string]string{
				"label":       "Amount",
				"type":        "text",
				"placeholder": "Enter the Amount here",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Banquet editor view: %s", err.Error())
	}

	return view, nil
}

func init() {
	item.Types["Banquet"] = func() interface{} { return new(Banquet) }
}

// String defines how a Banquet is printed. Update it using more descriptive
// fields from the Banquet struct type
func (b *Banquet) String() string {
	//return fmt.Sprintf("Banquet: %s", b.UUID)
	return b.Email
}

func (b *Banquet) Create(res http.ResponseWriter, req *http.Request) error {
	return nil
}
