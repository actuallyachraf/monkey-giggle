package eval

import "github.com/actuallyachraf/monkey-giggle/object"

var builtins = map[string]*object.BuiltIn{
	"len":    object.GetBuiltInByName("len"),
	"head":   object.GetBuiltInByName("head"),
	"tail":   object.GetBuiltInByName("tail"),
	"last":   object.GetBuiltInByName("last"),
	"append": object.GetBuiltInByName("append"),
}
