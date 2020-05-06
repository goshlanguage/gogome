package engine

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
