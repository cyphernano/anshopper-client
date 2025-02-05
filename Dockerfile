# Use a lightweight web server
FROM nginx:alpine

# Copy WebAssembly files
COPY anshopper_wasm/index.html /usr/share/nginx/html/
COPY anshopper_wasm/main.wasm /usr/share/nginx/html/
COPY anshopper_wasm/wasm.js /usr/share/nginx/html/

# Expose port 80
EXPOSE 8080

# Run Nginx
CMD ["nginx", "-g", "daemon off;"]
