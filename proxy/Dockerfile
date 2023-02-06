FROM nginx:1.13-alpine
EXPOSE 80
COPY default.conf /etc/nginx/conf.d/default.conf
COPY --from=front /app/dist /usr/share/nginx/html