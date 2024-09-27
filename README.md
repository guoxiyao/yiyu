### 代码协作

参考文档 [Git 分支与冲突解决探究](http://zxlmdonnie.cn/archives/1709691289692)

#### 首次获取项目
1. 从Github上克隆项目 `git clone https://github.com/Werun-backend/YiyuGoDemo.git`
2. 创建本地开发分支：`git checkout -b dev_zfr` (xxx一般为姓名首字母小写，如dev_zfr);

#### 修改代码前
3. 首先保证在自己的开发分支，每次写代码前，都要更新本地代码和远程main保持一致：
    1. 更新本地main分支：`git pull origin main`
    2. 将本地main分支合并到自己的分支 `git merge main dev_zfr`

#### 修改代码后
4. 代码提交需先push到自己的远程开发分支上：
   1. `git add 文件名` 或者使用 Goland 的图形化界面
   2. `git commit -m '日志信息'`
   3. `git push origin dev_zfr`;
6. 在GitHub的Pull Requests中手动将xxx_dev合并到main分支：点击合并请求合并分支到main，一路提交即可

### 前端对接

- 对入参、出参的 `json` 实体进行数据转换(dto vo)，以适配前端需求。
- 封装统一返回类组件 HTTP状态码
- 封装分页组件

把上述规范写一个文档（MarkDown）进行说明。


