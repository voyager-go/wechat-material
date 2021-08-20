package main

func main() {
	LoadConfig()
	Conn()
	Handler()
	// defer XEngine.Close()
	// http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	// http.Handle("/ueditor/", http.StripPrefix("/ueditor/", http.FileServer(http.Dir("ueditor"))))
	// http.HandleFunc("/index", ParseTemplate)
	// err := http.ListenAndServe(":8899", nil)
	// if err != nil {
	// 	log.Fatalln("ListenAndServe: ", err)
	// }
}
