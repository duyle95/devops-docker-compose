FROM node:20-alpine

WORKDIR /app

RUN chown -R node:node /app

COPY --chown=node:node package*.json ./

USER node

RUN npm install

COPY --chown=node:node . .

CMD [ "node", "index.js" ]