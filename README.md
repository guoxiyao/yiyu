### 代码协作

#### 首次拉取项目
1. 从git上拉取项目 `git clone https://github.com/Werun-backend/YiyuGoDemo.git`
2. 创建本地开发分支：`git checkout -b dev_zfr` (xxx一般为姓名首字母小写，如dev_zfr);

#### 修改代码前
3. 每次写代码前，都要更新本地代码和远程main保持一致：
    1. 更新本地main分支：`git pull origin main`
    2. 将本地main分支合并到自己的分支 `git merge main dev_zfr`

#### 修改代码后
4. 代码提交需先push到自己的开发分支上：`git add --all` && `git commit -m '日志信息'` && `git push origin dev_zfr`;
5. 在GitHub Pull request 手动将xxx_dev 合并到main分支，点击合并请求合并分支到 main