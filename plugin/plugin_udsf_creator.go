package plugin

// PluginUDSFCreator is an interface to get BQL UDSF function and to registration.
type PluginUDSFCreator interface {
	// CreateStreamFunction returns user plug-in function. THis function only deal
	// with generic function.
	CreateStreamFunction() interface{}
	// TypeName returns name of registration.
	// Example:
	//  a type name is "hoge_udsf", then
	//    SELECT * FROM hoge_udsf(...);
	TypeName() string
}
