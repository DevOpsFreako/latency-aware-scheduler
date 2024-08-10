# Use an official nginx image
FROM nginx:alpine

# Copy the HTML file to the nginx web root
COPY client-example.html /usr/share/nginx/html/index.html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
