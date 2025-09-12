# Use nginx to serve static files
FROM nginx:alpine

# Copy the HTML file to nginx web directory
COPY public/index.html /usr/share/nginx/html/

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
