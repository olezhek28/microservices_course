package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gomodule/redigo/redis"
)

const (
	fieldName     = "name"
	fieldLastName = "last_name"
	fieldAge      = "age"
	fieldEmail    = "email"
)

type User struct {
	Name     string `redis:"name"`
	LastName string `redis:"last_name"`
	Age      int    `redis:"age"`
	Email    string `redis:"email"`
}

func main() {
	// Подключаемся к Redis
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer conn.Close() // можно использовать closer из 3-й недели

	setAndGet(conn)
	hsetAndHGet(conn)
}

func setAndGet(conn redis.Conn) {
	key := gofakeit.UUID()
	value := gofakeit.FirstName()

	// Сохраняем пару ключ-значение
	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Fatalf("failed to set key: %v", err)
	}

	// Получаем значение по ключу
	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		log.Fatalf("failed to get key: %v", err)
	}

	fmt.Printf("Пара ключ-значение (%s: %s)\n\n", key, value)
}

func hsetAndHGet(conn redis.Conn) {
	hashKey := gofakeit.UUID()
	fields := map[string]string{
		fieldName:     gofakeit.FirstName(),
		fieldLastName: gofakeit.LastName(),
		fieldAge:      strconv.FormatInt(int64(gofakeit.IntRange(0, 100)), 10),
		fieldEmail:    gofakeit.Email(),
	}

	// Сохраняем значения в хеш-таблицу
	var err error
	for field, value := range fields {
		_, err = conn.Do("HSET", hashKey, field, value)
		if err != nil {
			log.Fatalf("failed to set hash field: %v", err)
		}
	}

	// Получаем значения из хеш-таблицы разными способами
	printMapFieldsByOne(conn, hashKey)
	fmt.Println()
	printMapFields(conn, hashKey)
	fmt.Println()
	printMapFieldsByStruct(conn, hashKey)
}

func printMapFieldsByOne(conn redis.Conn, hashKey string) {
	name, err := redis.String(conn.Do("HGET", hashKey, fieldName))
	if err != nil {
		log.Fatalf("failed to get hash field \"%v\": %v", fieldName, err)
	}

	lastName, err := redis.String(conn.Do("HGET", hashKey, fieldLastName))
	if err != nil {
		log.Fatalf("failed to get hash field \"%v\": %v", fieldLastName, err)
	}

	age, err := redis.String(conn.Do("HGET", hashKey, fieldAge))
	if err != nil {
		log.Fatalf("failed to get hash field \"%v\": %v", fieldAge, err)
	}

	email, err := redis.String(conn.Do("HGET", hashKey, fieldEmail))
	if err != nil {
		log.Fatalf("failed to get hash field \"%v\": %v", fieldEmail, err)
	}

	fmt.Printf("Данные пользователя с идентифкатором %s:\n", hashKey)
	fmt.Printf("Имя: %s\n", name)
	fmt.Printf("Фамилия: %s\n", lastName)
	fmt.Printf("Возраст: %s\n", age)
	fmt.Printf("Email: %s\n", email)
}

func printMapFields(conn redis.Conn, hashKey string) {
	hashMap, err := redis.StringMap(conn.Do("HGETALL", hashKey))
	if err != nil {
		log.Fatalf("failed to get all hash fields: %v", err)
	}

	fmt.Printf("Данные пользователя с идентифкатором (полученные разом) %s:\n", hashKey)
	fmt.Printf("%#v\n", hashMap)
}

func printMapFieldsByStruct(conn redis.Conn, hashKey string) {
	values, err := redis.Values(conn.Do("HGETALL", hashKey))
	if err != nil {
		log.Fatalf("failed to get all hash fields: %v", err)
	}

	var user User
	err = redis.ScanStruct(values, &user)
	if err != nil {
		log.Fatalf("failed to scan hash fields to struct: %v", err)
	}

	fmt.Printf("Данные пользователя с идентифкатором (распаршенные в структуру) %s:\n", hashKey)
	fmt.Printf("%#v\n", user)
}
