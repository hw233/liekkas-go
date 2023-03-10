### **gm操作指南**

1. 增加`GET/POST` 请求

```go
func (g *HttpGetHandler) Test(ctx context.Context, param *httputil.Params) (string, error) {
    glog.Info(param.Get("A"))
    return param.Encode(), nil
}
```
在`GET/POST`对应的HttpHandler增加函数,函数的全称且第一个字母小写就是http请求的path。
比如此时path是 `/test`,`GET`接受的参数在url中，
`POST`接受的参数在body以json格是发送
2. 返回格式

- success情况下 `err_code`为0，`err_msg`为空，`data`为任意json格式
```json
    {
      "err_code": 0,
      "err_msg": "",
      "data": ""
    }
 ```

- err情况下  `err_code`不为0，`err_msg`不为空，`data`为`null`
```json
    {
      "err_code": 1001,
      "err_msg": "Param error",
      "data": null
    }
```