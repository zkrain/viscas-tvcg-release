FROM node:14.18-alpine
WORKDIR /frontend
COPY . .
RUN npm ci && npm run postinstall
EXPOSE 8080
ENTRYPOINT ["npm", "run", "serve"]
