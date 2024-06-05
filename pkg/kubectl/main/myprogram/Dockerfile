# Start from the official Golang base image
FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download and install any required Go modules
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Expose port 12345 to the outside world
EXPOSE 12345

# Command to run the executable
CMD ["./main"]
