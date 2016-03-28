## 小拉勾
#### 目录结构
crawler 爬取拉勾网工作数据  
util    通用工具  
web     gin web框架，展示数据库中的数据  
templates   前端模板  
assets  css、js、img  

#### 功能
1. 爬取拉勾网工作数据保存到数据库
2. 前端使用Web的形式进行展示
3. 前端JS请求高德地图数据

#### TODO
1. job关键字，查询功能增加（目前写死的Go）
2. 高德地图起点关键字，查询功能增加（目前写死的北京-北苑）
3. 前端UI优化（增加查询按钮、使用AdminLTE替换Bootstrap）
4. 爬虫定时自动化（目前需要手动启动爬虫）
5. 部署到Digitalocean
