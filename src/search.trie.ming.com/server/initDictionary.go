package server

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//test - db-server
const (
	USERNAME = "root"
	PASSWORD = "123456"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "search"
)



func InitMysqlServer() *sql.DB {
	//init mysql
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("Open mysql failed,err:%v\n", err)
		return nil
	}
	DB.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(100)                  //设置最大连接数
	DB.SetMaxIdleConns(16)                   //设置闲置连接数
	return DB

}

type User struct {
	Id int64 `db:"id"`
}

type Store struct {
	Id   uint64 `db:"store_id"`
	Name string `db:"store_name"`
}
type Brand struct {
	Id   uint64 `db:"brand_id"`
	Name string `db:"brand_name"`
}

var BTrie *Trie
var DB = InitMysqlServer()

func InitQueryBrand() {
	for {
		start := time.Now()
		rows, err := DB.Query("select b.brand_id,b.brand_name from dict_brand b ")
		brand := new(Brand)
		defer func() {
			if rows != nil {
				rows.Close()
			}
		}()
		if err != nil {
			fmt.Printf("Query failed,err:%v", err)

		}
		BTrie = New()
		for rows.Next() {
			err := rows.Scan(&brand.Id, &brand.Name)
			if err != nil {
				fmt.Printf("Scan failed,err:%v", err)
				continue
			}

			response := Response{
				Id:   brand.Id,
				Name: brand.Name,
			}

			store_pinyin := ConvertPinyin(brand.Name)
			BTrie.Add(store_pinyin, response)

		}
		end := time.Now()
		fmt.Println("brand dictionary load done...spend: ", end.Sub(start))
		time.Sleep(time.Hour * 24)
	}

}


