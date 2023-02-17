# Use the official Golang image as the base image
FROM golang:1.20.1-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go app
RUN go build -o nandobot-go .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./nandobot-go"]