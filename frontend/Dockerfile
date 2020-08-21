# start build or layer from the Node.js 12 (LTS) image
FROM node:12

# Add Maintainer Info
LABEL maintainer="Olumide Ogundele <olumideralph@gmail.com>"

# Specify directory inside the image in which all commands will run
WORKDIR /usr/src/app

# install app dependencies
# Copy package files i.e package.json and package-lock.json
# and install dependencies.
# Copies package.json and package-lock.json to Docker environment
COPY package*.json ./
RUN npm install
RUN npm install react-scripts -g

# Copies everything over to Docker environment
COPY . .

# Expose port 3000 to the outside world
EXPOSE 3000

# Default command to run when starting the container
CMD ["npm", "start"]
