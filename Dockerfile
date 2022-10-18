FROM node:lts as dependencies
WORKDIR /web
COPY package.json yarn.lock ./
RUN yarn install

FROM node:lts as builder
WORKDIR /web
COPY . .
RUN yarn build

FROM node:lts as runner
WORKDIR /web
ENV NODE_ENV production
COPY --from=builder /web/public ./public
COPY --from=builder /web/node_modules ./node_modules
COPY --from=builder /web/package.json ./package.json

EXPOSE 3000
CMD ["yarn", "start"]