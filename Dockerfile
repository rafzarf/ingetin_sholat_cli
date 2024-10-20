# Use the official Golang image as a parent image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o ingetin_sholat_cli .

# Run the app when the container launches
CMD ["./ingetin_sholat_cli"]
