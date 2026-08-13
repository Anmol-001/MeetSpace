FROM node:22-alpine AS builder

WORKDIR /app

# Copy package.json and lockfile
COPY package.json package-lock.json* ./

# Install dependencies with legacy peer deps as required by Phase 0
RUN npm install --legacy-peer-deps

# Copy the rest of the frontend source
COPY . .

# Accept VITE_API_URL as a build argument
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build the application
RUN npm run build

# Final stage: Serve with Nginx
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy custom nginx configuration for SPA routing
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
