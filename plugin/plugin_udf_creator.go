package plugin

// PluginUDFCreator is an interface to get BQL UDF function and to registration.
type PluginUDFCreator interface {
	// CreateFunction returns user -lug-in function. This function only deal
	// with generic function.
	CreateFunction() interface{}
	// TypeName returns name of registration.
	// Example:
	//  a type name is "hoge_udf", then
	//    SELECT hoge_udf(...) FROM hoge_stream;
	TypeName() string
}
