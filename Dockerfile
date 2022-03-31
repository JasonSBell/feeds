FROM node:16-alpine

WORKDIR /app

# Add necessary modules and dependencies.
COPY package* .
RUN npm install

# Add the application files.
COPY bin/ ./bin
COPY lib/ ./lib
COPY env_defaults .env

# Add best practices.
ENV NODE_ENV production
USER node

ENTRYPOINT [ "node" ]
CMD ["bin/api.js"]