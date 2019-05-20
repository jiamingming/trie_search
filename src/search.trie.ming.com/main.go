package main

import (
	"fmt"
	"net"
	"net/http"
	"search.trie.ming.com/server"
)

func init() {


	go server.InitQueryBrand()

}
func main() {

	l, err := net.Listen("tcp", ":5678")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	fmt.Println("trie server start...")
	http.HandleFunc("/trie", http.HandlerFunc(server.TrieHandler))
	http.Serve(l, nil)

}
