# Use an official Golang image as the base
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy everything from the current directory to the container
COPY . .

# Download dependencies and build the application
RUN go mod tidy && go build -o user-service ./cmd/main.go

# Expose the port
EXPOSE 50051

# Command to run the application
CMD ["./user-service"]
