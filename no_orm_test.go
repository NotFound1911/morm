// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package morm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NotFound1911/morm/internal/valuer"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sqlTestSuite struct {
	suite.Suite

	// 配置字段
	driver string
	dsn    string

	// 初始化字段
	db *sql.DB
}

// TearDownTest 方法则是在每个测试用例执行完毕后自动调用的方法，用于完成与测试用例相对应的清理工作
func (s *sqlTestSuite) TearDownTest() {
	_, err := s.db.Exec("DELETE FROM test_model;")
	if err != nil {
		s.T().Fatal(err)
	}
}

// SetupSuite 方法会在测试套件中的所有测试运行之前执行
// 完成见表操作
func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *sqlTestSuite) TestCRUD() {
	t := s.T() // s.T()方法是一个测试辅助工具函数，用于获取当前测试环境中的*testing.T对象
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 或者 Exec(xxx)
	res, err := db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES (1, 'Tom', 18, 'Jerry')")
	if err != nil {
		t.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

	rows, err := db.QueryContext(context.Background(),
		"SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` LIMIT 1")
	if err != nil {
		t.Fatal()
	}
	for rows.Next() {
		tm := &TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "Tom", tm.FirstName)
	}

	// 或者 Exec(xxx)
	res, err = db.ExecContext(ctx, "UPDATE `test_model` SET `first_name` = 'changed' WHERE `id` = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	affected, err = res.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != 1 {
		t.Fatal(err)
	}

	row := db.QueryRowContext(context.Background(), "SELECT `id`, `first_name`,`age`, `last_name` FROM `test_model` LIMIT 1")
	if row.Err() != nil {
		t.Fatal(row.Err())
	}
	tm := &TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "changed", tm.FirstName)
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "sqlite3",
		dsn:    "file:test.db?cache=shared&mode=memory",
	})
}

func (TestModel) CreateSQL() string {
	return `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
}

// memoryDB 返回一个基于内存的 ORM，它使用的是 sqlite3 内存模式。
func memoryDB(t *testing.T, opts ...DBOption) *DB {
	orm, err := Open("sqlite3", "file:test.db?cache=shared&mode=memory", opts...)
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

func memoryDBWithDB(db string, t *testing.T) *DB {
	orm, err := Open("sqlite3", fmt.Sprintf("file:%s.db?cache=shared&mode=memory", db))
	if err != nil {
		t.Fatal(err)
	}
	return orm
}

// 执行 go test -bench=BenchmarkExec_Get -benchmem -benchtime=10000x
// goos: windows
// goarch: amd64
// pkg: github.com/NotFound1911/morm
// cpu: AMD Ryzen 5 2600X Six-Core Processor
// BenchmarkExec_Get/unsafe-12                10000            263590 ns/op            3287 B/op        110 allocs/op
// BenchmarkExec_Get/reflect-12               10000            724458 ns/op            3466 B/op        118 allocs/op
// PASS
// ok      github.com/NotFound1911/morm    10.579s
func BenchmarkExec_Get(b *testing.B) {
	db, err := Open("sqlite3", "file:benchmark_get.db?cache=shared&mode=memory")
	if err != nil {
		b.Fatal(err)
	}
	_, err = db.db.Exec(TestModel{}.CreateSQL()) // 建表
	if err != nil {
		b.Fatal(err)
	}
	res, err := db.db.Exec("INSERT INTO  `test_model`(`id`, `first_name`, `age`, `last_name`)VALUES (?,?,?,?) ",
		1, "aa", 18, "bb")
	if err != nil {
		b.Fatal(err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		b.Fatal(err)
	}
	if affected == 0 {
		b.Fatal()
	}
	b.Run("unsafe", func(b *testing.B) {
		db.valCreator = valuer.NewUnsafeValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("reflect", func(b *testing.B) {
		db.valCreator = valuer.NewReflectValue
		for i := 0; i < b.N; i++ {
			_, err = NewSelector[TestModel](db).Get(context.Background())
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
