# VividNote - 现代化的笔记应用

VividNote 是一个基于 Wails + Vue3 开发的现代化笔记应用，支持原始笔记和原子笔记两种记录方式。

## 功能特点

### 原始笔记
- 支持富文本编辑
- 实时保存
- 标题自动更新
- 支持删除操作

### 原子笔记
- 卡片式布局
- 自定义背景颜色
- 快速编辑和预览
- 支持搜索功能
- 支持颜色标记
- 动态高亮显示

## 技术栈

### 前端
- Vue 3
- Vue Router
- Tailwind CSS
- Heroicons
- TipTap Editor

### 后端
- Go
- Wails
- GORM
- SQLite

## 开发环境要求

- Go 1.17+
- Node.js 14+
- npm 或 yarn
- Wails CLI

## 安装和运行

1. 克隆项目
```bash
git clone https://github.com/yourusername/vividnote.git
cd vividnote
```
2. 安装依赖

- 安装前端依赖
```bash
cd frontend
npm install
cd ..
```
- 安装 Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
3. 运行开发环境
```bash
wails dev
```
4. 构建应用
```bash
wails build
```

## 项目结构
```
vividnote/
├── frontend/ # 前端 Vue 项目
│ ├── src/
│ │ ├── components/ # Vue 组件
│ │ ├── views/ # 页面视图
│ │ ├── router/ # 路由配置
│ │ └── App.vue # 根组件
├── internal/ # Go 后端代码
│ ├── data/ # 数据模型
│ └── service/ # 业务逻辑
├── main.go # 主入口
└── wails.json # Wails 配置
```

## 使用说明

### 原始笔记
1. 点击左侧"原始笔记"
2. 点击"+"创建新笔记
3. 使用富文本编辑器编写内容
4. 自动保存

### 原子笔记
1. 点击左侧"原子笔记"
2. 点击"+"创建新的原子笔记
3. 选择背景颜色
4. 编写标题和内容
5. 点击保存

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

[GNU License](LICENSE)

## 联系方式

- 项目维护者：[airhandsome]
- Email：[airhandsome@qq.com]

## 致谢

感谢所有为这个项目做出贡献的开发者。

## 项目截图

### 主界面
![主界面](images/origin-note.png)

### 原始笔记
![原始笔记编辑](images/origin-note.png)

### 原子笔记
![原子笔记展示](images/atomic.png)

### 搜索功能
![搜索功能](images/search.png)
