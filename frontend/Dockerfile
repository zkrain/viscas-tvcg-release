FROM node:14.18-alpine
WORKDIR /frontend
COPY . .
RUN npm install
EXPOSE 8080
ENTRYPOINT ["npm", "run", "serve"]

