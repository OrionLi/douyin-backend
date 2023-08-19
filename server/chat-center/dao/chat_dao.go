package dao

import (
	"bytes"
	"chat-center/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"strings"
	"time"
)

// GetAllMessagesByToUserId 查询指定 toUserId 的消息列表
// toUserId : 消息接收者
// fromUserId : 消息发送者
func GetAllMessagesByToUserId(toUserId int64, fromUserId int64) ([]model.Message, error) {

	query := fmt.Sprintf(`{
			"query": {
				"bool": {
					"must": [
						{ "match": { "toUserId": %d } },
						{ "match": { "fromUserId": %d } }
					]
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
// toUserId : 消息接收者
// fromUserId : 消息发送者
// time : 最后一次阅读消息的时间
func GetMessageByToUserId(time time.Time, toUserId int64, fromUserId int64) ([]model.Message, error) {
	query := fmt.Sprintf(`
		{
			"query": {
				"bool": {
					"must": [
						{ "match": { "toUserId": %d } },
						{ "match": { "fromUserId": %d } },
						{ "range": { "createTime": { "gt": "%s" } } }
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
func GetMessageList(res *esapi.Response) []model.Message {
	var messageList []model.Message
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

		messageList = append(messageList, model.Message{
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
func SendMessage(message model.Message) error {
	docJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index:   "douyin_messages",
		Body:    bytes.NewReader(docJSON),
		Refresh: "true", // 在索引之后刷新以使文档可搜索
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("error sending the request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("[%s] Error response: %s", res.Status(), res.String())
	}

	return nil
}
