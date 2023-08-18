package dao

import (
	"douyin-backend/chat-center/generated/message"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

// GetAllMessagesByToUserId 查询指定 toUserId 的消息列表
//
// 参数：toUserId
//
// 返回值：消息列表，错误
func GetAllMessagesByToUserId(toUserId int, fromUserId int) ([]message.Message, error) {

	query := fmt.Sprintf(`
	{
	  "query": {
	    "match": {
	      "toUserId": %d,
          "fromUserId": %d
	    }
	  }
	}`, toUserId, fromUserId)

	res, err := ESClient.Search(
		ESClient.Search.WithIndex("douyin_messages"), // 索引名
		ESClient.Search.WithBody(strings.NewReader(query)),
		ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	messageList := GetMessageList(res)

	return messageList, nil
}

// GetMessageByToUserId 查询指定 toUserId 的消息列表
//
// 参数：toUserId，preMsgTime
//
// 返回值：消息列表，错误
func GetMessageByToUserId(time time.Time, toUserId int, fromUserId int) ([]message.Message, error) {
	query := fmt.Sprintf(`
	{
	  "query": {
	    "bool": {
	      "must": [
	        {
	          "match": {
	            "toUserId": %d,
				"fromUserId": %d
	          }
	        },
	        {
	          "range": {
	            "createTime": {
	              "gte": "%s"
	            }
	          }
	        }
	      ]
	    }
	  }
	}`, toUserId, fromUserId, time.Format("2006-01-02 15:04:05"))

	res, err := ESClient.Search(
		ESClient.Search.WithIndex("douyin_messages"), // 索引名
		ESClient.Search.WithBody(strings.NewReader(query)),
		ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	messageList := GetMessageList(res)

	return messageList, nil
}

// GetMessageList 将 Elasticsearch 返回的 Response 转换为 message.Message 列表
func GetMessageList(res *esapi.Response) []message.Message {
	var messageList []message.Message
	var responseMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseMap); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	hits := responseMap["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		id := int(source["id"].(float64))
		toUserId := int(source["toUserId"].(float64))
		fromUserId := int(source["fromUserId"].(float64))
		content := source["content"].(string)
		createTime := source["createTime"].(string)

		messageList = append(messageList, message.Message{
			Id:         int64(id),
			ToUserId:   int64(toUserId),
			FromUserId: int64(fromUserId),
			Content:    content,
			CreateTime: createTime,
		})
	}
	return messageList
}

// SendMessage 发送信息
// TODO: 完成发送信息的数据库中添加操作
func SendMessage() error {
	return nil
}
