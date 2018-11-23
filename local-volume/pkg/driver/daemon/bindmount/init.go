package bindmount

import "github.com/elastic/stack-operators/local-volume/pkg/driver/model"

// Init returns a successful response when the driver is ready
func (d *Driver) Init() model.Response {
	return model.Response{
		Status:  model.StatusSuccess,
		Message: "driver is available",
		Capabilities: model.Capabilities{
			Attach: false, // only implement mount and unmount
		},
	}
}