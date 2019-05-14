# Mattermost Steedos Plugin

此插件用于获取steedos应用的authToken, 便于Mattermost访问steedos应用接口

## 使用说明
克隆代码:
```
git clone --depth 1 https://github.com/steedos/mattermost-plugin-steedos
```

构建你的插件:
```
make dist
```

这将生成一个插件文件（支持多种操作系统），以便上传到您的Mattermost服务器(如：`http://localhost:8065`):
```
dist/com.mattermost.steedos-0.1.0.tar.gz
```

访问你的Mattermost服务，新建一个[oauth2](https://docs.mattermost.com/developer/oauth-2-0-applications.html)应用，先[启用服务](https://docs.mattermost.com/administration/config-settings.html#enable-oauth-2-0-service-provider)才能新建

上传插件可以访问[系统控制台](https://docs.mattermost.com/developer/oauth-2-0-applications.html)

将生成的插件文件上传后，需要设置相应的参数，参数设置好之后启用插件：
```
- URL (steedos应用的服务地址)
- API URL (steedos应用的api地址)
- API Key (值为新建的oauth2应用的`客户端 ID`)
- API Secret (值为新建的oauth2应用的`客户端密钥`)
- Webhook Secret (用于验证Webhook到Mattermost的密钥)
```

同时为确保插件可用，steedos应用也需要新建一个OAuth2应用：
```
- 客户端ID (值为新建的oauth2应用的`客户端 ID`)
- 密钥 (值为新建的oauth2应用的`客户端密钥`)
```

## 开放接口
`GET /plugins/com.mattermost.steedos/startup`，用于获取steedos authToken等信息:
```json
{
    "authToken": "", (steedos应用的认证authToken，在调用steedos应用接口时传入)
    "url": "", (steedos应用的api地址)
    "userId": "" (steedos用户id)
}
```