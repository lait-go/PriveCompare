FROM node:20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache curl bash python3

COPY frontend/package.json ./
RUN npm install

COPY frontend/ .
COPY scripts/fetch-images-by-name.py scripts/fetch-images-by-name.py
RUN if [ "$(find public/images/products -name '*.jpg' 2>/dev/null | wc -l)" -lt 100 ]; then \
      pip3 install --break-system-packages ddgs 2>/dev/null || pip3 install --break-system-packages duckduckgo-search; \
      python3 scripts/fetch-images-by-name.py || true; \
    fi

ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=$VITE_API_URL

RUN npm run build

FROM nginx:1.25-alpine

COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
