FROM nginx:alpine

# Copy the frontend files to nginx html directory
COPY public /usr/share/nginx/html
COPY src /usr/share/nginx/html/src

# Copy custom nginx configuration if needed
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"] 