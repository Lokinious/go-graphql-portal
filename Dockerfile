# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download Go modules
RUN go mod download

# Build the Go app
RUN go build -o main .

# Expose port 8081 for the application
EXPOSE 8081

# Command to run the executable
CMD ["./main"]
