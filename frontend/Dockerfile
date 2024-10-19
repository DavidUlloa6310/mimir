# Dockerfile for Next.js frontend
FROM node:18-alpine as builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
ENV NEXT_PUBLIC_BACKEND_IP=https://mimir-backend-97586522204.us-central1.run.app/
RUN npm run build

FROM node:18-alpine as runner
WORKDIR /app

COPY --from=builder /app/ ./

EXPOSE 3000

CMD ["npm", "start"]
