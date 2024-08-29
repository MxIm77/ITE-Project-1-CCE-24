# Backend API Will be Run on a Ubuntu Server By Default.

# Get Golang Version x.xx as golang-builder to build the files:
FROM golang:1.23 AS dependancies

# Procedures that will likely Not Be Changed unless neccessary.
# DO NOT CHANGE THIS IF YOU DONT KNOW WHAT YOU ARE DOING!

# Set The Build Working Directory As /Backend/Source
WORKDIR /Backend/Source

# Run the go mod download Command To Install Neccessary Build Dependancies
COPY go.mod go.sum ./
RUN go mod download

# Copy All Files Before Start Of Build
COPY ./base ./base
COPY ./config ./config
COPY ./source ./source
COPY ./sql ./sql
COPY ./main.go ./main.go
COPY ./Makefile ./Makefile
COPY ./Phoenicia-Digital.log ./Phoenicia-Digital.log

# Run The Commands To Build Backend
RUN mkdir dist
RUN go build -v -o dist main.go
RUN mkdir -p dist/config
RUN cp config/.env dist/config/.env
RUN mkdir -p dist/sql
RUN cp sql/* dist/sql

# Use Ubuntu For Final Container To Run The Golang Backend
FROM ubuntu:22.04

# Create Working Directory To Host The Backend Files.
WORKDIR /srv/backend

# Copy Built Files from the Build Stage
COPY --from=dependancies /Backend/Source/dist .

# Do not change this Unless You Have Changed main.go <file name> or the default build output name!
CMD [ "./main" ]