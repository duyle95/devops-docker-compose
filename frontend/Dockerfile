FROM node:20 AS build

WORKDIR /app

COPY ./package.json /app/package.json
COPY ./package-lock.json /app/package-lock.json

RUN npm ci
COPY . .
RUN npm run build


FROM nginx:latest

RUN apt-get update -y && \
    apt-get install -y procps apache2-utils && \
    rm -rf /var/lib/apt/lists/*

ENV BASIC_USERNAME=user
ENV BASIC_PASSWORD=pass

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY --from=build /app/dist /usr/share/nginx/dist

COPY run.sh ./
RUN chmod 0755 ./run.sh
CMD [ "./run.sh" ]