package persistence

import (
    "fmt"
    "testing"
)

// 定义一个自定义结构体
type Person struct {
    Name string
    Age  int
}

func TestStorage(t *testing.T) {
    // 创建Person类型的存储实例
    personStorage := NewStorage[[]*Person]()

    // 创建要保存的数据
    person := []*Person{
        {
            Name: "张三",
            Age:  30,
        },
        {
            Name: "李四",
            Age:  28,
        },
    }

    // 保存数据
    err := personStorage.Save(person, "person.dat")
    if err != nil {
        fmt.Printf("保存失败: %v\n", err)
        return
    }

    // 加载数据，返回值类型是确定的Person
    loadedPerson, err := personStorage.Load("person.dat")
    if err != nil {
        fmt.Printf("加载失败: %v\n", err)
        return
    }

    for _, p := range loadedPerson {
        fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
    }

    //fmt.Printf("加载的数据: %+v\n", loadedPerson)
}
