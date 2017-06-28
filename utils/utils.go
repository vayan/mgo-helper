package utils

/*
PanicIfError panics when an unexpected error occurs
*/
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
