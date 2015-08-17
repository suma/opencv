package plugin

// UDFCreator is an interface to register the user defined function.
type UDFCreator interface {
	// CreateFunction returns user defined function. This function only deal
	// with generic function.
	CreateFunction() interface{}
	// TypeName returns name of registration.
	// Example:
	//  a type name is "hoge_udf", then
	//    SELECT hoge_udf(...) FROM hoge_stream;
	TypeName() string
}
