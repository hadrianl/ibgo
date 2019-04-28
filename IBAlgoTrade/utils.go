package IBAlgoTrade

import (
	"reflect"

	"github.com/hadrianl/ibgo/ibapi"
)

//Default try to init the object with the default tag
func Default(o interface{}) {
	t := reflect.TypeOf(o).Elem()
	v := reflect.ValueOf(o).Elem()

	fieldCount := t.NumField()

	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)

		// if v.Field(i).Kind() == reflect.Struct {
		// 	Create(v.Field(i).Interface())
		// 	fmt.Println(v.Field(i))
		// }

		if defaultValue, ok := field.Tag.Lookup("default"); ok {

			switch defaultValue {
			case "UNSETFLOAT":
				v.Field(i).SetFloat(ibapi.UNSETFLOAT)
			case "UNSETINT":
				v.Field(i).SetInt(ibapi.UNSETINT)
			case "-1":
				v.Field(i).SetInt(-1)
			case "true":
				v.Field(i).SetBool(true)
			default:
				panic("Unknown defaultValue:")
			}
		}
		// fmt.Printf("value:***%v***", field.Tag)
		// fmt.Printf("type:***%v***", field.Type)

	}
}
