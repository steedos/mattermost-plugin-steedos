# Mattermost Steedos Plugin

此插件功能：
- 提供`GET /plugins/steedos/startup`接口，获取steedos应用的authToken, 便于Mattermost访问steedos应用接口
- 提供`POST /plugins/steedos/workflow/webhook`接口，用于将审批王审批信息推送至华信，便于用户及时获取审批动态
- 提供`POST /plugins/steedos/creator/object_webhook`接口，用于将Creator系统中对象的增加修改删除信息推送至华信，便于用户及时获取对象动态

## 安装
转到此Github存储库的[发行版页面](https://github.com/steedos/mattermost-plugin-steedos/releases)，下载最新版本。 您可以在Mattermost[系统控制台](https://docs.mattermost.com/developer/oauth-2-0-applications.html)中上传此文件以安装插件。


## 开发说明
克隆代码到你的`$GOPATH`下:
```
git clone --depth 1 https://github.com/steedos/mattermost-plugin-steedos
```

构建你的插件:
```
make dist
```

这将生成一个插件文件（支持多种操作系统），以便上传到您的Mattermost服务器(如：`http://localhost:8065`):
```
dist/steedos-0.1.0.tar.gz
```

## 使用说明
访问你的Mattermost服务，新建一个[oauth2](https://docs.mattermost.com/developer/oauth-2-0-applications.html)应用，先[启用服务](https://docs.mattermost.com/administration/config-settings.html#enable-oauth-2-0-service-provider)才能新建

上传插件可以访问[系统控制台](https://docs.mattermost.com/developer/oauth-2-0-applications.html)

将插件文件上传后，需要设置相应的参数，参数设置好之后启用插件：
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
`GET /plugins/steedos/startup`，用于获取steedos authToken等信息，返回JSON对象:
```
{
    "authToken": "", (steedos应用的认证authToken，在调用steedos应用接口时传入)
    "url": "", (steedos应用的api地址)
    "userId": "" (steedos用户id)
}
```
`POST /plugins/steedos/workflow/webhook`接口，用于将审批王审批信息推送至华信，便于用户及时获取审批动态，使用此功能准备工作如下：
- 在华炎审批系统(如`https://cn.steedos.com`)中配置webhooks（需要工作区管理员身份），URL参数可配置为华信系统接受地址(如`https://messenger.steedos.cn/plugins/steedos/workflow/webhook`，即本接口)

`POST /plugins/steedos/creator/object_webhook`接口，用于将Creator系统中对象的增加修改删除信息推送至华信，便于用户及时获取对象动态，使用此功能准备工作如下：
- 在creator系统（如[华炎合同管理系统](https://github.com/steedos/steedos-contracts-app)）中配置object_webhooks（需要工作区管理员身份），URL参数可配置为华信系统接受地址(如`https://messenger.steedos.cn/plugins/steedos/creator/object_webhook`，即本接口)
## mattermost网页端调用`/plugins/steedos/startup`接口示例:
```js
    import request from 'superagent';
    doGet = async (url, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        try {
            const response = await request.
                get(url).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }

    let data = await doGet(`http://mattermostUrl/plugins/steedos/startup`);
```