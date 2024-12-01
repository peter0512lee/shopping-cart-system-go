FROM node:18-alpine

WORKDIR /app

# 複製 package.json 和 package-lock.json（如果存在）
COPY package*.json ./

# 安裝依賴
RUN npm install

# 複製源代碼和配置文件
COPY . .

# 確保構建腳本可以訪問 Tailwind
RUN npm install -g tailwindcss postcss autoprefixer

EXPOSE 3000

CMD ["npm", "start"]