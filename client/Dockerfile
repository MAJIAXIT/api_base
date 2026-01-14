# FROM node:22-alpine AS build

# WORKDIR /app

# COPY package*.json ./

# RUN npm install
# COPY . .
# RUN npm run build

# COPY nginx.conf /etc/nginx/templates/default.conf.template

# FROM nginx:stable-alpine

# COPY --from=build /app/out /usr/share/nginx/html
# COPY --from=build /etc/nginx/templates/default.conf.template /etc/nginx/templates/default.conf.template

# EXPOSE 80

# CMD ["nginx", "-g", "daemon off;"]


FROM nginx:stable-alpine

COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
