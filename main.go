/**
 * @author mch
 */

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)
const KEY = "换成自己的key"
type Message struct {
	Content string `form:"content" binding:"required"`
}
func main()  {
	g := gin.Default()
	g.GET("/gpt",gpt)
	g.Run(":8000")
}
func gpt(c *gin.Context) {
	fmt.Println("gpt....")
	message := new(Message)
	if err := c.ShouldBind(message);err != nil {
		c.JSON(500,gin.H{
			"code": "10000",
			"msg": err.Error(),
		})
		return
	}


	s,err := RequestPost(message.Content)
	if err != nil {
		fmt.Println(err.Error())

	}else {
		splitString(s)
	}
}
/*
   "model"=>"gpt-3.5-turbo",
   "messages"=>[["role"=>"user","content"=>"golang"]],
   "temperature" => 0,
   "stream"=>true,
   "n"=> 1,
 */

/*
$key="";
$header=[
    "Content-Type: application/json",
    "Authorization: Bearer ".$key
];
 */
// msg 输入的内容
func RequestPost(msg string) (string,error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	body := make(map[string]interface{})
	message := make([]map[string]interface{},0)
	message = append(message,map[string]interface{}{
		"role": "user",
		"content": msg,
	})

	body["model"] = "gpt-3.5-turbo"
	body["messages"] = message
	body["temperature"] = 0
	body["stream"] = true
	body["n"] = 1
	jsondata,_ := json.Marshal(body)
	req,err := http.NewRequest("POST","https://openai.1rmb.tk/v1/chat/completions",bytes.NewReader(jsondata))
	if err != nil {
		return "",err
	}
	req.Header.Set("Content-Type","application/json")
	req.Header.Set("Authorization","Bearer " + KEY)
	resp,err := client.Do(req)
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	byt,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "",err
	}
	return string(byt),nil
}
//解析chatgpt 返回的数据
func splitString(str string) (string,error) {
	result := ""
	if index := strings.Index(str,"[DONE]");index >= 0 {
		str = strings.Replace(str,"[DONE]","{",1)

		strArr := strings.Split(str,"data: {")

		contents := make([]string,0,len(strArr))
		for index,item := range strArr{
			fmt.Println("index:",index)

			item = "{" + item
			mp := make(map[string]interface{})
			_ = json.Unmarshal([]byte(item),&mp)
			if choices,ok := mp["choices"];ok {
				if delta,ok := choices.([]interface{});ok {
					if itm,ok := delta[0].(map[string]interface{});ok {
						if delta,ok := itm["delta"];ok {
							if content,ok := delta.(map[string]interface{});ok {
								if c,ok := content["content"];ok{
									contents = append(contents,c.(string))
								}

							}
						}
					}
				}
			}
		}
		result = strings.Join(contents,"")
		fmt.Println(result)
	}

	return result,nil
}