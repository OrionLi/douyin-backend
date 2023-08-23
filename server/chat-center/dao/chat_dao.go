package dao

import (
	"bytes"
	"chat-center/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/OrionLi/douyin-backend/pkg/pb"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

// GetAllMessagesByToUserId 查询指定 toUserId 的消息列表
// toUserId : 消息接收者
// fromUserId : 消息发送者
func GetAllMessagesByToUserId(toUserId int64, fromUserId int64) ([]*pb.Message, error) {

	query := fmt.Sprintf(`{
		"query": {
			"bool": {
				"must": [
					{
						"term": 
						{
							"to_user_id": %d 
						}
					},
					{
						"term": 
						{
							"from_user_id": %d 
						}
					}
				]
			}
		}
	}`, toUserId, fromUserId)

	// 创建排序条件
	sort := `{
		"create_time": {
			"order": "desc"
		}
	}`

	// 设置返回的结果数量
	size := 999

	res, err := ESClient.Search(
		ESClient.Search.WithIndex("douyin_messages"), // 索引名
		ESClient.Search.WithBody(strings.NewReader(query)),
		ESClient.Search.WithSort(sort),
		ESClient.Search.WithSize(size),
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
func GetMessageByToUserId(time time.Time, toUserId int64, fromUserId int64) ([]*pb.Message, error) {
	query := fmt.Sprintf(`
		{
			"query": {
				"bool": {
					"must": [
						{ "term": { "to_user_id": %d } },
						{ "term": { "from_user_id": %d } },
						{ "range": { "create_time": { "gt": "%s" } } }
					]
				}
			}
		}`, toUserId, fromUserId, time.Format("2006-01-02 15:04:05"))

	// 创建排序条件
	sort := `{
		"create_time": {
			"order": "desc"
		}
	}`

	// 设置返回的结果数量
	size := 999

	res, err := ESClient.Search(
		ESClient.Search.WithIndex("douyin_messages"), // 索引名
		ESClient.Search.WithBody(strings.NewReader(query)),
		ESClient.Search.WithSort(sort),
		ESClient.Search.WithSize(size),
		ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	messageList := GetMessageList(res)

	return messageList, nil
}

// GetMessageList 将 Elasticsearch 返回的 Response 转换为 message.Message 列表
func GetMessageList(res *esapi.Response) []*pb.Message {
	var messageList []*pb.Message
	var responseMap map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&responseMap); err != nil {
		utils.LogrusObj.Error("<GetMessage>, Error decoding JSON: ", err, " [be from]:", res)
		log.Fatalf("Error decoding JSON: %s", err)
	}

	hits := responseMap["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		id := int(source["id"].(float64))
		toUserId := int(source["to_user_id"].(float64))
		fromUserId := int(source["from_user_id"].(float64))
		content := source["content"].(string)
		createTime := source["create_time"].(string)

		messageList = append(messageList, &pb.Message{
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
func SendMessage(message *pb.Message) error {
	docJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "douyin_messages",
		DocumentID: strconv.FormatInt(message.Id, 10),
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true", // 在索引之后刷新以使文档可搜索
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("error sending the request: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing the body: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		return fmt.Errorf("[%s] Error response: %s", res.Status(), res.String())
	}

	return nil
}
