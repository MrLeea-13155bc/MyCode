package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// go get -u go.mongodb.org/mongo-driver/mongo
// mongoDB  依赖
var client *mongo.Client

func init() {
	// 客户端配置
	// url 格式： mongodb://(Username):(password)@(host):(port)
	// 官方例子：mongodb://user:pass@sample.host:27017
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	// 连接到数据库
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Panicln(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Link Success!")
}

//Client 获取client方便其他包调用
func Client() *mongo.Client {
	return client
}
func main() {

}

//CRUD
type User struct {
	UserName string
	Password string
}

func Insert() {
	s1 := User{"Peter", "peterspassword"}
	s2 := User{"John", "johnspassword"}
	// 获取Server数据库中的User数据集
	collection := Client().Database("Server").Collection("User")
	// 使用 InsertOne
	//lastInsertId,err := collection.InsertOne(context.TODO(),s1)
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(lastInsertId.InsertedID)
	//lastInsertId,err = collection.InsertOne(context.TODO(),s2)
	//if err != nil {
	//log.Fatal(err)
	//}
	//fmt.Println(lastInsertId.InsertedID)
	// 使用 InsertMany
	//需要通过使用空接口数组存储每条数据的内容
	users := []interface{}{s1, s2}
	lastInsertIds, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(lastInsertIds.InsertedIDs)
}

/**
  BSON => 二进制JSON类型
  D：一个BSON文档。这种类型应该在顺序重要的情况下使用，比如MongoDB命令。
  M：一张无序的map。它和D是一样的，只是它不保持顺序。
  A：一个BSON数组。
  E：D里面的一个元素。
*/

func Update() {
	// 更新数据的条件，类似于sql中的Where Username = 'Peter'，其中Key为条件需要的列名，Value为该列匹配数据
	filter := bson.D{{"UserName", "Peter"}}
	// 更新的方法 以及更新之后的内容
	// 其中Key为更新方法，Value为更新之后的内容 👇
	// 更新内容为bson.D 其中Key为修改的列名，Value为修改之后的数据
	update := bson.D{
		{"$set", bson.D{
			{"UserName", "March"},
		}},
	}
	collection := Client().Database("Server").Collection("User")
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	// 输出匹配的行数以及调整的行数
	fmt.Println(updateResult.MatchedCount, updateResult.ModifiedCount)
	// UpdateMany与UpdateOne类似
}

func Search() {
	var result User
	filter := bson.D{{"Username", "Peter"}}
	collection := Client().Database("Server").Collection("User")

	// FindOne
	findResult := collection.FindOne(context.TODO(), filter)
	err := findResult.Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	// Find
	//cur,err := collection.Find(context.TODO(),filter)
	////别忘了关闭
	//defer cur.Close(context.TODO())
	//var results []User
	//// 类似于 mysql 就不多说了
	//for cur.Next(context.TODO()) {
	//  err = cur.Decode(&result)
	//  if err != nil {
	//    log.Fatal(err)
	//  }
	//  results = append(results, result)
	//}
	////统一返回游标的错误信息
	//if err := cur.Err();err != nil{
	//  log.Fatal(err)
	//}
	//fmt.Println(results)
}

func Delete() {
	//删除的语法类似于查找
	filter := bson.D{{"Username", "Peter"}}
	collection := Client().Database("Server").Collection("User")

	// 与FindOne不同，这里会返回错误，需要进行处理
	findResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(findResult.DeletedCount)
	// DeleteMany 与 deleteOne 语法一致 删除时会将所有匹配的都进行删除
}
