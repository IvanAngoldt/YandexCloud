FROM node:20 as build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Копируем собранные файлы наружу
CMD ["sh", "-c", "cp -r /app/build/* /output && echo 'Frontend built and copied'"]
