# start build or layer from the Node.js 12 (LTS) image
FROM node:12-alpine

# Specify directory inside the image in which all commands will run
WORKDIR /app

# install app dependencies
# Copy package files i.e package.json and package-lock.json
# and install dependencies.
COPY package.json ./
COPY package-lock.json ./
RUN npm install
RUN npm install react-scripts -g

# Copy all of the app files into the image
COPY . .

# Expose port 3000 to the outside world
EXPOSE 3000

# Default command to run when starting the container
CMD ["npm", "start"]
