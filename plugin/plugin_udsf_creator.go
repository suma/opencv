package plugin

// UDSFCreator is an interface to register the user defined stream function.
type UDSFCreator interface {
	// CreateStreamFunction returns user defined stream function. This function
	// only deal with generic function.
	CreateStreamFunction() interface{}
	// TypeName returns name of registration.
	// Example:
	//  a type name is "hoge_udsf", then
	//    SELECT * FROM hoge_udsf(...);
	TypeName() string
}
