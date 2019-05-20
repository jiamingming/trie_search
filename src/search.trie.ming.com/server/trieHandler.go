package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Result struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data []Response `json:"data"`
}
type Response struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func TrieHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("error...")
			w.WriteHeader(http.StatusBadGateway)
		}
	}()

	defer r.Body.Close()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}


	dict_type := r.Form.Get("t")
	prefix := r.Form.Get("pre")


	result := Result{}

	resp := []Response{}
	if dict_type == `brand` {
		resp = prefixSearch(prefix, BTrie)
	}

	result.Data = resp
	if len(resp) == 0 {
		result.Code = (int(20))
		result.Msg = "error"
	} else {
		result.Code = (int(10))
		result.Msg = `success`
	}

	resu, _ := json.Marshal(result)
	fmt.Fprint(w, string(resu))

}

// todo 拼音前缀搜索 例如 故宫博物院 可根据gu 检索出来 名称拼音前缀为 gu 的条目
func prefixSearch(pre string, trie *Trie) []Response {
	result := trie.PrefixSearch(pre)
	//对检索结果进行排序
	sort.Strings(result)

	length := len(result)
	name_result := []Response{}
	if length > 10 {
		length = 10
	}

	for i := 0; i < length; i++ {
		res1, _ := trie.Find(result[i])
		if res1 == nil {
			return name_result
		}
		val1 := res1.meta
		response := Response{
			Id:   val1.(Response).Id,
			Name: val1.(Response).Name,
		}

		//jsondata,_ := json.Marshal(response)
		//fmt.Println(string(jsondata))
		name_result = append(name_result, response)
	}
	return name_result

}



//todo 拼音模糊搜索，例如 故宫博物院 可根据 ggbwy 检索出来
func fuzzySearchAndSort(fuzzy string, trie *Trie) []string {
	name_result := []string{}
	if fuzzy == `` {
		return name_result
	}
	start := time.Now()
	result := trie.FuzzySearch(fuzzy)
	//fmt.Println(len(result))
	for i := 0; i < len(result); i++ {
		res1, _ := trie.Find(result[i])
		val1 := res1.meta
		name_result = append(name_result, val1.(string))
		//fmt.Println(val1)
	}
	fmt.Println(len(result))
	end := time.Now()
	fmt.Println(end.Sub(start))
	return name_result
}
