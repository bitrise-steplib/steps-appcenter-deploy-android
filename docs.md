only single universal apk is supported
https://docs.microsoft.com/en-us/appcenter/build/react-native/android/#63-building-multiple-apks

aab is not supported, it is possible to upload, but impossible to set it to a tester group:
```
POST /v0.1/apps/tamaspapik-bitrise.io/myapp/releases/22/groups HTTP/1.1
Host: api.appcenter.ms

{"id":"7f720936-ff27-4d93-b7b9-af64890c18bb","mandatory_update":false,"notify_testers":false}
HTTP/2.0 400 Bad Request
Content-Length: 100
Access-Control-Allow-Credentials: false
Access-Control-Allow-Headers:
Access-Control-Allow-Methods: POST
Access-Control-Allow-Origin: *
Access-Control-Max-Age: 0
Content-Type: application/json; charset=utf-8
Date: Wed, 20 Nov 2019 09:53:07 GMT
Server: nginx/1.15.8
Strict-Transport-Security: max-age=15724800; includeSubDomains
X-Powered-By: Express

{"code":"not_supported","message":"Error: Distribution of .aab is not supported for groups/testers"}
invalid status code: 400, url: https://api.appcenter.ms/v0.1/apps/tamaspapik-bitrise.io/myapp/releases/22/groups
```