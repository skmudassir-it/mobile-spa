# Build stage: compile static site (Astro / Vite / plain static)
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY . .
RUN npm run build

# Runtime: tiny nginx
FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
RUN sed -i 's/listen\s*80;/listen 3000;/' /etc/nginx/conf.d/default.conf \
    && sed -i 's/listen\s*\[::\]:80;/listen [::]:3000;/' /etc/nginx/conf.d/default.conf
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
