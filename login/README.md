## **登录服接口**

### 说明
1. 每个环境都有对应的login服，比如内网开发和内网测试是不同的
2. post的接口参数都放到form中
3. 返回content格式统一为 json``` {"code":0,"time":1648447997,"data":{"user_id":9,"token":"5cacaca70d6638016e763cadd185efcff8f8bee2"}}```,code=0的时候表示没有错误，code>0的时候表示有错误。
### 接口
- 内部注册
  - url: `/v1/inner-register`
  - 参数：
  - `account` 账号
  - `password` 密码
  - 返回：
  ```json
  {"code":0,"time":1648447997,"data":{"user_id":9,"token":"5cacaca70d6638016e763cadd185efcff8f8bee2"}}
  ```

- 内部登录
  - url:`/v1/inner-login`
  - 参数：
  - `account` 账号
  - `password` 密码
  - 返回：
  ```json
  {"code":0,"time":1648447997,"data":{"user_id":9,"token":"5cacaca70d6638016e763cadd185efcff8f8bee2"}}
  ```

- 第三方登录
  - url:`/v1/third-party-login`
  - 参数：
  - `third_party` 第三方sdk名称，小萌就填 koMoe
  - `open_id` 第三方openId
  - `token` 第三方登录token
  - `ios` true or false 表示是否是ios
  
  - 返回：
  ```json
  {"code":0,"time":1648447997,"data":{"user_id":9,"token":"5cacaca70d6638016e763cadd185efcff8f8bee2"}}
  ```

